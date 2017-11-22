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

//工作量统计
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
	if !util.IsEmpty(req.AgentId) {
		split := strings.Split(req.AgentId, ",")
		pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End},
			"ag_id": bson.M{"$in": split},},},
			{"$group":	bson.M{"_id": "$worker_id", "group_ids": bson.M{"$first": "$group_ids"},
				"worker_id": bson.M{"$first":"$worker_id"},
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
		pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End}}},
			{"$group":	bson.M{"_id": "$worker_id", "group_ids": bson.M{"$first": "$group_ids"},
				"worker_id": bson.M{"$first":"$worker_id"},
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
	allErr := pipe.All(&arr)
	if allErr != nil {
		log4go.Warn("allErr:%v", allErr)
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

	res.Data.Rows = make([]interface{}, 0)
	var tmpMap map[string]interface{}
	for _, v := range arr {
		tmpMap = make(map[string]interface{})
		tmpMap["group_ids"] = util.GetString(v.GroupId)
		tmpMap["worker_id"] = util.GetString(v.WorkerId)
		tmpMap["online_secs"] = util.GetIntString(v.OnlineSecs)
		tmpMap["busy_secs"] = util.GetIntString(v.BusySecs)
		tmpMap["in_session_num"] = util.GetIntString(v.InSessionNum)
		tmpMap["invalid_session_num"] = util.GetIntString(v.InvalidSessionNum)
		tmpMap["receive_session_num"] = util.GetIntString(v.ReceiveSessionNUm)
		tmpMap["trans_in_session_num"] = util.GetIntString(v.TransInSessionNum)
		tmpMap["trans_out_session_num"] = util.GetIntString(v.TransOutSessionNum)
		tmpMap["reply_news_num"] = util.GetIntString(v.ReplyNewsNum)
		res.Data.Rows = append(res.Data.Rows, tmpMap)
	}

	context.JSON(http.StatusOK, res)
	log4go.Info("AgentWorkStatics suc return. req:%+v, res:%+v:", agentReq, res)
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
		ar.Start = mq.TransTime(now, constants.DATE_FORMATE)
	}else {
		start, err = time.Parse(constants.DATE_FORMATE, req.Start)
		ar.Start = req.Start
	}

	if err != nil {
		return nil, errors.New("start is wrong")
	}

	if util.IsEmpty(req.End){
		end = time.Now()
		ar.End = mq.TransTime(end, constants.DATE_FORMATE)
	}else {
		end, err = time.Parse(constants.DATE_FORMATE, req.End)
		ar.End = req.End
	}

	if start.After(end) {
		return nil, errors.New("start after end")
	}

	if err != nil {
		return nil, errors.New("end is wrong")
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	ar.PageSize = req.PageSize
	ar.AgentId = req.AgentId

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
	}else {
		ar.SortField = req.SortField
	}

	if req.SortOrder != constants.SORT_ORDER_DESCEND{
		ar.SortOrder = constants.SORT_ORDER_DEFAULT
	}

	ar.Status = req.Status
	ar.DepId = req.DepId
	return &ar, nil
}

//绩效统计
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
	if !util.IsEmpty(req.AgentId) {
		split := strings.Split(req.AgentId, ",")
		pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End},
			"agentId": bson.M{"$in": split},}},
			{"$group":	bson.M{"_id": "$worker_id",  "group_ids": bson.M{"$first": "$group_ids"},
				"worker_id": bson.M{"$first": "$worker_id"},
				"receive_msg_times": bson.M{"$sum":"$receive_msg_times"},
				"reply_news_num": bson.M{"$sum": "$reply_news_num"},
				"in_session_num": bson.M{"$sum": "$in_session_num"},
				"require_eval_times": bson.M{"$sum": "$require_eval_times"},
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
		pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End}}},
			{"$group":	bson.M{"_id": "$worker_id", "group_ids": bson.M{"$first": "$group_ids"},
				"worker_id": bson.M{"$first": "$worker_id"},
				"receive_msg_times": bson.M{"$sum":"$receive_msg_times"},
				"reply_news_num": bson.M{"$sum": "$reply_news_num"},
				"in_session_num": bson.M{"$sum": "$in_session_num"},
				"require_eval_times": bson.M{"$sum": "$require_eval_times"},
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
	res.Data.Total = getPageTotal(total, agentReq.PageSize)
	res.Code = 0
	res.Message = "suc"

	if agentReq.StartIndex < total {
		if agentReq.EndIndex >= total {
			arr = arr[agentReq.StartIndex:]
		} else {
			arr = arr[agentReq.StartIndex: agentReq.EndIndex]
		}
	}

	res.Data.Rows = make([]interface{}, 0)
	var tmpMap map[string]interface{}
	for _, v := range arr {
		v.AvgResponseSecs = v.FirstRespSecs / int64(util.GetDefaultInt(v.TotalSessionNum, 1))
		v.AvgSessionSecs = v.SessionKeepSecs / int64(util.GetDefaultInt(v.TotalSessionNum, 1))
		v.ArRatio = float64(v.ReceiveMsgTimes) / float64(util.GetDefaultInt64(v.ReplyNewsNum, 1))
		v.EvalRatio = float64(v.ReceiveEvalTimes) / float64(util.GetDefaultInt(v.InSessionNum, 1))
		v.InvalidRatio = float64(v.InvalidSessionNum) / float64(util.GetDefaultInt(v.TotalSessionNum, 1))

		tmpMap = make(map[string]interface{})

		tmpMap["worker_id"] = util.GetString(v.WorkerId)
		tmpMap["group_id"] = util.GetString(v.GroupId)
		tmpMap["receive_msg_times"] = util.GetIntString(v.ReceiveMsgTimes)
		tmpMap["reply_news_num"] = util.GetIntString(v.ReplyNewsNum)
		tmpMap["ar_ratio"] = fmt.Sprintf("%.2f", v.ArRatio)
		tmpMap["require_eval_times"] = util.GetIntString(v.RequireEvalTimes)
		tmpMap["total_session_num"] = util.GetIntString(v.TotalSessionNum)
		tmpMap["invalid_session_num"] = util.GetIntString(v.InvalidSessionNum)
		tmpMap["invalid_ratio"] = fmt.Sprintf("%.2f", v.InvalidRatio)
		tmpMap["avg_response_secs"] = util.GetIntString(v.AvgResponseSecs)
		tmpMap["avg_session_secs"] = util.GetIntString(v.AvgSessionSecs)
		tmpMap["serv_user_num"] = util.GetIntString(v.ServeUserNum)
		tmpMap["one_serv_client_num"] = util.GetIntString(v.OneServeClientNum)
		tmpMap["receive_eval_times"] = util.GetIntString(v.ReceiveEvalTimes)
		tmpMap["eval_ratio"] = fmt.Sprintf("%.2f", v.EvalRatio)
		res.Data.Rows = append(res.Data.Rows, tmpMap)
	}

	context.JSON(http.StatusOK, res)
	log4go.Info("AgentAchievements suc return. req:%+v, res:%+v:", agentReq, res)
	return
}
func getPageTotal(total int, pageSize int) int {
	if total % pageSize == 0 {
		return total / pageSize
	}else {
		return total / pageSize + 1
	}
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
		err = collection.Find(bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End}}).
			Sort(agentReq.GetSortString()).All(&arr)
	}else {
		split := strings.Split(req.AgentId, ",")
		err = collection.Find(bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End},
			"agentId": bson.M{"$in": split},}).
			Sort(agentReq.GetSortString()).All(&arr)
	}

	if err != nil {
		log4go.Warn("AgentSessionRecord find err:%v", err)
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
	log4go.Info("AgentStatusRecord suc return. req:%+v, res:%+v:", agentReq, res)
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

	collection := db.GetSession().DB("").C(constants.STATICS_AE_TABLE_NAME)
	var pipe *mgo.Pipe
	if !util.IsEmpty(req.AgentId) {
		split := strings.Split(req.AgentId, ",")
		pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vcc_id": agentReq.VccId,
			"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End},
			"ag_id": bson.M{"$in": split},}},
			{"$group":	bson.M{"_id": "$worker_id", "group_ids": bson.M{"$first": "$group_ids"},
				"worker_id": bson.M{"$first": "$worker_id"},
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
			"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End}}},
			{"$group":	bson.M{"_id": "$worker_id", "group_ids": bson.M{"$first": "$group_ids"},
				"worker_id": bson.M{"$first": "$worker_id"},
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
	allErr := pipe.All(&arr)
	if allErr != nil {
		log4go.Warn("allErr:%v", allErr)
	}

	total := len(arr)
	res.Data.Total = getPageTotal(total, agentReq.PageSize)
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
		v.RateEvaluate = float64(v.EvaluateNum) / float64(util.GetDefaultInt(v.IndepSessionNum, 1))
		tmpMap = make(map[string]interface{})
		tmpMap["group_ids"] = v.GroupId
		tmpMap["worker_id"] = v.WorkerId
		tmpMap["indep_session_num"] = util.GetIntString(v.IndepSessionNum)
		tmpMap["require_eval_times"] = util.GetIntString(v.RequireEvalTimes)
		tmpMap["evaluate_num"] = util.GetIntString(v.RequireEvalTimes)
		tmpMap["one_star_num"] = util.GetIntString(v.OneStarNum)
		tmpMap["two_star_num"] = util.GetIntString(v.TwoStarNum)
		tmpMap["three_star_num"] = util.GetIntString(v.ThreeStarNum)
		tmpMap["four_star_num"] = util.GetIntString(v.FourStarNum)
		tmpMap["five_star_num"] = util.GetIntString(v.FiveStarNum)
		tmpMap["six_star_num"] = util.GetIntString(v.SixStarNum)
		tmpMap["rate_evaluate"] = fmt.Sprintf("%.2f", v.RateEvaluate)
		res.Data.Rows = append(res.Data.Rows, tmpMap)
	}

	context.JSON(http.StatusOK, res)
	log4go.Info("AgentEval suc return. req:%+v, res:%+v:", agentReq, res)
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

	collection := db.GetSession().DB("").C(constants.STATICS_CE_TABLE_NAME)
	var pipe *mgo.Pipe
	pipe = collection.Pipe([]bson.M{{"$match": bson.M{"vcc_id": agentReq.VccId,
		"date": bson.M{"$gte": agentReq.Start, "$lte": agentReq.End}}},
		{"$group":	bson.M{"_id": "$channel_id", "group_ids": bson.M{"$first": "%group_ids"},
			"source_name": bson.M{"$first": "$source_name"},
			"source_type": bson.M{"$first": "$source_type"},
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
		v.RateEvaluate = float64(v.EvaluateNum) / float64(util.GetDefaultInt(v.EndSessionNum, 1))

		tmpMap = make(map[string]interface{})
		tmpMap["source_name"] = v.SourceName
		tmpMap["source_type"] = v.SourceType
		tmpMap["end_session_num"] = util.GetIntString(v.EndSessionNum)
		tmpMap["evaluate_num"] = util.GetIntString(v.EvaluateNum)
		tmpMap["one_star_num"] = util.GetIntString(v.OneStarNum)
		tmpMap["two_star_num"] = util.GetIntString(v.TwoStarNum)
		tmpMap["three_star_num"] = util.GetIntString(v.ThreeStarNum)
		tmpMap["four_star_num"] = util.GetIntString(v.FourStarNum)
		tmpMap["five_star_num"] = util.GetIntString(v.FiveStarNum)
		tmpMap["six_star_num"] = util.GetIntString(v.SixStarNum)
		tmpMap["rate_evaluate"] = fmt.Sprintf("%.2f", v.RateEvaluate)
		res.Data.Rows = append(res.Data.Rows, tmpMap)
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
	keys, err := db.GetClient().Keys(key).Result()
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
			log4go.Warn("AgentMonitor hGetAll err:%v, key:%s", err, v)
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

	now := time.Now().Unix()
	var tmpMap map[string]interface{}
	for _, v := range arr {
		tmpMap = make(map[string]interface{})
		tmpMap["name"] = v[constants.AGENT_MONITOR_FIELD_NAME]

		depId := v[constants.AGENT_MONITOR_FIELD_DEP_ID]
		tmpMap["dep_id"] = v[constants.AGENT_MONITOR_FIELD_DEP_ID]
		if !util.IsEmpty(agentReq.DepId) && agentReq.DepId != depId {
			continue
		}
		tmpMap["worker_id"] = v[constants.AGENT_MONITOR_FIELD_WORKER_ID]
		status := v[constants.AGENT_MONITOR_FIELD_STATUS]
		if !util.IsEmpty(agentReq.Status) && agentReq.Status != status {
			continue
		}
		tmpMap["status"] = v[constants.AGENT_MONITOR_FIELD_STATUS]
		statusStartTime := int64(util.GetInt(v[constants.AGENT_MONITOR_FIELD_STATUS_START_TIME]))
		tmpMap["status_keep_time"] = fmt.Sprintf("%d", now - statusStartTime)
		tmpMap["cur_session_num"] = util.GetStringDefault(v[constants.AGENT_MONITOR_FIELD_CUR_SESSION_NUM], "0")
		tmpMap["max_session_num"] = util.GetStringDefault(v[constants.AGENT_MONITOR_FIELD_MAX_SESSION_NUM], "0")
		tmpMap["invalid_session_num"] = util.GetStringDefault(v[constants.AGENT_MONITOR_FIELD_INVALID_SESSION_NUM], "0")
		tmpMap["indep_session_num"] = util.GetStringDefault(v[constants.AGENT_MONITOR_FIELD_INDEP_SESSION_NUM], "0")
		tmpMap["dep_session_num"] = util.GetStringDefault(v[constants.AGENT_MONITOR_FIELD_DEP_SESSION_NUM], "0")
		res.Data.Rows = append(res.Data.Rows, tmpMap)
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
			log4go.Warn("ChannelMonitor HGetAll err:%v, key:%s", err, v)
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
		tmpMap["session_num"] = util.GetStringDefault(v[constants.CHANNEL_MONITOR_FIELD_SESSION_NUM], "0")
		tmpMap["queue_num"] = util.GetStringDefault(v[constants.CHANNEL_MONITOR_FIELD_QUEUE_NUM], "0")
		tmpMap["invalid_session_num"] = util.GetStringDefault(v[constants.CHANNEL_MONITOR_FIELD_INVALID_SESSION_NUM], "0")
		tmpMap["end_session_num"] = util.GetStringDefault(v[constants.CHANNEL_MONITOR_FIELD_END_SESSION_NUM], "0")
		tmpMap["serve_user_num"] = util.GetStringDefault(v[constants.CHANNEL_MONITOR_FIELD_SERVE_USER_NUM], "0")
		tmpMap["giveup_queue_num"] = util.GetStringDefault(v[constants.CHANNEL_MONITOR_FIELD_GIVEUP_QUEUE_NUM], "0")
		tmpMap["add_msg_num"] = util.GetStringDefault(v[constants.CHANNEL_MONITOR_FIELD_ADD_MSG_NUM], "0")
		tmpMap["deal_msg_num"] = util.GetStringDefault(v[constants.CHANNEL_MONITOR_FIELD_DEAL_MSG_NUM], "0")
		res.Data.Rows = append(res.Data.Rows, tmpMap)
	}

	log4go.Info("ChannelMonitor suc return. req:%+v, res:%+v:", req, res)
	context.JSON(http.StatusOK, res)
}






