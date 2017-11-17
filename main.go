package icsoc_monitor

import (
	"paas/icsoc-monitor/calc_task"
	"paas/icsoc-monitor/config"
	"paas/icsoc-monitor/mq"
	"paas/icsoc-monitor/server"

	_ "paas/icsoc-monitor/db"
)

func main(){
	config.Init("config.yml")

	go server.StartHttpServer()

	go calc_task.StartTask()

	mq.StartMq()
}