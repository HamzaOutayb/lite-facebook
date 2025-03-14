package service

import (
	"fmt"

	"social-network/internal/models"
	utils "social-network/pkg"
)

func (S *Service) GreatedGroup(Group *models.Group) (err error) {
	err = S.Database.SaveGroup(Group)
	return
}

func (S *Service) AllGroups() ([]models.Group, error) {
	rows, err := S.Database.Getallgroup()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		if err := rows.Scan(utils.GetScanFields(&group)...); err != nil {
			fmt.Println(err)
			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}

func (S *Service) GetGroupsById(Group *models.Group) (*models.Group, error) {
	row := S.Database.GetGroupById(Group.ID)
	if row == nil {
		return nil, fmt.Errorf("no group found with ID: %s", Group.ID)
	}

	// Scan the row into the Group struct
	if err := row.Scan(utils.GetScanFields(Group)...); err != nil {
		return nil, fmt.Errorf("error scanning group data: %v", err)
	}

	return Group, nil
}

func (S *Service) GetMemberById(GroupId int) ([]models.Group, error) {
	rows, err := S.Database.Getmember(GroupId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fmt.Println(rows)

	Group :=make(map[string]int )
	for rows.Next() {
		var groupIDScan int
		if err := rows.Scan(&groupIDScan); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		Group["id"] = groupIDScan

	}

	var groups []models.Group

	for _, v := range Group {
		var group models.Group
		rowGroupe:= S.Database.GetGroupById(v)
		if err := rowGroupe.Scan(utils.GetScanFields(&group)...); err != nil {
			return nil, fmt.Errorf("error scanning group data: %v", err)
		}
		groups = append(groups, group)
	}
	
	return groups,nil
}

