package server

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"paas/icsoc-monitor/constants"
	"paas/icsoc-monitor/db"
	"paas/icsoc-monitor/mq"
	"paas/icsoc-monitor/util"

	"fmt"

	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func AgentWorkStatics(context *gin.Context){
	var req MonitorRequest
	var res MonitorResponse
	context.Bind(&req)

	agentReq, err := transAgentReq(&req, constants.STATICS_AW_FILED_WORKER_ID)
	if err != nil {
		res.Code = -1
		res.Message = err.Error()
		context.JSON(http.StatusOK, res)
		log4go.Warn("tranAgentReq err:%v, req:%+v", err, req)
		return
	}

	collection := db.GetSession().DB("").C(constants.STATICS_AW_TABLE_NAME)
	var pipe *mgo.Pipe
	if util.IsEmpty(req.AgentId) {
		split := strings.Split(req.AgentId, ",")
		pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vccId": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lt": agentReq.End}},
			"agentId": bson.M{"$in": split},},
			{"$group":	bson.M{"_id": "$worker_id", "group_ids": bson.M{"$first": ""},
				"online_secs": bson.M{"$sum":"$online_secs"},
				"busy_secs": bson.M{"$sum": "$busy_secs"},
				"in_session_num": bson.M{"$sum": "$in_session_num"},
				"invalid_session_num": bson.M{"$sum": "$invalid_session_num"},
				"receive_session_num": bson.M{"$sum": "$receive_session_num"},
				"trans_in_session_num": bson.M{"$sum": "$trans_in_session_num"},
				"trans_out_session_num": bson.M{"$sum": "$trans_out_session_num"},
				"reply_news_num": bson.M{"$sum": "$reply_news_num"}}},
			{"$sort": bson.M{agentReq.SortField: agentReq.GetSortSym()}},
		})
	}else {
		pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vccId": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lt": agentReq.End}}},
			{"$group":	bson.M{"_id": "$worker_id", "group_ids": bson.M{"$first": ""},
				"online_secs": bson.M{"$sum":"$online_secs"},
				"busy_secs": bson.M{"$sum": "$busy_secs"},
				"in_session_num": bson.M{"$sum": "$in_session_num"},
				"invalid_session_num": bson.M{"$sum": "$invalid_session_num"},
				"receive_session_num": bson.M{"$sum": "$receive_session_num"},
				"trans_in_session_num": bson.M{"$sum": "$trans_in_session_num"},
				"trans_out_session_num": bson.M{"$sum": "$trans_out_session_num"},
				"reply_news_num": bson.M{"$sum": "$reply_news_num"}}},
			{"$sort": bson.M{agentReq.SortField: agentReq.GetSortSym()}},
		})
	}

	var arr = make([]mq.AgentWork, 0)
	pipe.All(&arr)

	total := len(arr)
	res.Data.Total = total / agentReq.PageSize
	res.Code = 0
	res.Message = "suc"

	if agentReq.StartIndex < total {
		if agentReq.EndIndex >= total {
			arr = arr[agentReq.StartIndex:]
		}else {
			arr = arr[agentReq.StartIndex: agentReq.EndIndex]
		}
	}

	res.Data.Rows = make([]interface{}, 0)
	for _, v := range arr {
		res.Data.Rows = append(res.Data.Rows, v)
	}

	context.JSON(http.StatusOK, res)
	log4go.Info("AgentWorkStatics suc return. req:%+v, res:%+v:", req, res)
	return
}

