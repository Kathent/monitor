package mq

import (
	"encoding/json"
	"fmt"

	"paas/icsoc-monitor/constants"
	"paas/icsoc-monitor/db"

	"strconv"
	"time"

	"paas/icsoc-monitor/util"

	log "github.com/alecthomas/log4go"
	"github.com/streadway/amqp"
	"strings"
	"gopkg.in/mgo/bson"
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
	routeMap[routeKey_session] = syncSession
	routeMap[routeKey_inQueue] = syncInQueue
	routeMap[routeKey_outQueue] = syncOutQueue
	routeMap[routeKey_message] = syncMessage
	routeMap[routeKey_evaluate] = syncEvaluate
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

	//坐席监控
	agentMonitorKey := fmt.Sprintf(constants.AGENT_MONITOR_HASH_KEY, agentChangeTime,
		agentStatus.VccId, agentStatus.AgentId)

	pipe := db.GetClient().Pipeline()

	//坐席状态
	pipe.HSet(agentMonitorKey, constants.AGENT_MONITOR_FIELD_STATUS,
		fmt.Sprintf("%d", agentStatus.Status))

	//坐席状态开始时间
	pipe.HSet(agentMonitorKey, constants.AGENT_MONITOR_FIELD_STATUS_START_TIME,
		fmt.Sprintf("%d", agentStatus.StampTime))

	arr, err := db.GetClient().HMGet(agentMonitorKey, constants.AGENT_MONITOR_FIELD_STATUS,
		constants.AGENT_MONITOR_FIELD_STATUS_START_TIME).Result()
	if err != nil {
		log.Warn("syncAgentStatus hGet err:%v, mgoStruct:%+v", insertErr, mgoStruct)
		needAck = false
		return
	}

	preStatus := util.GetString(arr[0])
	startTime := util.GetString(arr[1])

	var startTimeInt64 int64
	startTimeInt, err := strconv.Atoi(startTime)
	if err != nil {
		startTimeInt64 = now.Unix()
	}else {
		startTimeInt64 = int64(startTimeInt)
	}

	if preStatus == "2" && agentStatus.Status != 2 {
		pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_ONLINE_TIME_TOTAL,
			now.Unix() - startTimeInt64)
	}

	if preStatus == "1" && agentStatus.Status != 1{
		pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_BUSY_TIME_TOTAL,
			now.Unix() - startTimeInt64)
	}

	_, err = pipe.Exec()
	if err != nil {
		log.Warn("syncAgentStatus pipe exec err:%v, mgoStruct:%+v", insertErr, mgoStruct)
		needAck = false
		return
	}
}

//会话结束
func syncSession(delivery amqp.Delivery) {
	jsonStr := string(delivery.Body)
	if strings.Contains(jsonStr, constants.EVENT_TYPE_SESSION_END) {
		syncSessionEnd(delivery)
	}else if strings.Contains(jsonStr, constants.EVENT_TYPE_SESSION_CONTENT) {
		syncSessionContent(delivery)
	}else {
		log.Warn("syncSession event type wrong. body:%s", jsonStr)
	}
}

//会话内容
func syncSessionContent(delivery amqp.Delivery) {
	var needAck bool
	defer deferAck(delivery, needAck)

	sessionContent := SessionContentMQ{}
	unmarshalErr := json.Unmarshal(delivery.Body, &sessionContent)
	if unmarshalErr != nil {
		log.Warn("syncSessionContent unmarshal err:%v, body:%v", unmarshalErr, delivery.Body)
		return
	}

	database := db.GetSession().DB("")
	query := database.C(constants.STATICS_SR_TABLE_NAME).Find(bson.M{"sid": sessionContent.SessionID})
	var res = SessionRecord{}
	err := query.One(&res)
	if err != nil {
		log.Warn("syncSessionContent cannot find session record. err:%v, content:%+v", err, sessionContent)
		return
	}

	//TODO: 会话内容类型
	ct := SessionContent{
		SessionId: sessionContent.SessionID,
		Index: sessionContent.Index,
		Type: "",
		//Content: sessionContent.Content
	}

	err = database.C(constants.STATICS_SC_TABLE_NAME).Insert(&ct)
	if err != nil {
		log.Warn("syncSessionContent insert err:%v, ct:%+v", err, ct)
		return
	}
}

