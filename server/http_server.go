package server

import (
	"paas/icsoc-monitor/config"
	"paas/icsoc-monitor/util"

	"net/http"

	"fmt"

	log "github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
)

func StartHttpServer() {
	conf := config.GetConf()
	if util.IsEmpty(conf.Http.Addr) {
		log.Exit("local host is emtpy")
	}

	log.Info("start http server at:%s", conf.Http.Addr)

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	registerHandler(engine)
	runErr := engine.Run(conf.Http.Addr)
	log.Exit(fmt.Sprintf("StartHttpServer service stop. err: %v", runErr))
}

func registerHandler(en *gin.Engine) {
	en.GET("/ping", successHandler)
	en.POST("/im-cs/statics/work", AgentWorkStatics)
	en.POST("/im-cs/statics/achievements", AgentAchievements)
	en.POST("/im-cs/statics/agent_status_record", AgentStatusRecord)
	en.POST("/im-cs/statics/session_record", AgentSessionRecord)
	en.POST("/im-cs/statics/agent_eval", AgentEval)
	en.POST("/im-cs/statics/source_eval", ChannelEval)
	en.POST("/im-cs/statics/agent_monitor", AgentMonitor)
	en.POST("/im-cs/statics/source_monitor", ChannelMonitor)
}

func successHandler(context *gin.Context) {
	context.String(http.StatusOK, "success")
}