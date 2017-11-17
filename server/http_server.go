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
}

func successHandler(context *gin.Context) {
	context.String(http.StatusOK, "success")
}