package mq

import (
	"fmt"
	"net"
	"time"

	logger "github.com/alecthomas/log4go"
	"github.com/streadway/amqp"
)

const (
	vhost = "im"
)

type RmqConsumer struct {
	core    *amqp.Connection
	mqchan  *amqp.Channel
	connerr chan *amqp.Error

	uri       string
	exchange string
	routingKeys []string
	MsgRecv chan amqp.Delivery
}

func NewRmqConsumer(uri string, exchange string, routingkeys []string) (rmq *RmqConsumer) {
	rmq = &RmqConsumer{
		uri:       uri,
		exchange: exchange,
		routingKeys: routingkeys,
		MsgRecv:   make(chan amqp.Delivery, 256),
	}

	go rmq.mainRoutine()
	return
}

func (rmq *RmqConsumer) mainRoutine() {
	var deliverychan <-chan amqp.Delivery
	for {
		if rmq.core == nil {
			if err := rmq.initCore(); err != nil {
				logger.Debug(fmt.Sprintf("initCore err:%s", err))
			} else if err = rmq.setupContext(); err != nil {
				logger.Debug(fmt.Sprintf("setupContext err:%s", err))
			} else if deliverychan, err = rmq.mqchan.Consume(
				queueName,    // queue
				"",           // consumer
				false,         // auto-ack
				false,        // exclusive
				false,        // no local
				false,        // no wait
				amqp.Table{}, // args
			); err != nil {
				logger.Debug(fmt.Sprintf("consumer(%p) failed consuming queue(%s), err(%v)", rmq, queueName, err))
				rmq.closeConn()
			}
		} else {
			select {
			case err := <-rmq.connerr:
				logger.Error(fmt.Sprintf("consumer(%p) lost connection lost with rabbitmq, err(%v)", rmq, err))
				rmq.closeConn()
			case msg := <-deliverychan: // XXX for range channel会在channel中没有数据时阻塞
				logger.Debug(fmt.Sprintf("got rmq message on exchange(%s) routingkey(%s), data(%v)", msg.Exchange, msg.RoutingKey, msg.Body))
				rmq.MsgRecv <- msg
			}
		}
	}
}

func (rmq *RmqConsumer) closeConn() {
	logger.Debug(fmt.Sprintf("consumer(%p) closing connection", rmq))
	if rmq.core != nil {
		rmq.core.Close()
		rmq.core = nil
	}
	rmq.mqchan = nil
	rmq.connerr = nil
	<-time.After(3 * time.Second)
}

func (rmq *RmqConsumer) initCore() (err error) {
	if rmq.core, err = amqp.DialConfig(rmq.uri, amqp.Config{
		Vhost: vhost,
		Heartbeat: 3 * time.Second,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 2*time.Second)
		},
	}); err != nil {
		logger.Debug(fmt.Sprintf("initCore dial config err:%v", err))
	} else {
		rmq.connerr = make(chan *amqp.Error)
		rmq.core.NotifyClose(rmq.connerr)
	}

	if err != nil {
		rmq.closeConn()
	}
	logger.Info(fmt.Sprintf("consumer(%p) got connection(%p) with rabbitmq(%s) err(%v)", rmq, rmq.core, rmq.uri, err))
	return
}

func (rmq *RmqConsumer) setupContext() (err error) {
	//XXX the creation of channel and queue MUST be in same goroutine with `Consume`, otherwise the channel will disappear
	//XXX block or not, this is the only choice
	if rmq.mqchan, err = rmq.core.Channel(); err != nil {
	} else if err = rmq.mqchan.ExchangeDeclare(
		rmq.exchange, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		logger.Error(fmt.Sprintf("consumer(%p) failed declaring exchange(%s), err(%v)", rmq, rmq.exchange, err))
	} else if err = rmq.setupQueue(rmq.mqchan, rmq.exchange); err != nil {
		// error already log in setupQueues function
	} else if err = rmq.mqchan.Qos(1, 0, false); err != nil {
		logger.Error(fmt.Sprintf("consumer(%p) failed setting Qos, err(%v)", rmq, err))
	}

	if err != nil {
		rmq.closeConn()
	}
	return
}

func (rmq *RmqConsumer) setupQueue(mqchan *amqp.Channel, exchange string) (err error) {
	//XXX the creation of channel and queue MUST be in same goroutine with `Consume`, otherwise the channel will disappear
	//XXX block or not, this is the only choice
	if consumerqueue, err := mqchan.QueueDeclare(queueName, false, false, false, false, amqp.Table{"x-expires": int32(3000)}); err != nil {
		logger.Error("consumer(%p) failed declaring queue(%s) on exchange(%s), err(%v)", rmq, "im-cs-consumer", exchange, err)
	} else {
		for _, routingkey := range rmq.routingKeys {
			if err = mqchan.QueueBind(consumerqueue.Name, routingkey, exchange, false, nil); err != nil {
				logger.Error(fmt.Sprintf("consumer(%p) failed binding routingkey(%s) on queue(%s) in exchange(%s), err(%v)", rmq, routingkey, consumerqueue.Name, exchange, err))
				break
			} else {
				logger.Debug(fmt.Sprintf("consumer(%p) queue(%s) declared on exchange(%s) binding routingkey(%s)", rmq, consumerqueue.Name, exchange, routingkey))
			}
		}
	}
	return
}
