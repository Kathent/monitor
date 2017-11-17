package mq

import (
	"encoding/json"
	"fmt"

	"paas/icsoc-monitor/constants"
	"paas/icsoc-monitor/db"

	log "github.com/alecthomas/log4go"
	"github.com/streadway/amqp"
	"time"
	"gopkg.in/mgo.v2/bson"
	"paas/icsoc-monitor/util"
	"strconv"
)

const (
	routeKey_connSuccess        = "routeKey_connSuccess"
	routeKey_agentBreak         = "routeKey_agentBreak"
	routeKey_agentBreakPassive  = "routeKey_agentBreakPassive"
	routeKey_userBreak          = "routeKey_userBreak"
	routeKey_connRedirect       = "routeKey_connRedirect"
	routeKey_session            = "routeKey_session"
	routeKey_evaluate           = "routeKey_evaluate"
	routeKey_message            = "routeKey_message"
	routeKey_agentStatus        = "routeKey_agentStatus"
	routeKey_userFirstSpeak     = "routeKey_userFirstSpeak"
	routeKey_agentFirstSpeak    = "routeKey_agentFirstSpeak"
	routeKey_agentFirstResponse = "routeKey_agentFirstResponse"
	routeKey_sessionCreate      = "routeKey_sessionCreate"
	routeKey_inQueue            = "routeKey_inQueue"
	routeKey_outQueue           = "routeKey_outQueue"
)

var routeMap map[string]func(delivery amqp.Delivery)
var routeKeyArr []string

func init() {
	routeMap = make(map[string]func(delivery amqp.Delivery))
	routeKeyArr = []string{
		routeKey_connSuccess,
		routeKey_agentBreak,
		routeKey_agentBreakPassive,
		routeKey_userBreak,
		routeKey_connRedirect,
		routeKey_session,
		routeKey_evaluate,
		routeKey_message,
		routeKey_agentStatus,
		routeKey_userFirstSpeak,
		routeKey_agentFirstSpeak,
		routeKey_agentFirstResponse,
		routeKey_sessionCreate,
		routeKey_inQueue,
		routeKey_outQueue,
	}

	routeMap[routeKey_connSuccess] = syncConSuc
	routeMap[routeKey_agentBreak] = syncAgentBreak
	routeMap[routeKey_agentStatus] = syncAgentStatus
	routeMap[routeKey_session] = syncSessionEnd
}

func GetRouteKeys() []string{
	return routeKeyArr
}

func GetRouteFunc(key string) func(delivery amqp.Delivery){
	return routeMap[key]
}

//坐席和用户连接成功
func syncConSuc(delivery amqp.Delivery) {
	var needAck = true

	defer func() {
		if needAck {
			delivery.Acknowledger.Ack(delivery.DeliveryTag, false)
		}else {
			delivery.Acknowledger.Nack(delivery.DeliveryTag, false, true)
		}
	}()

	conSuc := ConnSuccessMQ{}
	unmarshalErr := json.Unmarshal(delivery.Body, &conSuc)
	if unmarshalErr != nil {
		log.Warn("syncConSuc unmarshal err:%v, body:%v", unmarshalErr, delivery.Body)
		return
	}

	pipe := db.GetClient().Pipeline()

	//坐席监控 当前会话数+1
	timeTrans := TransDate(conSuc.ConnSuccessTime, constants.DATE_FORMATE)
	key := fmt.Sprintf(constants.AGENT_MONITOR_HASH_KEY, timeTrans, conSuc.VccID, conSuc.AgentID)
	pipe.HIncrBy(key, constants.AGENT_MONITOR_FIELD_CUR_SESSION_NUM, 1)

	//总会话数+1
	pipe.HIncrBy(key, constants.AGENT_MONITOR_FIELD_TOTAL_SESSION_NUM, 1)

	//渠道监控 当前会话数+1
	sourceKey := fmt.Sprintf(constants.CHANNEL_MONITOR_HASH_KEY, timeTrans, conSuc.VccID, conSuc.ChannelID)
	pipe.HIncrBy(sourceKey, constants.CHANNEL_MONITOR_FIELD_SESSION_NUM, 1)
	//今日处理留言数 判断 后+1
	if true {
		// TODO: 判断是否是留言分配过来的
		pipe.HIncrBy(sourceKey, constants.CHANNEL_MONITOR_FIELD_DEAL_MSG_NUM, 1)
	}

	_, err := pipe.Exec()
	if err != nil {
		log.Warn("syncConSuc pipe exec. err:%v body:%+v", err, conSuc)
		needAck = false
		return
	}
}

