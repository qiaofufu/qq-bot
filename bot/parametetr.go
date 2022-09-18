package bot

import (
	"strings"
)

type Parameter struct {
}

func (p *Parameter) ParameterParse(parameter map[string]interface{}) {
	switch parameter["post_type"] {
	case "message":
		p.messageParse(parameter)
	case "request":
	case "notice":
	case "meta_event":
	}
}

func (p Parameter) messageParse(param map[string]interface{}) {
	switch param["message_type"] {
	case "group":
		p.messageGroupParse(param)
	case "private":
	}
}

func (p Parameter) messageGroupParse(param map[string]interface{}) {
	content := param["message"].(string)
	if strings.HasPrefix(content, "bot") { // 机器人指令
		comm := strings.Fields(content)
		bot.Command.CommandAnalysis(comm, param)
	} else if strings.HasPrefix(content, "通过") {
		comm := strings.Fields(content)
		bot.Command.AuditPass(comm, param)
	} else { // 非机器人指令
	}
}
