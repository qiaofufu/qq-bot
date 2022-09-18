package db

import (
	"QQ-BOT/bot/model"
	"log"
)

var dbs []Database

var device = make(map[string]func() Database)

func RegisterDevice(name string, initFunc func() Database) {
	device[name] = initFunc
}

func Init() {
	for name, init := range device {
		log.Println(name)
		db := init()
		if db != nil {
			dbs = append(dbs, db)
		}
	}
	Open()
}

type Database interface {
	Open()

	AddAdminAccount(account int64, parentAccount int64) error
	GetAdminAccountList() ([]model.AdminAccount, error)
	DeleteAdminAccount(account int64) error

	GetAssociationGroupList() ([]model.AssociationGroup, error)
	AddAssociationGroup(groupID int64, operatorAccount int64, school string) error
	RemoveAssociationGroup(groupID int64) error

	AddSchoolWallMessage(msg string, school string) (msgID int64, err error)
	GetSchoolWallMessage(msgID int64) (model.WallMessage, error)

	AddAuditRecord(mID int64, status string, operatorID int64, qqMID []int64) error
	AuditRecordIsExist(msgID int64) bool
	GetAuditRecord(mID int64) (model.AuditRecord, error)
}

func Open() {
	for _, v := range dbs {
		v.Open()
	}
}

func GetAuditRecord(mID int64) (model.AuditRecord, error) {
	if len(dbs) == 0 {
		panic("database not connect.")
	}
	return dbs[0].GetAuditRecord(mID)
}

func AddAuditRecord(mID int64, status string, operatorID int64, qqMID []int64) error {
	if len(dbs) == 0 {
		panic("database not connect.")
	}
	for _, v := range dbs {
		err := v.AddAuditRecord(mID, status, operatorID, qqMID)
		if err != nil {
			return err
		}
	}
	return nil
}

func AuditRecordIsExist(msgID int64) bool {
	if len(dbs) == 0 {
		panic("database not connect.")
	}
	return dbs[0].AuditRecordIsExist(msgID)
}

func AddSchoolWallMessage(msg string, school string) (msgID int64, err error) {
	if len(dbs) == 0 {
		panic("database not connect.")
	}
	for _, v := range dbs {
		msgID, err = v.AddSchoolWallMessage(msg, school)
		if err != nil {
			return
		}
	}
	return
}

func GetSchoolWallMessage(msgID int64) (model.WallMessage, error) {
	if len(dbs) == 0 {
		panic("database not connect.")
	}
	return dbs[0].GetSchoolWallMessage(msgID)
}

func AddAdminAccount(account int64, parentAccount int64) error {
	if len(dbs) == 0 {
		panic("database not connect.")
	}
	for _, v := range dbs {
		err := v.AddAdminAccount(account, parentAccount)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAdminAccountList() ([]model.AdminAccount, error) {
	if len(dbs) == 0 {
		panic("database not connect.")
	}
	return dbs[0].GetAdminAccountList()
}

func DeleteAdminAccount(account int64) error {
	if len(dbs) == 0 {
		panic("database not connect.")
	}
	for _, v := range dbs {
		err := v.DeleteAdminAccount(account)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetAssociationGroupList() ([]model.AssociationGroup, error) {
	if len(dbs) == 0 {
		panic("database not connect.")
	}
	return dbs[0].GetAssociationGroupList()
}

func AddAssociationGroup(groupID int64, operatorAccount int64, school string) error {
	if len(dbs) == 0 {
		panic("database not connect.")
	}
	for _, v := range dbs {
		err := v.AddAssociationGroup(groupID, operatorAccount, school)
		if err != nil {
			return err
		}
	}
	return nil
}

func RemoveAssociationGroup(groupID int64) error {
	if len(dbs) == 0 {
		panic("database not connect.")
	}
	for _, v := range dbs {
		err := v.RemoveAssociationGroup(groupID)
		if err != nil {
			return err
		}
	}
	return nil
}