func transAgentReq(req *MonitorRequest, defaultField string) (*AgentRequest, error) {
	if util.IsEmpty(req.VccId) {
		return nil, errors.New("vccId is wrong")
	}

	var (
		ar AgentRequest
		err error
	)

	ar.VccId = req.VccId
	now := time.Now()
	var start time.Time
	var end time.Time
	if util.IsEmpty(req.Start) {
		start = now
	}else {
		start, err = time.Parse(constants.DATE_FORMATE, req.Start)
	}

	if err != nil {
		return nil, errors.New("start is wrong")
	}

	if util.IsEmpty(req.End){
		end = time.Now()
	}else {
		end, err = time.Parse(constants.DATE_FORMATE, req.End)
	}

	if start.After(end) {
		return nil, errors.New("start after end")
	}

	if err != nil {
		return nil, errors.New("end is wrong")
	}

	ar.Start = req.Start
	ar.End = req.End

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	ar.PageSize = req.PageSize

	if req.Page <= 1 {
		ar.Page = 1
		ar.StartIndex = 0
		ar.EndIndex = req.PageSize
	}else {
		ar.Page = req.Page
		ar.StartIndex = req.PageSize * (req.Page - 1)
		ar.EndIndex = req.PageSize * req.Page
	}

	if util.IsEmpty(req.SortField) {
		ar.SortField = defaultField
	}

	if req.SortOrder != constants.SORT_ORDER_DESCEND{
		ar.SortOrder = constants.SORT_ORDER_DEFAULT
	}

	if util.IsEmpty(req.Status) {
		ar.Status = "0"
	}

	return &ar, nil
}

func AgentAchievements(context *gin.Context){
	var req MonitorRequest
	var res MonitorResponse
	context.Bind(&req)

	agentReq, err := transAgentReq(&req, constants.STATICS_AW_FILED_WORKER_ID)
	if err != nil {
		res.Code = -1
		res.Message = err.Error()
		context.JSON(http.StatusOK, res)
		log4go.Warn("AgentAchievements tranAgentReq err:%v, req:%+v", err, req)
		return
	}

	collection := db.GetSession().DB("").C(constants.STATICS_AW_TABLE_NAME)
	var pipe *mgo.Pipe
	if util.IsEmpty(req.AgentId) {
		split := strings.Split(req.AgentId, ",")
		pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vccId": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lt": agentReq.End}},
			"agentId": bson.M{"$in": split},},
			{"$group":	bson.M{"_id": "$worker_id", "group_ids": bson.M{"$first": ""},
				"receive_msg_times": bson.M{"$sum":"$receive_msg_times"},
				"reply_news_num": bson.M{"$sum": "$reply_news_num"},
				"in_session_num": bson.M{"$sum": "$in_session_num"},
				"total_session_num": bson.M{"$sum": "$total_session_num"},
				"invalid_session_num": bson.M{"$sum": "$invalid_session_num"},
				"serv_user_num": bson.M{"$sum": "$serv_user_num"},
				"one_serv_client_num": bson.M{"$sum": "$one_serv_client_num"},
				"receive_eval_times": bson.M{"$sum": "$receive_eval_times"},
				"first_response_secs": bson.M{"$sum": "$first_response_secs"},
				"session_keep_secs": bson.M{"$sum": "$session_keep_secs"},
				}},
			{"$sort": bson.M{agentReq.SortField: agentReq.GetSortSym()}},
		})
	}else {
		pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vccId": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lt": agentReq.End}}},
			{"$group":	bson.M{"_id": "$worker_id", "group_ids": bson.M{"$first": ""},
				"receive_msg_times": bson.M{"$sum":"$receive_msg_times"},
				"reply_news_num": bson.M{"$sum": "$reply_news_num"},
				"in_session_num": bson.M{"$sum": "$in_session_num"},
				"total_session_num": bson.M{"$sum": "$total_session_num"},
				"invalid_session_num": bson.M{"$sum": "$invalid_session_num"},
				"serv_user_num": bson.M{"$sum": "$serv_user_num"},
				"one_serv_client_num": bson.M{"$sum": "$one_serv_client_num"},
				"receive_eval_times": bson.M{"$sum": "$receive_eval_times"},
				"first_response_secs": bson.M{"$sum": "$first_response_secs"},
				"session_keep_secs": bson.M{"$sum": "$session_keep_secs"}}},
			{"$sort": bson.M{agentReq.SortField: agentReq.GetSortSym()}},
		})
	}

	var arr = make([]mq.AgentAchievements, 0)
	pipe.All(&arr)

	total := len(arr)
	res.Data.Total = total / agentReq.PageSize
	res.Code = 0
	res.Message = "suc"

	if agentReq.StartIndex < total {
		if agentReq.EndIndex >= total {
			arr = arr[agentReq.StartIndex:]
		}else {
			arr = arr[agentReq.StartIndex: agentReq.EndIndex]
		}
	}

	res.Data.Rows = make([]interface{}, 0)
	for _, v := range arr {
		v.AvgResponseSecs = v.FirstRespSecs / int64(v.TotalSessionNum)
		v.AvgSessionSecs = v.SessionKeepSecs / int64(v.TotalSessionNum)
		v.ArRatio = float64(v.ReceiveMsgTimes) / float64(v.ReplyNewsNum)
		v.EvalRatio = float64(v.ReceiveEvalTimes) / float64(v.InSessionNum)
		res.Data.Rows = append(res.Data.Rows, v)
	}

	context.JSON(http.StatusOK, res)
	log4go.Info("AgentAchievements suc return. req:%+v, res:%+v:", req, res)
	return
}

