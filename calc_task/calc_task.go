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
	t := time.NewTicker(time.Minute)
	select {
	case t.C:
		syncRedisToMongoDb()
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
			return
		}

		if len(keys) <= 0 {
			log4go.Info("syncRedisToMongoDb scan empty keys.")
			return
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
			aw.InSessionNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_AGENTID])
			aw.ReceiveEvalTimes = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_RECEIVE_EVAL_TIMES])
			aw.EvalRatio = float64(aw.ReceiveEvalTimes) / float64(aw.InSessionNum)
			aw.ReplyNewsNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_REPLY_MSG_TIMES])
			aw.ReceiveMsgTimes = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_RECEIVE_MSG_TIMES])
			aw.ArRatio = float64(aw.ReceiveMsgTimes) / float64(aw.ReplyNewsNum)
			aw.TotalSessionNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_TOTAL_SESSION_NUM])
			aw.SessionKeepSecs = int64(util.GetInt(strings[constants.AGENT_MONITOR_FIELD_SESSION_TIME_TOTAL]))
			aw.AvgSessionSecs = aw.SessionKeepSecs / int64(aw.TotalSessionNum)
			aw.FirstRespSecs = int64(util.GetInt(strings[constants.AGENT_MONITOR_FIELD_FIRST_RESP_TIME]))
			aw.AvgResponseSecs = aw.FirstRespSecs / int64(aw.TotalSessionNum)
			aw.Date = tranTime
			aw.DeptId = util.GetString(strings[constants.AGENT_MONITOR_FIELD_DEP_ID])
			aw.BusySecs = int64(util.GetInt(strings[constants.AGENT_MONITOR_FIELD_BUSY_TIME_TOTAL]))
			aw.GroupId = util.GetString(strings[constants.AGENT_MONITOR_FIELD_GROUP_IDS])
			aw.InvalidSessionNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_INVALID_SESSION_NUM])
			aw.InvalidRatio = float64(aw.InvalidSessionNum) / float64(aw.TotalSessionNum)
			aw.OneServeClientNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_ONE_SERV_CLIENT_NUM])
			aw.OnlineSecs = int64(util.GetInt(strings[constants.AGENT_MONITOR_FIELD_ONLINE_TIME_TOTAL]))
			aw.ServeUserNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_SERV_USER_NUM])
			aw.TransInSessionNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_TRANSFER_IN_TIMES])
			aw.TransOutSessionNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_TRANSFER_OUT_TIMES])
			aw.ReceiveSessionNUm = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_DEP_SESSION_NUM])

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
			ae.RateEvaluate = float64(ae.EvaluateNum) / float64(ae.IndepSessionNum)
			ae.OneStarNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_ONE_STAR_NUM])
			ae.TwoStarNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_TWO_STAR_NUM])
			ae.ThreeStarNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_THREE_STAR_NUM])
			ae.FourStarNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_FOUR_STAR_NUM])
			ae.FiveStarNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_FIVE_STAR_NUM])
			ae.SixStarNum = util.GetInt(strings[constants.AGENT_MONITOR_FIELD_SIX_STAR_NUM])
			agentEvalC.Upsert(bson.M{"vcc_id": aw.VccId,
				"date": tranTime,
				"ag_id": aw.AgentId}, ae)
		}
	}

	keyPattern = fmt.Sprintf(constants.CHANNEL_MONITOR_HASH_KEY, tranTime, "*", "*")
	for {
		keys, cursor, err = db.GetClient().Scan(cursor, keyPattern, count).Result()
		if err != nil {
			log4go.Warn("syncRedisToMongoDb scan err:%v", err)
			return
		}

		if len(keys) <= 0 {
			log4go.Info("syncRedisToMongoDb scan empty keys.")
			return
		}

		for _, v := range keys {
			strings, err := db.GetClient().HGetAll(v).Result()
			if err != nil {
				log4go.Warn("syncRedisToMongoDb scan err:%v, key:%s", err, v)
				continue
			}

			ce := mq.ChannelEval{}
			ce.ChannelId = util.GetString(strings[constants.CHANNEL_MONITOR_FIELD_CHANNEL_ID])
			ce.VccId = util.GetString(strings[constants.CHANNEL_MONITOR_FIELD_VCCID])
			ce.SourceName = util.GetString(strings[constants.CHANNEL_MONITOR_FIELD_SOURCE_NAME])
			ce.SourceType = util.GetString(strings[constants.CHANNEL_MONITOR_FIELD_SOURCE_TYPE])
			ce.EndSessionNum = util.GetInt(strings[constants.CHANNEL_MONITOR_FIELD_END_SESSION_NUM])
			ce.EvaluateNum = util.GetInt(strings[constants.CHANNEL_MONITOR_FIELD_EVALUATE_NUM])
			ce.RateEvaluate = float64(ce.EvaluateNum) / float64(ce.EndSessionNum)
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
	}
}