func deferAck(delivery amqp.Delivery, needAck bool) {
	if needAck {
		delivery.Acknowledger.Ack(delivery.DeliveryTag, false)
	}else {
		delivery.Acknowledger.Nack(delivery.DeliveryTag, false, true)
	}
}

func syncAgentBreak(delivery amqp.Delivery) {
	var needAck = true
	defer deferAck(delivery, needAck)

	agentBreak := SessionBreakMQ{}

	unmarshalErr := json.Unmarshal(delivery.Body, &agentBreak)
	if unmarshalErr != nil {
		log.Warn("syncAgentBreak unmarshal err:%v, body:%v", unmarshalErr, delivery.Body)
		return
	}

	//TODO:
}

//坐席状态变化
func syncAgentStatus(delivery amqp.Delivery) {
	var needAck = true
	defer deferAck(delivery, needAck)

	agentStatus := AgentStatusMQ{}
	unmarshalErr := json.Unmarshal(delivery.Body, &agentStatus)
	if unmarshalErr != nil {
		log.Warn("syncAgentStatus unmarshal err:%v, body:%v", unmarshalErr, delivery.Body)
		return
	}

	now := time.Now()
	timeTrans := TransTime(now, constants.DATE_FORMATE)
	key := fmt.Sprintf(constants.AGENT_MONITOR_HASH_KEY, timeTrans, agentStatus.VccId, agentStatus.AgentId)
	fields := []string{
		constants.AGENT_MONITOR_FIELD_WORKER_ID, constants.AGENT_MONITOR_FIELD_DEP_ID,
		constants.AGENT_MONITOR_FIELD_STATUS, constants.AGENT_MONITOR_FIELD_STATUS_START_TIME,
	}

	client := db.GetClient()
	res, getErr := client.HMGet(key, fields...).Result()
	if getErr != nil {
		log.Warn("syncAgentStatus getErr:%v, obj:%+v", getErr, agentStatus)
		needAck = false
		return
	}

	preStartInt, _ := strconv.Atoi(fields[3])
	collection := db.GetSession().DB("").C(constants.STATICS_AS_TABLE_NAME)
	agentChangeTime := TransDate(agentStatus.StampTime, constants.DATE_FORMATE)
	//TODO: OP_TYPE.
	mgoStruct := AgentStatus{
		VccId: agentStatus.VccId,
		AgentId: agentStatus.AgentId,
		WorkerId: util.GetString(res[0]),
		DeptId: util.GetString(res[1]),
		Date: agentChangeTime,
		PreStatus: fields[2],
		Status: string(agentStatus.Status + 1),
		OpType: agentStatus.EventType,
		PreStatusSecs: now.Unix() - int64(preStartInt),
	}

	insertErr := collection.Insert(&mgoStruct)
	if insertErr != nil {
		log.Warn("syncAgentStatus insert err:%v, mgoStruct:%+v", insertErr, mgoStruct)
		return
	}
}