//坐席状态变更
func AgentStatusRecord(context *gin.Context){
	log4go.Info("AgentStatusRecord enter ctx:%+v",context)
	var req MonitorRequest
	var res MonitorResponse
	bindErr := context.Bind(&req)
	if bindErr != nil {
		log4go.Warn("AgentStatusRecord bindErr:%v ctx:%+v", bindErr, context)
	}

	agentReq, err := transAgentReq(&req, constants.STATICS_AW_FILED_WORKER_ID)
	if err != nil {
		res.Code = -1
		res.Message = err.Error()
		context.JSON(http.StatusOK, res)
		log4go.Warn("AgentStatusRecord tranAgentReq err:%v, req:%+v", err, req)
		return
	}

	log4go.Info("AgentStatusRecord agentReq:%+v", agentReq)
	var arr = make([]mq.AgentStatus, 0)

	total := len(arr)
	res.Data.Total = total / agentReq.PageSize
	res.Code = 0
	res.Message = "suc"

	if agentReq.StartIndex < total {
		if agentReq.EndIndex >= total {
			arr = arr[agentReq.StartIndex:]
		}else {
			arr = arr[agentReq.StartIndex: agentReq.EndIndex]
		}
	}

	res.Data.Rows = make([]interface{}, 0)
	for _, v := range arr {
		res.Data.Rows = append(res.Data.Rows, v)
	}

	collection := db.GetSession().DB("").C(constants.STATICS_AS_TABLE_NAME)
	if util.IsEmpty(req.AgentId) {
		collection.Find(bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End}}).
			Sort(agentReq.GetSortString()).All(&arr)
	}else {
		split := strings.Split(req.AgentId, ",")
		collection.Find(bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End},
			"agentId": bson.M{"$in": split},}).
			Sort(agentReq.GetSortString()).All(&arr)
	}

	var tmpMap map[string]interface{}
	for _, v := range arr {
		tmpMap = make(map[string]interface{})
		tmpMap["group_ids"] = v.GroupId
		tmpMap["worker_id"] = v.WorkerId
		tmpMap["pre_status"] = v.PreStatus
		tmpMap["status"] = v.Status
		tmpMap["op_type"] = v.OpType
		tmpMap["time"] = v.Time
		tmpMap["pre_status_secs"] = fmt.Sprintf("%d", v.PreStatusSecs)
		res.Data.Rows = append(res.Data.Rows, tmpMap)
	}

	context.JSON(http.StatusOK, res)
	log4go.Info("AgentStatusRecord suc return. req:%+v, res:%+v:", req, res)
	return
}

