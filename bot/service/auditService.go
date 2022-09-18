package service

import (
	"QQ-BOT/bot/db"
	"QQ-BOT/bot/model"
)

type AuditService struct {
}

func (a AuditService) AddAuditRecord(mID int64, status string, operatorID int64, qqMID []int64) error {
	return db.AddAuditRecord(mID, status, operatorID, qqMID)
}

func (a AuditService) AuditRecordIsExist(msgID int64) bool {
	return db.AuditRecordIsExist(msgID)
}

func (a AuditService) GetAuditRecord(mID int64) (model.AuditRecord, error) {
	return db.GetAuditRecord(mID)
}
