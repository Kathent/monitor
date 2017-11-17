package db

import (
	"time"

	"paas/icsoc-monitor/config"

	"github.com/alecthomas/log4go"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

func LoadSession() {
	conf := config.GetConf()

	s, err := mgo.DialWithTimeout(conf.MongoDb.Addr, time.Second * 2)
	if err != nil {
		log4go.Exit("dial mongodb err:%v", err)
	}

	session = s
}

func GetSession() *mgo.Session {
	return session
}