//会话记录查询
func AgentSessionRecord(context *gin.Context){
	var req MonitorRequest
	var res MonitorResponse
	context.Bind(&req)

	agentReq, err := transAgentReq(&req, constants.STATICS_SR_FIELD_WORKER_ID)
	if err != nil {
		res.Code = -1
		res.Message = err.Error()
		context.JSON(http.StatusOK, res)
		log4go.Warn("AgentAchievements tranAgentReq err:%v, req:%+v", err, req)
		return
	}

	var arr = make([]mq.SessionRecord, 0)
	res.Data.Rows = make([]interface{}, 0)

	collection := db.GetSession().DB("").C(constants.STATICS_SR_TABLE_NAME)
	if util.IsEmpty(req.AgentId) {
		collection.Find(bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End}}).
			Sort(agentReq.GetSortString()).All(&arr)
	}else {
		split := strings.Split(req.AgentId, ",")
		collection.Find(bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End},
			"agentId": bson.M{"$in": split},}).
			Sort(agentReq.GetSortString()).All(&arr)
	}

	total := len(arr)
	res.Data.Total = total / agentReq.PageSize
	res.Code = 0
	res.Message = "suc"

	if agentReq.StartIndex < total {
		if agentReq.EndIndex >= total {
			arr = arr[agentReq.StartIndex:]
		}else {
			arr = arr[agentReq.StartIndex: agentReq.EndIndex]
		}
	}

	var tmpMap map[string]interface{}
	for _, v := range arr {
		tmpMap = make(map[string]interface{})
		tmpMap["name"] = v.Name
		tmpMap["worker_id"] = v.WorkerId
		tmpMap["client_news_num"] = fmt.Sprintf("%d", v.ClientNewsNum)
		tmpMap["agent_news_num"] = fmt.Sprintf("%d", v.AgentNewsNum)
		tmpMap["session_start_time"] = v.SessionStartTime
		tmpMap["session_end_time"] = v.SessionEndTime
		tmpMap["first_response_secs"] = fmt.Sprintf("%d", v.FirstRespSecs)
		tmpMap["session_keep_secs"] = fmt.Sprintf("%d", v.SessionKeepSecs)
		tmpMap["session_create_type"] = v.CreateType
		tmpMap["session_end_type"] = v.EndType
		tmpMap["source_type"] = fmt.Sprintf("%d", v.SourceType)
		tmpMap["source_name"] = v.SourceName
		tmpMap["is_evaluate"] = fmt.Sprintf("%d", v.IsEvaluate)
		tmpMap["evaluate"] = v.Evaluate
		tmpMap["evaluate_explain"] = v.EvalExplain
		res.Data.Rows = append(res.Data.Rows, tmpMap)
	}

	context.JSON(http.StatusOK, res)
	log4go.Info("AgentStatusRecord suc return. req:%+v, res:%+v:", req, res)
	return
}


