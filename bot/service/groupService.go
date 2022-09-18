package service

import (
	"QQ-BOT/bot/db"
	"QQ-BOT/bot/model"
)

type GroupService struct {
}

func (g GroupService) AddAssociationGroup(groupID int64, operatorAccount int64, school string) error {
	return db.AddAssociationGroup(groupID, operatorAccount, school)
}

func (g GroupService) GetAssociationGroupList() ([]model.AssociationGroup, error) {
	return db.GetAssociationGroupList()
}

func (g GroupService) RemoveAssociationGroup(groupID int64) error {
	return db.RemoveAssociationGroup(groupID)
}
