package main

import (
	"paas/icsoc-monitor/calc_task"
	"paas/icsoc-monitor/config"
	"paas/icsoc-monitor/db"
	"paas/icsoc-monitor/mq"
	"paas/icsoc-monitor/server"
)

func main(){
	config.Init("config.yml")
	db.Init()

	go server.StartHttpServer()

	go calc_task.StartTask()

	mq.StartMq()
}