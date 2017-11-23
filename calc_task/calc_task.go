package calc_task

import (
	"fmt"
	"time"

	"paas/icsoc-monitor/constants"
	"paas/icsoc-monitor/db"
	"paas/icsoc-monitor/mq"
	"paas/icsoc-monitor/util"

	"github.com/alecthomas/log4go"
	"gopkg.in/mgo.v2/bson"
)

func StartTask(){
	go timerWork()
	redisToMongoDb()
}

func timerWork() {
	log4go.Info("timerWork enter.")

	now := time.Now()
	duration := util.TomorrowDuration(now)
	time.AfterFunc(duration, timerWork)

	transTime := mq.TransTime(now.AddDate(0,0,-1), constants.DATE_FORMATE)
	nextTransTime := mq.TransTime(now, constants.DATE_FORMATE)
	//迁移坐席监控表
	keyPattern := fmt.Sprintf(constants.AGENT_MONITOR_HASH_KEY, transTime, "*", "*")
	var cursor uint64
	var count int64 = 20
	var keys []string
	var err error

	pipe := db.GetClient().Pipeline()

	for {
		keys, cursor, err = db.GetClient().Scan(cursor, keyPattern, count).Result()
		if err != nil {
			log4go.Warn("syncRedisToMongoDb scan err:%v", err)
			break
		}

		if len(keys) <= 0 {
			break
		}

		for _, v := range keys {
			resMap, err := db.GetClient().HGetAll(v).Result()
			if err != nil {
				log4go.Warn("timerWork HGetAll err:%v", err)
				continue
			}

			vccId := util.GetString(resMap[constants.AGENT_MONITOR_FIELD_VCCID])
			agentId := util.GetString(resMap[constants.AGENT_MONITOR_FIELD_AGENTID])
			newKey := fmt.Sprintf(constants.AGENT_MONITOR_HASH_KEY, nextTransTime, vccId, agentId)

			exist, _ := db.GetClient().HExists(newKey, constants.AGENT_MONITOR_FIELD_AGENTID).Result()
			if exist {//如果存在了就不再处理
				continue
			}
			statusStartTime := int64(util.GetInt(resMap[constants.AGENT_MONITOR_FIELD_STATUS_START_TIME]))
			newMap := map[string]interface{}{
				constants.AGENT_MONITOR_FIELD_VCCID:
					util.GetString(resMap[constants.AGENT_MONITOR_FIELD_VCCID]),
				constants.AGENT_MONITOR_FIELD_AGENTID:
					util.GetString(resMap[constants.AGENT_MONITOR_FIELD_AGENTID]),
				constants.AGENT_MONITOR_FIELD_WORKER_ID:
					util.GetString(resMap[constants.AGENT_MONITOR_FIELD_WORKER_ID]),
				constants.AGENT_MONITOR_FIELD_NAME:
					util.GetString(resMap[constants.AGENT_MONITOR_FIELD_NAME]),
				constants.AGENT_MONITOR_FIELD_DEP_ID:
					util.GetString(resMap[constants.AGENT_MONITOR_FIELD_DEP_ID]),
				constants.AGENT_MONITOR_FIELD_GROUP_IDS:
					util.GetString(resMap[constants.AGENT_MONITOR_FIELD_GROUP_IDS]),
				constants.AGENT_MONITOR_FIELD_MAX_SESSION_NUM:
					util.GetString(resMap[constants.AGENT_MONITOR_FIELD_MAX_SESSION_NUM]),
				constants.AGENT_MONITOR_FIELD_STATUS:
					util.GetString(resMap[constants.AGENT_MONITOR_FIELD_STATUS]),
				constants.AGENT_MONITOR_FIELD_STATUS_START_TIME:
					util.GetIntString(now.Unix()),
				constants.AGENT_MONITOR_FIELD_CUR_SESSION_NUM:
					util.GetString(resMap[constants.AGENT_MONITOR_FIELD_CUR_SESSION_NUM]),
			}
			pipe.HMSet(newKey, newMap)
			pipe.HIncrBy(v, constants.AGENT_MONITOR_FIELD_ONLINE_TIME_TOTAL, statusStartTime - now.Unix())
		}
		pipe.Exec()

		if cursor <= 0 {
			break
		}
	}


	//迁移渠道监控表
	keyPattern = fmt.Sprintf(constants.CHANNEL_MONITOR_HASH_KEY, transTime, "*", "*")
	for {
		keys, cursor, err = db.GetClient().Scan(cursor, keyPattern, count).Result()
		if err != nil {
			log4go.Warn("syncRedisToMongoDb scan err:%v", err)
			break
		}

		if len(keys) <= 0 {
			break
		}

		for _, v := range keys {
			resMap, err := db.GetClient().HGetAll(v).Result()
			if err != nil {
				log4go.Warn("timerWork HGetAll err:%v", err)
				continue
			}

			vccId := util.GetString(resMap[constants.CHANNEL_MONITOR_FIELD_VCCID])
			channelId := util.GetString(resMap[constants.CHANNEL_MONITOR_FIELD_CHANNEL_ID])
			newKey := fmt.Sprintf(constants.CHANNEL_MONITOR_HASH_KEY, nextTransTime, vccId, channelId)
			newMap := map[string]interface{}{
				constants.CHANNEL_MONITOR_FIELD_VCCID:
				util.GetString(resMap[constants.CHANNEL_MONITOR_FIELD_VCCID]),
				constants.CHANNEL_MONITOR_FIELD_CHANNEL_ID:
				util.GetString(resMap[constants.CHANNEL_MONITOR_FIELD_CHANNEL_ID]),
				constants.CHANNEL_MONITOR_FIELD_SOURCE_NAME:
				util.GetString(resMap[constants.CHANNEL_MONITOR_FIELD_SOURCE_NAME]),
				constants.CHANNEL_MONITOR_FIELD_SOURCE_TYPE:
				util.GetString(resMap[constants.CHANNEL_MONITOR_FIELD_SOURCE_TYPE]),
				constants.CHANNEL_MONITOR_FIELD_SESSION_NUM:
				util.GetString(resMap[constants.CHANNEL_MONITOR_FIELD_SESSION_NUM]),
				constants.CHANNEL_MONITOR_FIELD_QUEUE_NUM:
				util.GetString(resMap[constants.CHANNEL_MONITOR_FIELD_QUEUE_NUM]),
			}
			pipe.HMSet(newKey, newMap)
		}
		pipe.Exec()

		if cursor <= 0 {
			break
		}
	}

	//将昨天数据拷过来
	log4go.Info("timerWork over. next duration:%d", duration)
}

