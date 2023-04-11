package service

import (
	"QQ-BOT/bot/db"
	"QQ-BOT/bot/model"
	"github.com/pkg/errors"
	"log"
)

type AdminService struct {
}

func (a AdminService) AddAdminAccount(account int64, parentAccount int64) error {
	if err := db.AddAdminAccount(account, parentAccount); err != nil {
		return errors.Wrap(err, "添加管理员失败")
	}
	return nil
}

func (a AdminService) GetAdminAccountList() ([]model.AdminAccount, error) {
	return db.GetAdminAccountList()
}

func (a AdminService) DeleteAdminAccount(account int64) error {
	return db.DeleteAdminAccount(account)
}

func (a AdminService) Auth(account int64) bool {
	if account == 707402933 {
		return true
	}
	flag := false
	list, err := a.GetAdminAccountList()
	if err != nil {
		log.Fatal("获取管理员账号列表出错 " + err.Error())
		return false
	}
	for _, v := range list {
		if v.Account == account {
			flag = true
		}
	}
	return flag
}
