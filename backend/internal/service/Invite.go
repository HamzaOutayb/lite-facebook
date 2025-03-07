package service

import (
	"fmt"

	"social-network/internal/models"
)

func (service *Service) CreateInvite(Invites models.Invite) (err error) {
	resA := service.Database.IsFollow(Invites.Sender, Invites.Receiver)
	resB := service.Database.IsFollow(Invites.Receiver, Invites.Sender)
	if resA && resB {

		Invites.Status = "pending"
		id := 0
		id, err = service.Database.Saveinvite(&Invites)
		if err != nil {
			return fmt.Errorf("bad request")
		}
		if id == 0 {
			if Invites.Sender == Invites.Receiver {
				return fmt.Errorf("bad request")
			}
			err = service.Database.SaveInvite(&Invites)
		} else {
			Invites.ID = id
			err = service.Database.DeleteInvites(&Invites)
		}
	}else {
		return fmt.Errorf("not follow")
	}
	return
}

func (service *Service) InviderDecision(Invites *models.Invite) (err error) {
	if Invites.Status == "accepted" {
		err = service.Database.AcceptInviteRequest(Invites)
		fmt.Println(err)
	} else if Invites.Status == "rejected" {
		err = service.Database.DeleteInvites(Invites)
	} else {
		err = fmt.Errorf("bad request")
	}
	return
}
