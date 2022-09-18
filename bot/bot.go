package bot

import (
	"QQ-BOT/bot/command"
	"QQ-BOT/bot/config"
	"QQ-BOT/bot/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
)

type Bot struct {
	Config  *config.Config
	Command *command.Command

	Parameter *Parameter
}

var bot *Bot

func Default() *Bot {
	return &Bot{Config: &config.Config{}, Command: &command.Command{}, Parameter: &Parameter{}}
}

func Run(b *Bot) error {
	b.Config.Init()
	b.Command.Init()
	bot = b
	r := gin.Default()
	r.Use(Cors())
	gin.SetMode(gin.ReleaseMode)

	port := ":" + viper.GetString("server.port")
	httpsPort := ":" + viper.GetString("server.httpsPort")
	pemPath := viper.GetString("server.pemPath")
	keyPath := viper.GetString("server.keyPath")
	r.POST("", bot.ReportCall)
	r.POST("qqMessage", bot.WallMessage)
	go r.Run(port)
	r.RunTLS(httpsPort, pemPath, keyPath)
	return nil
}

func (b *Bot) ReportCall(ctx *gin.Context) {
	// 读取body
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("read body error")
		return
	}
	var parameter map[string]interface{}

	json.Unmarshal(body, &parameter)

	b.Parameter.ParameterParse(parameter)
}

func (b Bot) WallMessage(ctx *gin.Context) {
	type wallMessage struct {
		Content string   `json:"content"`
		School  string   `json:"school"`
		Images  []string `json:"images"`
	}
	var dto wallMessage
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "数据绑定失败",
		})
		return
	}
	// 加入待审核信息
	baseMsg := dto.Content
	for _, v := range dto.Images {
		baseMsg += fmt.Sprintf("[CQ:image,file=imagefile,subType=0,url=%s]", v)
	}
	msgID, err := service.SchoolWallService{}.PostWallMessage(baseMsg, dto.School)
	if err != nil {
		log.Println("加入审核信息错误, 错误码0100")
		return
	}

	// 发送审核群
	msg := fmt.Sprintf("表白墙投稿(待审核) ID:[%d]\n学校:%s\n内容:\n%s", msgID, dto.School, baseMsg)
	service.BaseService{}.SendGroupMessage(msg, viper.GetInt64("auditGroup"), false)

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "投稿成功",
	})
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
