package command

import (
	"QQ-BOT/bot/tootls"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"strconv"
	"time"
)

// AtALL @全体关联群组
// command format: bot at所有人 <内容>
func (c *Command) AtALL(command []string, param map[string]interface{}) {
	group, err := c.service.GroupService.GetAssociationGroupList()
	if err != nil {
		c.ErrorComment(err, "获取关联群组失败", param)
	}
	log.Println(group)
	for _, v := range group {
		c.service.BaseService.AtALL(command[2], v.GroupID)
	}
}

// AddAdminAccount 添加管理员账号
// command format: bot 添加管理员账号 <qq账号>
func (c *Command) AddAdminAccount(command []string, param map[string]interface{}) {
	ac, err := strconv.ParseInt(command[2], 10, 64)
	if err != nil {
		c.ErrorComment(err, "参数解析失败", param)
		return
	}
	parentAccount := int64(param["user_id"].(float64))
	err = c.service.AdminService.AddAdminAccount(ac, parentAccount)
	if err != nil {
		c.ErrorComment(err, "添加管理员账号失败", param)
		return
	}
}

// DeleteAdminAccount 删除管理员账号
// command format: bot 删除管理员账号 <qq账号>
func (c *Command) DeleteAdminAccount(command []string, param map[string]interface{}) {
	account, err := strconv.ParseInt(command[2], 10, 64)
	if err != nil {
		c.ErrorComment(err, "参数解析失败", param)
		return
	}
	userID := int64(param["user_id"].(float64))
	if userID != 707402933 {
		c.ErrorComment(errors.New("权限错误"), "没有操作权限, 联系707402933获取权限", param)
		return
	}
	c.service.AdminService.DeleteAdminAccount(account)
}

// CancelAudit 撤销审核
// command format: bot 撤销审核 <id>
func (c *Command) CancelAudit(command []string, param map[string]interface{}) {
	mid, err := strconv.ParseInt(command[2], 10, 64)
	if err != nil {
		c.ErrorComment(err, "解析id失败", param)
		return
	}
	record, err := c.service.AuditService.GetAuditRecord(mid)
	if err != nil {
		c.ErrorComment(err, "获取审核记录失败", param)
		return
	}
	for _, v := range record.QmID {
		c.service.BaseService.DeleteMessage(v)
	}

}

// Audit 审核
// command format: bot 审核 <id> <状态:通过/拒绝>
func (c *Command) Audit(mID string, status string, param map[string]interface{}) {
	id, err := strconv.ParseInt(mID, 10, 64)
	if err != nil {
		c.ErrorComment(err, "ID解析失败", param)
		return
	}
	if c.service.AuditService.AuditRecordIsExist(id) == true {
		return
	}

	userID := int64(param["user_id"].(float64))

	var msgID []int64

	if status == "通过" {
		msg, err := c.service.SchoolWallService.GetSchoolWallMessage(id)
		message := fmt.Sprintf("%s\n投稿ID: %d", msg.Msg, msg.MsgID)
		if err != nil {
			c.ErrorComment(err, "获取审核信息错误", param)
			return
		}
		group, err := c.service.GroupService.GetAssociationGroupList()
		if err != nil {
			c.ErrorComment(err, "获取关联群组列表失败", param)
			return
		}
		for _, v := range group {
			if v.School == msg.School {
				t := c.service.BaseService.SendGroupMessage(message, v.GroupID, true)
				msgID = append(msgID, t)
			}
		}
	}

	err = c.service.AuditService.AddAuditRecord(id, status, userID, msgID)
	if err != nil {
		c.ErrorComment(err, "插入审核记录错误", param)
		return
	}
	groupID := int64(param["group_id"].(float64))
	c.service.BaseService.SendGroupMessage(fmt.Sprintf("消息ID: %s\n审核人: %d\n审核结果: %s", mID, userID, status), groupID, false)
}

func (c *Command) QuickAuditAgree(command []string, param map[string]interface{}) {
	idArr := command[1:]
	log.Println(idArr)
	for _, v := range idArr {
		c.Audit(v, "通过", param)
	}
}

func (c *Command) QuickAuditDeny(command []string, param map[string]interface{}) {
	idArr := command[2:]
	for _, v := range idArr {
		c.Audit(v, "拒绝", param)
	}
}

func (c *Command) ErrorComment(err error, msg string, param map[string]interface{}) {
	groupID := int64(param["group_id"].(float64))
	userID := int64(param["user_id"].(float64))
	c.service.BaseService.SendGroupMessage(fmt.Sprintf("[ERROR] 用户[%d]\n%s\nerr:%s", userID, msg, err), groupID, false)
}

func (c Command) Finish(action string, param map[string]interface{}) {
	groupID := int64(param["group_id"].(float64))
	userID := int64(param["user_id"].(float64))
	c.service.BaseService.SendGroupMessage(fmt.Sprintf("[SUCCESS]\n用户: [CQ:at,qq=%d]\n动作: %s", userID, action), groupID, false)
}

// AddAssociationGroup 添加关联群组
// command format: bot 添加关联群组 <群号> <所属学校>
// eg: bot 添加关联群组 707402933 四川轻化工大学
func (c *Command) AddAssociationGroup(command []string, param map[string]interface{}) {
	groupID, err := strconv.ParseInt(command[2], 10, 64)
	if err != nil {
		c.ErrorComment(err, "解析groupID错误", param)
		return
	}
	if len(command) < 4 {
		c.ErrorComment(errors.New("参数错误"), "缺少学校字段", param)
		return
	}
	school := command[3]
	operatorAccount := int64(param["user_id"].(float64))
	c.service.GroupService.AddAssociationGroup(groupID, operatorAccount, school)
}

