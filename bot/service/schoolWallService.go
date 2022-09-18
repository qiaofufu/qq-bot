package service

import (
	"QQ-BOT/bot/db"
	"QQ-BOT/bot/model"
)

type SchoolWallService struct {
}

func (s SchoolWallService) PostWallMessage(msg string, school string) (id int64, err error) {
	return db.AddSchoolWallMessage(msg, school)
}

func (s SchoolWallService) GetSchoolWallMessage(msgID int64) (model.WallMessage, error) {
	return db.GetSchoolWallMessage(msgID)
}
