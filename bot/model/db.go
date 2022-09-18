package model

import (
	"fmt"
)

type (
	AdminAccount struct {
		Account       int64 `json:"account"`       // 账号
		ParentAccount int64 `json:"parentAccount"` // 父账号
	}

	AssociationGroup struct {
		GroupID  int64  `bson:"group-id"`         // qq群号
		Operator int64  `bson:"operator-account"` // 添加人
		School   string `bson:"school"`           // 所属学校
	}

	WallMessage struct {
		MsgID  int64  `bson:"msg-id"`
		Msg    string `bson:"msg"`
		School string `bson:"school"`
	}

	NextID struct {
		ID    string `bson:"_id"`
		Value int64  `bson:"value"`
	}

	AuditRecord struct {
		MsgID      int64   `bson:"msg-id"`
		Status     string  `bson:"status"`
		OperatorID int64   `bson:"operatorID"`
		QmID       []int64 `bson:"qm-id"`
	}
)

func (a AdminAccount) String() string {
	return fmt.Sprintf("-account: %d\n\t-parentAccount: %d", a.Account, a.ParentAccount)
}

func (a AssociationGroup) String() string {
	return fmt.Sprintf("-group-id: %d\n\t-creator: %d\n\t-school: %s", a.GroupID, a.Operator, a.School)
}