// DeleteAssociationGroup 删除关联群组
// command format: bot 删除关联群组 <群号>
// eg: bot 删除关联群组 707402933
func (c *Command) DeleteAssociationGroup(command []string, param map[string]interface{}) {
	groupID, err := strconv.ParseInt(command[2], 10, 64)
	if err != nil {
		c.ErrorComment(err, "解析groupID错误", param)
		return
	}
	err = c.service.GroupService.RemoveAssociationGroup(groupID)
	if err != nil {
		c.ErrorComment(err, "移除群组失败", param)
		return
	}
}

// GetAssociationGroupList
// command format: bot 获取关联群组列表
func (c *Command) GetAssociationGroupList(command []string, param map[string]interface{}) {
	list, err := c.service.GroupService.GetAssociationGroupList()
	groupID := int64(param["group_id"].(float64))
	if err != nil {
		c.ErrorComment(err, "获取失败", param)
		return
	}
	str := "<===关联群组列表===>\n"
	for _, v := range list {
		str += v.String() + "\n"
	}
	c.service.BaseService.SendGroupMessage(str, groupID, false)
}

func (c *Command) Auth(account int64) bool {
	flag := c.service.AdminService.Auth(account)
	return flag
}

// GetAdminAccountList 获取管理员账号列表
// command format: bot 获取管理员账号列表
func (c *Command) GetAdminAccountList(command []string, param map[string]interface{}) {
	list, err := c.service.AdminService.GetAdminAccountList()
	groupID := int64(param["group_id"].(float64))
	if err != nil {
		c.ErrorComment(err, "获取失败", param)
		return
	}
	str := "<===管理员账号列表===>\n"
	for _, v := range list {
		str += v.String() + "\n"
	}
	c.service.BaseService.SendGroupMessage(str, groupID, false)
}

// Help 帮助
// command format: bot 帮助
func (c *Command) Help(param map[string]interface{}) {
	groupID := int64(param["group_id"].(float64))
	c.service.BaseService.SendGroupMessage(getHelp(), groupID, false)
}

func getHelp() string {

	menu := "" +
		"@全体关联群组\n" +
		"command format: bot at所有人 <内容>\n\n" +
		"添加管理员账号\n" +
		"command format: bot 添加管理员账号 <qq账号>\n\n" +
		"删除管理员账号\n" +
		"command format: bot 删除管理员账号 <qq账号>\n\n" +
		"审核\n" +
		"command format: bot 审核 <id> <状态:通过/拒绝>\n\n" +
		"添加关联群组\n" +
		"command format: bot 添加关联群组 <群号> <所属学校>\n\n" +
		"删除关联群组\n" +
		"command format: bot 删除关联群组 <群号>\n\n" +
		"获取关联群组列表\n" +
		"command format: bot 获取关联群组列表\n\n" +
		"获取管理员账号列表\n" +
		"command format: bot 获取管理员账号列表\n\n" +
		"发布广告\n" +
		"command format: bot 发布广告 <内容> [广告间隔(秒)] [重复次数]\n\n" +
		"删除广告\n" +
		"command format: bot 删除广告 <广告id>\n\n" +
		"获取所有广告\n" +
		"command format: bot 获取所有广告\n\n" +
		"帮助\n" +
		"command format: bot 帮助\n\n" +
		"快捷审核通过\n" +
		"command format: bot + <id1> <id2> ....\n\n" +
		"快捷审核拒绝\n" +
		"command format: bot - <id1> <id2> ....\n\n" +
		"快捷@全体\n" +
		"command format: bot @ <内容>\n\n" +
		"快捷撤销审核\n" +
		"command format: bot ~ <id>\n\n"

	return "Help\n" + menu
}

// PostAD 发布广告
// command format: bot 发布广告 <内容> [广告间隔(秒)] [重复次数]
func (c *Command) PostAD(command []string, param map[string]interface{}) {
	ad := command[2]
	second := 60 * 60 * 3
	count := 8
	if len(command) == 5 {
		if command[3] != "0" {
			s, err := strconv.Atoi(command[3])
			if err != nil {
				c.ErrorComment(err, "解析时间间隔出错", param)
				return
			}
			second = s
		}
		if command[4] != "0" {
			co, err := strconv.Atoi(command[4])
			if err != nil {
				c.ErrorComment(err, "解析次数出错", param)
				return
			}
			count = co
		}
	}
	log.Println(second)
	log.Println(count)
	tootls.AddCornFunc(time.Second*time.Duration(second), count, func() {
		group, err := c.service.GroupService.GetAssociationGroupList()
		if err != nil {
			log.Println("获取关联群组失败")
			return
		}
		vis := map[int64]bool{}
		for _, v := range group {
			if vis[v.GroupID] != true {
				vis[v.GroupID] = true
				c.service.BaseService.SendGroupMessage(ad, v.GroupID, false)
			}
		}
	}, ad)
}

// DeleteAD 删除广告
// command format: bot 删除广告 <广告id>
func (c *Command) DeleteAD(command []string, param map[string]interface{}) {
	adID, err := strconv.Atoi(command[2])
	if err != nil {
		c.ErrorComment(err, "解析广告ID错误", param)
		return
	}
	err = tootls.StopCorn(adID)
	if err != nil {
		c.ErrorComment(err, "删除广告错误", param)
		return
	}
}

// GetAllAD 获取所有广告
// command format: bot 获取所有广告
func (c *Command) GetAllAD(param map[string]interface{}) {
	data := tootls.GetTaskDataList()
	groupID := int64(param["group_id"].(float64))
	str := "<===广告信息===>\n"
	for key, value := range data {
		str += fmt.Sprintf("广告ID: %d\n\t内容: %s\n", key, value)
	}
	c.service.BaseService.SendGroupMessage(str, groupID, false)
}
