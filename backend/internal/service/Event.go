package service

import (
	"database/sql"
	"fmt"
	"log"

	"social-network/internal/models"
	utils "social-network/pkg"
)

func (service *Service) CreateEvent(Events models.Event) (err error) {
	// valid, err := service.Database.GetCreatorGroup(Events.GroupID, Events.UserID)
	err2 := service.VerifyGroup(Events.GroupID, Events.UserID)
	if err2.Err != nil {
		fmt.Println("err => ", err)
		err = err2.Err
		return
	}
	err = service.Database.SaveEvent(&Events)
	return
}

func (service *Service) AllEvents(Event models.Event) ([]models.Event, error) {
	rows, err := service.Database.GetallEvents(Event.GroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var event models.Event
		if err := rows.Scan(utils.GetScanFields(&event)...); err != nil {
			fmt.Println(err)
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (S *Service) GetEventsById(Event *models.Event) (*models.Event, error) {
	row := S.Database.GetEventById(Event.ID)
	if row == nil {
		return nil, fmt.Errorf("no group found with ID: %s", Event.ID)
	}

	// Scan the row into the Event struct
	if err := row.Scan(utils.GetScanFields(Event)...); err != nil {
		return nil, fmt.Errorf("error scanning Event data: %v", err)
	}

	return Event, nil
}

func (S *Service) PostEventsOption(OptionEvent models.EventOption) (err error) {
	fmt.Println("OptionEvent", OptionEvent)
	booll, err := S.Database.CheckEvent(OptionEvent.EventID, OptionEvent.UserID)
	fmt.Println("booll", booll)
	if err != nil {
		if err == sql.ErrNoRows {
			err = S.Database.SaveOptionEvent(&OptionEvent)
			return
		}
	}
	if booll == OptionEvent.Going {
		return
	} else if booll != OptionEvent.Going {
		err = S.Database.UpdateOptionEvent(OptionEvent)
		fmt.Println(err)
		return
	}
	err = S.Database.SaveOptionEvent(&OptionEvent)
	return
}

func (S *Service) GetEventsOption(OptionEvent models.EventOption) ([]models.EventOption, error) {
	rows, err := S.Database.OptionEvent(OptionEvent.EventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.EventOption

	for rows.Next() {
		var event models.EventOption
		if err := rows.Scan(utils.GetScanFields(&event)...); err != nil {
			fmt.Println(err)
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (S *Service) GetEventgoing(OptionEvent models.EventOption, user_id int) (int, string, error) {
	rows, user, err := S.Database.ChoiseEvent(OptionEvent.EventID, OptionEvent.Going)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("rows", rows)
			fmt.Println("user", user)
			fmt.Println("err", nil)
			return 0, "not", nil
		} else {
			return 0, "", err
		}
	}
	if user == user_id {
		fmt.Println("user", user)
		return rows, "action", nil
	}
	return rows, "not", nil
}