//坐席满意度
func AgentEval(context *gin.Context){
	var req MonitorRequest
	var res MonitorResponse
	context.Bind(&req)

	agentReq, err := transAgentReq(&req, constants.STATICS_AW_FILED_WORKER_ID)
	if err != nil {
		res.Code = -1
		res.Message = err.Error()
		context.JSON(http.StatusOK, res)
		log4go.Warn("AgentAchievements tranAgentReq err:%v, req:%+v", err, req)
		return
	}

	collection := db.GetSession().DB("").C(constants.STATICS_AW_TABLE_NAME)
	var pipe *mgo.Pipe
	if util.IsEmpty(req.AgentId) {
		split := strings.Split(req.AgentId, ",")
		pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lt": agentReq.End}},
			"agentId": bson.M{"$in": split},},
			{"$group":	bson.M{"_id": "$worker_id", "group_ids": bson.M{"$first": ""},
				"indep_session_num": bson.M{"$sum":"$indep_session_num"},
				"require_eval_times": bson.M{"$sum": "$require_eval_times"},
				"evaluate_num": bson.M{"$sum": "$evaluate_num"},
				"one_star_num": bson.M{"$sum": "$one_star_num"},
				"two_star_num": bson.M{"$sum": "$two_star_num"},
				"three_star_num": bson.M{"$sum": "$three_star_num"},
				"four_star_num": bson.M{"$sum": "$four_star_num"},
				"five_star_num": bson.M{"$sum": "$five_star_num"},
				"six_star_num": bson.M{"$sum": "$six_star_num"},
			}},
			{"$sort": bson.M{agentReq.SortField: agentReq.GetSortSym()}},
		})
	}else {
		pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lt": agentReq.End}}},
			{"$group":	bson.M{"_id": "$worker_id", "group_ids": bson.M{"$first": ""},
				"indep_session_num": bson.M{"$sum":"$indep_session_num"},
				"require_eval_times": bson.M{"$sum": "$require_eval_times"},
				"evaluate_num": bson.M{"$sum": "$evaluate_num"},
				"one_star_num": bson.M{"$sum": "$one_star_num"},
				"two_star_num": bson.M{"$sum": "$two_star_num"},
				"three_star_num": bson.M{"$sum": "$three_star_num"},
				"four_star_num": bson.M{"$sum": "$four_star_num"},
				"five_star_num": bson.M{"$sum": "$five_star_num"},
				"six_star_num": bson.M{"$sum": "$six_star_num"},
			}},
			{"$sort": bson.M{agentReq.SortField: agentReq.GetSortSym()}},
		})
	}

	var arr = make([]mq.AgentEval, 0)
	pipe.All(&arr)

	total := len(arr)
	res.Data.Total = total / agentReq.PageSize
	res.Code = 0
	res.Message = "suc"

	if agentReq.StartIndex < total {
		if agentReq.EndIndex >= total {
			arr = arr[agentReq.StartIndex:]
		}else {
			arr = arr[agentReq.StartIndex: agentReq.EndIndex]
		}
	}

	res.Data.Rows = make([]interface{}, 0)
	var tmpMap map[string]interface{}
	for _, v := range arr {
		v.RateEvaluate = float64(v.EvaluateNum) / float64(v.IndepSessionNum)
		tmpMap = make(map[string]interface{})
		tmpMap["group_ids"] = v.GroupId
		tmpMap["worker_id"] = v.WorkerId
		tmpMap["indep_session_num"] = fmt.Sprintf("%d", v.IndepSessionNum)
		tmpMap["require_eval_times"] = fmt.Sprintf("%d", v.RequireEvalTimes)
		tmpMap["evaluate_num"] = fmt.Sprintf("%d", v.RequireEvalTimes)
		tmpMap["one_star_num"] = fmt.Sprintf("%d", v.OneStarNum)
		tmpMap["two_star_num"] = fmt.Sprintf("%d", v.TwoStarNum)
		tmpMap["three_star_num"] = fmt.Sprintf("%d", v.ThreeStarNum)
		tmpMap["four_star_num"] = fmt.Sprintf("%d", v.FourStarNum)
		tmpMap["five_star_num"] = fmt.Sprintf("%d", v.FiveStarNum)
		tmpMap["six_star_num"] = fmt.Sprintf("%d", v.SixStarNum)
		tmpMap["rate_evaluate"] = fmt.Sprintf("%.2f", v.RateEvaluate)
	}

	context.JSON(http.StatusOK, res)
	log4go.Info("AgentEval suc return. req:%+v, res:%+v:", req, res)
	return
}