func redisToMongoDb() {
	t := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-t.C:
			syncRedisToMongoDb()
		}
	}
}

func syncRedisToMongoDb() {
	now := time.Now()
	tranTime := mq.TransTime(now, constants.DATE_FORMATE)
	keyPattern := fmt.Sprintf(constants.AGENT_MONITOR_HASH_KEY, tranTime, "*", "*")
	var cursor uint64
	var count int64 = 20
	var keys []string
	var err error

	agentWorkC := db.GetSession().DB("").C(constants.STATICS_AW_TABLE_NAME)
	agentEvalC := db.GetSession().DB("").C(constants.STATICS_AE_TABLE_NAME)
	channelEvalC := db.GetSession().DB("").C(constants.STATICS_CE_TABLE_NAME)
	for {
		keys, cursor, err = db.GetClient().Scan(cursor, keyPattern, count).Result()
		if err != nil {
			log4go.Warn("syncRedisToMongoDb scan err:%v", err)
			break
		}

		if len(keys) <= 0 {
			break
		}

		for _, v := range keys{
			strings, err := db.GetClient().HGetAll(v).Result()
			if err != nil {
				log4go.Warn("syncRedisToMongoDb scan err:%v, key:%s", err, v)
				continue
			}

			aw := mq.AgentWorkMG{}
			aw.VccId = util.GetString(strings[constants.AGENT_MONITOR_FIELD_VCCID])
			aw.AgentId = util.GetString(strings[constants.AGENT_MONITOR_FIELD_AGENTID])
			aw.WorkerId = util.GetString(strings[constants.AGENT_MONITOR_FIELD_WORKER_ID])
			aw.InSessionNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_INDEP_SESSION_NUM])
			aw.ReceiveEvalTimes = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_RECEIVE_EVAL_TIMES])

			aw.EvalRatio = float64(aw.ReceiveEvalTimes) / float64(util.GetDefaultInt(aw.InSessionNum, 1))
			aw.ReplyNewsNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_REPLY_MSG_TIMES])
			aw.ReceiveMsgTimes = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_RECEIVE_MSG_TIMES])
			aw.ArRatio = float64(aw.ReceiveMsgTimes) / float64(util.GetDefaultInt(aw.ReplyNewsNum, 1))
			aw.TotalSessionNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_TOTAL_SESSION_NUM])
			aw.SessionKeepSecs = int64(util.GetInt(strings[constants.AGENT_MONITOR_FIELD_SESSION_TIME_TOTAL]))
			aw.AvgSessionSecs = aw.SessionKeepSecs / int64(util.GetDefaultInt(aw.TotalSessionNum, 1))
			aw.FirstRespSecs = int64(util.GetInt(strings[constants.AGENT_MONITOR_FIELD_FIRST_RESP_TIME]))
			aw.AvgResponseSecs = aw.FirstRespSecs / int64(util.GetDefaultInt(aw.TotalSessionNum, 1))
			aw.Date = tranTime
			aw.DeptId = util.GetString(strings[constants.AGENT_MONITOR_FIELD_DEP_ID])
			aw.BusySecs = int64(util.GetInt(strings[constants.AGENT_MONITOR_FIELD_BUSY_TIME_TOTAL]))
			aw.GroupId = util.GetString(strings[constants.AGENT_MONITOR_FIELD_GROUP_IDS])
			aw.InvalidSessionNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_INVALID_SESSION_NUM])
			aw.InvalidRatio = float64(aw.InvalidSessionNum) / float64(util.GetDefaultInt(aw.TotalSessionNum, 1))
			aw.OneServeClientNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_ONE_SERV_CLIENT_NUM])
			aw.OnlineSecs = int64(util.GetInt(strings[constants.AGENT_MONITOR_FIELD_ONLINE_TIME_TOTAL]))
			aw.ServeUserNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_SERV_USER_NUM])
			aw.TransInSessionNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_TRANSFER_IN_TIMES])
			aw.TransOutSessionNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_TRANSFER_OUT_TIMES])
			aw.ReceiveSessionNUm = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_DEP_SESSION_NUM])
			aw.RequireEvalTimes = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_REQUIRE_EVAL_TIMES])

			agentWorkC.Upsert(bson.M{"vcc_id": aw.VccId,
			                         "date": tranTime,
			                         "ag_id": aw.AgentId}, aw)


			ae := mq.AgentEval{}
			ae.VccId = util.GetString(strings[constants.AGENT_MONITOR_FIELD_VCCID])
			ae.AgentId = util.GetString(strings[constants.AGENT_MONITOR_FIELD_AGENTID])
			ae.DepId = util.GetString(strings[constants.AGENT_MONITOR_FIELD_DEP_ID])
			ae.GroupId = util.GetString(strings[constants.AGENT_MONITOR_FIELD_GROUP_IDS])
			ae.WorkerId = util.GetString(strings[constants.AGENT_MONITOR_FIELD_WORKER_ID])
			ae.IndepSessionNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_INDEP_SESSION_NUM])
			ae.RequireEvalTimes = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_REQUIRE_EVAL_TIMES])
			ae.EvaluateNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_RECEIVE_EVAL_TIMES])
			ae.RateEvaluate = float64(ae.EvaluateNum) / float64(util.GetDefaultInt(ae.IndepSessionNum, 1))
			ae.OneStarNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_ONE_STAR_NUM])
			ae.TwoStarNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_TWO_STAR_NUM])
			ae.ThreeStarNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_THREE_STAR_NUM])
			ae.FourStarNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_FOUR_STAR_NUM])
			ae.FiveStarNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_FIVE_STAR_NUM])
			ae.SixStarNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_SIX_STAR_NUM])
			ae.Date = tranTime
			agentEvalC.Upsert(bson.M{"vcc_id": aw.VccId,
				"date": tranTime,
				"ag_id": aw.AgentId}, ae)
		}

		if cursor <= 0 {
			break
		}
	}

	cursor = 0
	keyPattern = fmt.Sprintf(constants.CHANNEL_MONITOR_HASH_KEY, tranTime, "*", "*")
	for {
		keys, cursor, err = db.GetClient().Scan(cursor, keyPattern, count).Result()
		if err != nil {
			log4go.Warn("syncRedisToMongoDb scan err:%v", err)
			break
		}

		if len(keys) <= 0 {
			break
		}

		for _, v := range keys {
			strings, err := db.GetClient().HGetAll(v).Result()
			if err != nil {
				log4go.Warn("syncRedisToMongoDb scan err:%v, key:%s", err, v)
				continue
			}

			ce := mq.ChannelEval{}
			ce.Date = tranTime
			ce.ChannelId = util.GetString(strings[constants.CHANNEL_MONITOR_FIELD_CHANNEL_ID])
			ce.VccId = util.GetString(strings[constants.CHANNEL_MONITOR_FIELD_VCCID])
			ce.SourceName = util.GetString(strings[constants.CHANNEL_MONITOR_FIELD_SOURCE_NAME])
			ce.SourceType = util.GetString(strings[constants.CHANNEL_MONITOR_FIELD_SOURCE_TYPE])
			ce.EndSessionNum = util.GetInt(strings[constants.CHANNEL_MONITOR_FIELD_END_SESSION_NUM])
			ce.EvaluateNum = util.GetInt(strings[constants.CHANNEL_MONITOR_FIELD_EVALUATE_NUM])
			ce.RateEvaluate = float64(ce.EvaluateNum) / float64(util.GetDefaultInt(ce.EndSessionNum, 1))
			ce.OneStarNum = util.GetInt(strings[constants.CHANNEL_MONITOR_FIELD_ONE_STAR_NUM])
			ce.TwoStarNum = util.GetInt(strings[constants.CHANNEL_MONITOR_FIELD_TWO_STAR_NUM])
			ce.ThreeStarNum = util.GetInt(strings[constants.CHANNEL_MONITOR_FIELD_THREE_STAR_NUM])
			ce.FourStarNum = util.GetInt(strings[constants.CHANNEL_MONITOR_FIELD_FOUR_STAR_NUM])
			ce.FiveStarNum = util.GetInt(strings[constants.CHANNEL_MONITOR_FIELD_FIVE_STAR_NUM])
			ce.SixStarNum = util.GetInt(strings[constants.CHANNEL_MONITOR_FIELD_SIX_STAR_NUM])
			channelEvalC.Upsert(bson.M{"vcc_id": ce.VccId,
				"date": tranTime,
				"channel_id": ce.ChannelId}, ce)
		}

		if cursor <= 0 {
			break
		}
	}
}