//会话结束
func syncSessionEnd(delivery amqp.Delivery) {
	var needAck bool
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
		needAck = false
		return
	}

	pipe := db.GetClient().Pipeline()
	//会话结束 坐席监控 当前连接数-1
	agentMonitorKey := fmt.Sprintf(constants.AGENT_MONITOR_HASH_KEY,
		transTime, sessionEndMq.VccId, sessionEndMq.AgentId)

	//会话结束 渠道监控
	channelKey := fmt.Sprintf(constants.CHANNEL_MONITOR_HASH_KEY, transTime, sessionEndMq.VccId,
		sessionEndMq.Source.SourceId)

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
			//if sessionEndMq.EvaluateStatus == "1" {
			//	pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_RECEIVE_EVAL_TIMES, 1)
			//}
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

		add, _ := db.GetClient().SAdd(serveClientKey, sessionEndMq.UserId).Result()
		if add > 0 {//add 成功 坐席服务客户数加1
			pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_SERV_USER_NUM, 1)
			pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_SERVE_USER_NUM, 1)
		}

		//客户对应父会话id
		clientCidSetKey := fmt.Sprintf(constants.CLIENT_CID_SET, transTime, sessionEndMq.VccId,
			sessionEndMq.UserId)
		db.GetClient().SAdd(clientCidSetKey, sessionEndMq.Cid)
		cidNum, _ := db.GetClient().SCard(clientCidSetKey).Result()
		//一次会话客户数
		if cidNum <= 1 {
			pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_ONE_SERV_CLIENT_NUM, 1)
			pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_C_SESSION_NUM, 1)

			//if sessionEndMq.EvaluateStatus == "1" {//已评价父会话数
			//	pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_EVALUATE_NUM, 1)
			//}
		}else {
			//一次会话客户数
			pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_ONE_SERV_CLIENT_NUM, -1)
		}
	}

	//转入会话次数
	if !util.IsEmpty(sessionEndMq.SessionFrom) {
		pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_TRANSFER_IN_TIMES, 1)
	}

	//转出会话次数
	if !util.IsEmpty(sessionEndMq.Next) {
		pipe.HIncrBy(agentMonitorKey, constants.AGENT_MONITOR_FIELD_TRANSFER_OUT_TIMES, 1)
	}

	pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_SESSION_NUM, -1)

	//无效会话
	if sessionEndMq.UserSpeakNum <= 0 {
		pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_INVALID_SESSION_NUM, 1)
	}else {
		//今日完成会话数
		pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_END_SESSION_NUM, 1)

		//首次响应时长
		pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_FIRST_RESP_TOTAL, sessionEndMq.FirstRespTime)

		//会话持续时长
		pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_SESSION_TIME_TOTAL, sr.SessionKeepSecs)
	}

	//排队放弃数
	if sessionEndMq.GiveUpQueueing == 1 {
		pipe.HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_GIVEUP_QUEUE_NUM, 1)
	}
}

//进入排队
func syncInQueue(delivery amqp.Delivery) {
	var needAck bool
	defer deferAck(delivery, needAck)

	im := InQueueMQ{}
	unmarshalErr := json.Unmarshal(delivery.Body, &im)
	if unmarshalErr != nil {
		log.Warn("syncInQueue unmarshal err:%v, body:%v", unmarshalErr, delivery.Body)
		return
	}

	tranTime := TransDate(im.InQueueTime, constants.DATE_FORMATE)
	channelKey := fmt.Sprintf(constants.CHANNEL_MONITOR_HASH_KEY, tranTime, im.VccID, im.ChannelID)

	db.GetClient().HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_QUEUE_NUM, 1)
}

func syncOutQueue(delivery amqp.Delivery) {
	var needAck bool
	defer deferAck(delivery, needAck)

	im := OutQueueMQ{}
	unmarshalErr := json.Unmarshal(delivery.Body, &im)
	if unmarshalErr != nil {
		log.Warn("syncOutQueue unmarshal err:%v, body:%v", unmarshalErr, delivery.Body)
		return
	}

	tranTime := TransDate(im.OutQueueTime, constants.DATE_FORMATE)
	channelKey := fmt.Sprintf(constants.CHANNEL_MONITOR_HASH_KEY, tranTime, im.VccID, im.ChannelID)

	db.GetClient().HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_QUEUE_NUM, -1)
}

//留言
func syncMessage(delivery amqp.Delivery) {
	var needAck bool
	defer deferAck(delivery, needAck)

	im := WebMessageMQ{}
	unmarshalErr := json.Unmarshal(delivery.Body, &im)
	if unmarshalErr != nil {
		log.Warn("syncMessage unmarshal err:%v, body:%v", unmarshalErr, delivery.Body)
		return
	}

	tranTime := TransDate(im.MessageTime, constants.DATE_FORMATE)
	channelKey := fmt.Sprintf(constants.CHANNEL_MONITOR_HASH_KEY, tranTime, im.VccID, im.ChannelID)

	db.GetClient().HIncrBy(channelKey, constants.CHANNEL_MONITOR_FIELD_ADD_MSG_NUM, 1)
}

//满意度评价
func syncEvaluate(delivery amqp.Delivery){
	var needAck bool
	defer deferAck(delivery, needAck)

	eval := EvaluateMsgMQ{}
	unmarshalErr := json.Unmarshal(delivery.Body, &eval)
	if unmarshalErr != nil {
		log.Warn("syncEvaluate unmarshal err:%v, body:%v", unmarshalErr, delivery.Body)
		return
	}

	transTime := TransDate(eval.ConnSuccessTime, constants.DATE_FORMATE)
	agentKey := fmt.Sprintf(constants.AGENT_MONITOR_HASH_KEY, transTime, eval.VccID, eval.AgentID)

	pipe := db.GetClient().Pipeline()
	pipe.HIncrBy()
}