//渠道满意度
func ChannelEval(context *gin.Context){
	var req MonitorRequest
	var res MonitorResponse
	context.Bind(&req)

	agentReq, err := transAgentReq(&req, constants.STATICS_CE_FIELD_SOURCE_TYPE)
	if err != nil {
		res.Code = -1
		res.Message = err.Error()
		context.JSON(http.StatusOK, res)
		log4go.Warn("AgentAchievements tranAgentReq err:%v, req:%+v", err, req)
		return
	}

	collection := db.GetSession().DB("").C(constants.STATICS_AW_TABLE_NAME)
	var pipe *mgo.Pipe
	pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vcc_id": agentReq.VccId,
		"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End}}},
		{"$group":	bson.M{"_id": "$worker_id", "group_ids": bson.M{"$first": ""},
			"source_name": bson.M{"$first": ""},
			"source_type": bson.M{"$first": ""},
			"end_session_num": bson.M{"$sum":"$end_session_num"},
			"evaluate_num": bson.M{"$sum": "$evaluate_num"},
			"one_star_num": bson.M{"$sum": "$one_star_num"},
			"two_star_num": bson.M{"$sum": "$two_star_num"},
			"three_star_num": bson.M{"$sum": "$three_star_num"},
			"four_star_num": bson.M{"$sum": "$four_star_num"},
			"five_star_num": bson.M{"$sum": "$five_star_num"},
			"six_star_num": bson.M{"$sum": "$six_star_num"},
		}},
		{"$sort": bson.M{agentReq.SortField: agentReq.GetSortSym()}},
	})

	var arr = make([]mq.ChannelEval, 0)
	pipe.All(&arr)

	total := len(arr)
	res.Data.Total = total / agentReq.PageSize
	res.Code = 0
	res.Message = "suc"

	if agentReq.StartIndex < total {
		if agentReq.EndIndex >= total {
			arr = arr[agentReq.StartIndex:]
		}else {
			arr = arr[agentReq.StartIndex: agentReq.EndIndex]
		}
	}

	res.Data.Rows = make([]interface{}, 0)
	var tmpMap map[string]interface{}
	for _, v := range arr {
		v.RateEvaluate = float64(v.EvaluateNum) / float64(v.EndSessionNum)

		tmpMap = make(map[string]interface{})
		tmpMap["source_name"] = v.SourceName
		tmpMap["source_type"] = v.SourceType
		tmpMap["end_session_num"] = fmt.Sprintf("%d", v.EndSessionNum)
		tmpMap["evaluate_num"] = fmt.Sprintf("%d", v.EvaluateNum)
		tmpMap["one_star_num"] = fmt.Sprintf("%d", v.OneStarNum)
		tmpMap["two_star_num"] = fmt.Sprintf("%d", v.TwoStarNum)
		tmpMap["three_star_num"] = fmt.Sprintf("%d", v.ThreeStarNum)
		tmpMap["four_star_num"] = fmt.Sprintf("%d", v.FourStarNum)
		tmpMap["five_star_num"] = fmt.Sprintf("%d", v.FiveStarNum)
		tmpMap["six_star_num"] = fmt.Sprintf("%d", v.SixStarNum)
		tmpMap["rate_evaluate"] = fmt.Sprintf("%.2f", v.RateEvaluate)
	}

	context.JSON(http.StatusOK, res)
	log4go.Info("ChannelEval suc return. req:%+v, res:%+v:", req, res)
	return
}

//坐席监控
func AgentMonitor(context *gin.Context) {
	var req MonitorRequest
	var res MonitorResponse
	context.Bind(&req)

	agentReq, err := transAgentReq(&req, constants.AGENT_MONITOR_FIELD_WORKER_ID)
	if err != nil {
		res.Code = -1
		res.Message = err.Error()
		context.JSON(http.StatusOK, res)
		log4go.Warn("AgentMonitor tranAgentReq err:%v, req:%+v", err, req)
		return
	}

	transTime := mq.TransTime(time.Now(), constants.DATE_FORMATE)
	key := fmt.Sprintf(constants.AGENT_MONITOR_HASH_KEY, transTime, agentReq.VccId, "*")
	keys, err := db.GetClient().HKeys(key).Result()
	if err != nil {
		res.Code = -1
		res.Message = err.Error()
		context.JSON(http.StatusOK, res)
		log4go.Warn("AgentMonitor hkeys err:%v, req:%+v", err, req)
		return
	}

	arr := make([]map[string]string, 0)
	for _, v := range keys {
		resMap, err := db.GetClient().HGetAll(v).Result()
		if err != nil {
			continue
		}

		arr = append(arr, resMap)
	}

	total := len(arr)

	if agentReq.StartIndex < total {
		if agentReq.EndIndex >= total {
			arr = arr[agentReq.StartIndex:]
		}else {
			arr = arr[agentReq.StartIndex: agentReq.EndIndex]
		}
	}
	res.Data.Total = total / agentReq.PageSize
	res.Code = 0
	res.Message = "suc"
	res.Data.Rows = make([]interface{}, 0)
	for _, v := range arr {
		res.Data.Rows = append(res.Data.Rows, v)
	}

	log4go.Info("AgentMonitor suc return. req:%+v, res:%+v:", req, res)
	context.JSON(http.StatusOK, res)
}

