package mq

import (
	"paas/icsoc-monitor/config"

	log "github.com/alecthomas/log4go"
)

const (
	exchangeName = "icsoc.imcsmq.ex"
	queueName = "icsoc-monitor"
)

func StartMq() {
	conf := config.GetConf()

	routeKeys := GetRouteKeys()

	rmq := NewRmqConsumer(conf.RabbitMQ.Addr, exchangeName, routeKeys)
	for msg := range rmq.MsgRecv {
		log.Info("got rmq message on exchange(%s) routingkey(%s), data(%v)",
			msg.Exchange, msg.RoutingKey, msg.Body)
		GetRouteFunc(msg.RoutingKey)(msg)
	}
}
