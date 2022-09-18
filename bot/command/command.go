/*
	test
*/
package command

import (
	"QQ-BOT/bot/service"
	"github.com/pkg/errors"
)

type Command struct {
	service *service.Service
}

func (c *Command) Init() {
	c.service = service.Init()
}

func (c Command) CommandAnalysis(command []string, param map[string]interface{}) {
	account := param["user_id"].(float64)
	if c.Auth(int64(account)) == false { // 鉴权失败
		c.ErrorComment(errors.New("鉴权失败"), "", param)
		return
	}
	if len(command) <= 1 {
		c.ErrorComment(errors.New("缺少参数"), "", param)
		return
	}
	switch command[1] {
	case "添加管理员账号":
		c.AddAdminAccount(command, param)
	case "获取管理员账号列表":
		c.GetAdminAccountList(command, param)
	case "删除管理员账号":
		c.DeleteAdminAccount(command, param)
	case "at所有人":
		c.AtALL(command, param)

	case "添加关联群组":
		c.AddAssociationGroup(command, param)
	case "删除关联群组":
		c.DeleteAssociationGroup(command, param)
	case "获取关联群组列表":
		c.GetAssociationGroupList(command, param)
	case "审核":
		if len(command) < 4 {
			c.ErrorComment(errors.New("缺少参数"), "", param)
			return
		}
		c.Audit(command[2], command[3], param)
	case "撤销审核":
		c.CancelAudit(command, param)
	case "发布广告":
		c.PostAD(command, param)
	case "删除广告":
		c.DeleteAD(command, param)
	case "获取所有广告":
		c.GetAllAD(param)
	case "帮助":
		c.Help(param)
	//case "+":
	//	c.QuickAuditAgree(command, param)
	case "-":
		c.QuickAuditDeny(command, param)
	case "@":
		c.AtALL(command, param)
	case "~":
		c.CancelAudit(command, param)
	default:
		c.ErrorComment(errors.New("无效的指令"), "请输入:bot 帮助 \n查看所有命令", param)
		return
	}
	c.Finish(command[1], param)
}

func (c Command) AuditPass(command []string, param map[string]interface{}) {
	account := param["user_id"].(float64)
	if c.Auth(int64(account)) == false { // 鉴权失败
		c.ErrorComment(errors.New("鉴权失败"), "", param)
		return
	}
	c.QuickAuditAgree(command, param)
}