//渠道监控
func ChannelMonitor(context *gin.Context){
	var req MonitorRequest
	var res MonitorResponse
	context.Bind(&req)

	agentReq, err := transAgentReq(&req, constants.AGENT_MONITOR_FIELD_WORKER_ID)
	if err != nil {
		res.Code = -1
		res.Message = err.Error()
		context.JSON(http.StatusOK, res)
		log4go.Warn("ChannelMonitor tranAgentReq err:%v, req:%+v", err, req)
		return
	}

	transTime := mq.TransTime(time.Now(), constants.DATE_FORMATE)
	key := fmt.Sprintf(constants.CHANNEL_MONITOR_HASH_KEY, transTime, agentReq.VccId, "*")
	log4go.Info("ChannelMonitor key:%s", key)
	keys, err := db.GetClient().Keys(key).Result()
	if err != nil {
		res.Code = -1
		res.Message = err.Error()
		context.JSON(http.StatusOK, res)
		log4go.Warn("ChannelMonitor hkeys err:%v, req:%+v", err, req)
		return
	}

	arr := make([]map[string]string, 0)
	for _, v := range keys {
		resMap, err := db.GetClient().HGetAll(v).Result()
		if err != nil {
			continue
		}

		arr = append(arr, resMap)
	}

	total := len(arr)

	if agentReq.StartIndex < total {
		if agentReq.EndIndex >= total {
			arr = arr[agentReq.StartIndex:]
		}else {
			arr = arr[agentReq.StartIndex: agentReq.EndIndex]
		}
	}
	res.Data.Total = total / agentReq.PageSize
	res.Code = 0
	res.Message = "suc"
	res.Data.Rows = make([]interface{}, 0)

	var tmpMap map[string]interface{}
	for _, v := range arr {
		tmpMap = make(map[string]interface{})
		tmpMap["source_name"] = v[constants.CHANNEL_MONITOR_FIELD_SOURCE_NAME]
		tmpMap["source_type"] = v[constants.CHANNEL_MONITOR_FIELD_SOURCE_TYPE]
		tmpMap["session_num"] = v[constants.CHANNEL_MONITOR_FIELD_SESSION_NUM]
		tmpMap["queue_num"] = v[constants.CHANNEL_MONITOR_FIELD_QUEUE_NUM]
		tmpMap["invalid_session_num"] = v[constants.CHANNEL_MONITOR_FIELD_INVALID_SESSION_NUM]
		tmpMap["end_session_num"] = v[constants.CHANNEL_MONITOR_FIELD_END_SESSION_NUM]
		tmpMap["serve_user_num"] = v[constants.CHANNEL_MONITOR_FIELD_SERVE_USER_NUM]
		tmpMap["giveup_queue_num"] = v[constants.CHANNEL_MONITOR_FIELD_GIVEUP_QUEUE_NUM]
		tmpMap["add_msg_num"] = v[constants.CHANNEL_MONITOR_FIELD_ADD_MSG_NUM]
		tmpMap["deal_msg_num"] = v[constants.CHANNEL_MONITOR_FIELD_DEAL_MSG_NUM]
	}

	log4go.Info("ChannelMonitor suc return. req:%+v, res:%+v:", req, res)
	context.JSON(http.StatusOK, res)
}






