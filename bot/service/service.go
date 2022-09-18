package service

import (
	"QQ-BOT/bot/db"
	"github.com/spf13/viper"
)

type Service struct {
	AdminService      AdminService
	BaseService       BaseService
	GroupService      GroupService
	SchoolWallService SchoolWallService
	AuditService      AuditService
}

var BaseURL string

func Init() *Service {
	db.Init()
	BaseURL = viper.GetString("baseURL")
	return &Service{
		AdminService: AdminService{},
		BaseService:  BaseService{},
		GroupService: GroupService{},
	}
}
