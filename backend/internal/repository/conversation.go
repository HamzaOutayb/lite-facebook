package repository

import (
	"fmt"

	"social-network/internal/models"
	utils "social-network/pkg"
)

func (data *Database) CreateConversation(conv *models.Conversation) (err error) {
	args := utils.GetExecFields(conv, "ID")

	res, err := data.Db.Exec(fmt.Sprintf(`
		INSERT INTO conversations
		VALUES (NULL, %v) 
	`, utils.Placeholders(len(args))),
		args...)
	if err != nil {
		return
	}

	id, err := res.LastInsertId()
	conv.ID = int(id)

	return
}

func (data *Database) VerifyConversation(id1, id2 int, type_obj string) (err error) {
	param := `WHERE entitie_two_group = ?`
	if type_obj == "private" {
		param = `WHERE entitie_two_user = ?`
	}
	var result int
	err = data.Db.QueryRow(fmt.Sprintf(`
		SELECT id
		FROM conversations
		WHERE entitie_one = ? AND %v`,
		param),
		id1, id2).Scan(&result)
	return
}

func (data *Database) GetConversations(id int) (conversations []models.ConversationsInfo, err error) {
	query := `
		SELECT
			C.*,COALESCE((SELECT content FROM messages M WHERE M.conversation_id = C.id),"") as last_message
		FROM
			conversations C
		WHERE
			(
				C.type = 'private'
				AND (
					C.entitie_one = ?
					OR C.entitie_two_user = ?
				)
			)
			OR (
				C.type = 'group'
				AND (
					C.entitie_one = ?
					OR
					EXISTS (
					SELECT
						1
					FROM
						invites I
					WHERE
						I.group_id = C.entitie_two_group
						AND (
							I.receiver = ?
							OR I.sender = ?
						)
						AND I.status = 'accepted'
					)
				)
			)
	`
	rows, err := data.Db.Query(query, id, id, id, id, id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var conv models.ConversationsInfo
		tab := append(utils.GetScanFields(&conv.Conversation),&conv.LastMessage)
		err1 := rows.Scan(tab...)
		if err1 != nil {
			fmt.Println(err1)
		}
		conversations = append(conversations, conv)
	}

	for i, conv := range conversations {
		if conv.Conversation.Type == "group" {
			row := data.GetGroupById(*conv.Conversation.Entitie_two_group)
			err1 := row.Scan(utils.GetScanFields(&conversations[i].Group)...)
			if err1 != nil {
				fmt.Println(err1)
			}
		} else {
			var err1 error
			idUser := conv.Conversation.Entitie_one
			if id == idUser {
				idUser = *conv.Conversation.Entitie_two_user
			}
			conversations[i].UserInfo, err1 = data.GetUserByID(idUser)
			if err1 != nil {
				fmt.Println(err1)
			}
		}
	}

	return
}
