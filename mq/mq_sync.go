package mq

import (
	"paas/icsoc-monitor/config"

	log "github.com/alecthomas/log4go"
)

const (
	exchangeName = "icsoc.imcsmq.ex"
	queueName = "icsoc-icsoc-monitor"
)

func StartMq() {
	conf := config.GetConf()

	routeKeys := GetRouteKeys()

	rmq := NewRmqConsumer(conf.RabbitMQ.Addr, exchangeName, routeKeys)
	for msg := range rmq.MsgRecv {
		log.Info("got rmq message on exchange(%s) routingkey(%s), data(%s)",
			msg.Exchange, msg.RoutingKey, string(msg.Body))
		routeFunc := GetRouteFunc(msg.RoutingKey)
		if routeFunc != nil {
			routeFunc(msg)
		}else {
			msg.Acknowledger.Ack(msg.DeliveryTag, false)
		}
	}
	log.Warn("StartMq end...")
}