//会话结束
func syncSessionEnd(delivery amqp.Delivery) {
	var needAck = true
	defer deferAck(delivery, needAck)

	sessionEndMq := SessionEndMq{}
	unmarshalErr := json.Unmarshal(delivery.Body, &sessionEndMq)
	if unmarshalErr != nil {
		log.Warn("syncSessionEnd unmarshal err:%v, body:%v", unmarshalErr, delivery.Body)
		return
	}

	transTime := TransDate(sessionEndMq.ConSucTime, constants.DATE_FORMATE)
	//插入会话记录
	sr := SessionRecord{
		Sid: sessionEndMq.SessionId,
		Cid: "",//TODO:
		Date: transTime,
		VccId: sessionEndMq.VccId,
		Name: sessionEndMq.Name,
		AgentId: sessionEndMq.AgentId,
		WorkerId: sessionEndMq.WorkerId,
		ClientNewsNum: sessionEndMq.UserSpeakNum,
		AgentNewsNum: sessionEndMq.AgentSpeakTime,
		SessionStartTime: TransDate(sessionEndMq.ConSucTime, constants.DATE_FORMATE_ALL),
		SessionEndTime: TransDate(sessionEndMq.SessionEndTime, constants.DATE_FORMATE_ALL),
		FirstRespSecs: sessionEndMq.FirstRespTime,
		SessionKeepSecs: sessionEndMq.SessionEndTime - sessionEndMq.ConSucTime,
		CreateType: sessionEndMq.CreateType,
		EndType: sessionEndMq.EndType,
		SourceType: sessionEndMq.SourceType,
	}

	collection := db.GetSession().DB("").C(constants.STATICS_SR_TABLE_NAME)
	insertErr := collection.Insert(sr)
	if insertErr != nil {
		log.Warn("syncSessionEnd insert err:%v, mgoStruct:%+v", insertErr, sr)
		return
	}

	pipe := db.GetClient().Pipeline()
	//会话结束 坐席监控 当前连接数-1
	agentMonitorKey := fmt.Sprintf(constants.AGENT_MONITOR_HASH_KEY,
		transTime, sessionEndMq.VccId, sessionEndMq.AgentId)

	pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_CUR_SESSION_NUM, -1)

	//今日无效会话数
	if sessionEndMq.UserSpeakNum <= 0 {
		pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_INVALID_SESSION_NUM, 1)
	}else {
		//今日独立会话数
		if util.IsEmpty(sessionEndMq.Next) && util.IsEmpty(sessionEndMq.SessionFrom) {
			pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_INDEP_SESSION_NUM, 1)

			//邀评数
			if sessionEndMq.EvaluateStatus == "0" || sessionEndMq.EvaluateStatus == "1" {
				pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_REQUIRE_EVAL_TIMES, 1)
			}

			//受评数
			if sessionEndMq.EvaluateStatus == "1" {
				pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_RECEIVE_EVAL_TIMES, 1)
			}
		}else {
			//参与转接数
			pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_DEP_SESSION_NUM, 1)
		}

		//今日首次响应时长总计
		pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_FIRST_RESP_TIME, sessionEndMq.FirstRespTime)
		pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_SESSION_TIME_TOTAL, sr.SessionKeepSecs)

		//回复消息数
		pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_REPLY_MSG_TIMES, int64(sessionEndMq.AgentSpeakTime))

		//用户消息数
		pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_RECEIVE_MSG_TIMES, int64(sessionEndMq.UserSpeakNum))

		//坐席服务客户数
		serveClientKey := fmt.Sprintf(constants.AGENT_SERVE_CLIENT_SET, transTime, sessionEndMq.VccId,
			sessionEndMq.AgentId)
		pipe.SAdd(serveClientKey, sessionEndMq.UserId)

		//客户对应父会话id
		clientCidSetKey := fmt.Sprintf(constants.CLIENT_CID_SET, transTime, sessionEndMq.VccId,
			sessionEndMq.UserId)
		pipe.SAdd(clientCidSetKey, sessionEndMq.Cid)
	}

	//转入会话次数
	if !util.IsEmpty(sessionEndMq.SessionFrom) {
		pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_TRANSFER_IN_TIMES, 1)
	}

	//转出会话次数
	if !util.IsEmpty(sessionEndMq.Next) {
		pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_TRANSFER_OUT_TIMES, 1)
	}


	//会话结束 渠道监控
	channelKey := fmt.Sprintf(constants.CHANNEL_MONITOR_HASH_KEY, transTime, sessionEndMq.VccId,
		sessionEndMq.Source.SourceId)

	pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_SESSION_NUM, 1)

	//无效会话
	if sessionEndMq.UserSpeakNum <= 0 {
		pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_INVALID_SESSION_NUM, 1)
	}else {
		//今日完成会话数
		pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_END_SESSION_NUM, 1)
	}

	//排队放弃数
	if sessionEndMq.GiveUpQueueing == 1 {
		pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_GIVEUP_QUEUE_NUM, 1)
	}
}