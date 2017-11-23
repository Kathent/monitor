package mq

type AgentStatus struct {
	GroupId string `bson:"group_ids"`
	VccId string `bson:"vcc_id"`
	AgentId string `bson:"ag_id"`
	WorkerId string `bson:"worker_id"`
	DeptId string `bson:"dept_id"`
	Date string	`bson:"date"`
	PreStatus string `bson:"pre_status"`
	Status string	`bson:"status"`
	OpType string	`bson:"action"`
	PreStatusSecs int64 `bson:"pre_status_secs"`
	Time string	`bson:"time"`
}

type SessionRecord struct {
	Sid string `bson:"sid"`
	Cid string	`bson:"c_id"`
	Date string `bson:"date"`
 	VccId string `bson:"vcc_id"`
	Name string `bson:"name"`
	AgentId string `bson:"ag_id"`
	WorkerId string `bson:"worker_id"`
	ClientNewsNum int `bson:"client_news_num"`
	AgentNewsNum int `bson:"agent_news_num"`
	SessionStartTime string `bson:"session_start_time"`
	SessionEndTime string `bson:"session_end_time"`
	FirstRespSecs int64 `bson:"first_response_secs"`
	SessionKeepSecs int64 `bson:"session_keep_secs"`
	CreateType string `bson:"session_create_type"`
	EndType string `bson:"session_end_type"`
	SourceType int8 `bson:"source_type"`
	SourceName string `bson:"source_name"`
	IsEvaluate int8 `bson:"is_evaluate"`
	Evaluate string `bson:"evaluate"`
	EvalExplain string `bson:"evaluate_explain"`
}

type SessionContent struct {
	SessionId string `bson:"sid"`
	Index int `bson:"index"`
	Type string `bson:"type"`
	Content string `bson:"content"`
}

type AgentWorkMG struct {
	VccId string `bson:"vcc_id"`
	AgentId string `bson:"ag_id"`
	DeptId string `bson:"dept_id"`
	GroupId string  `bson:"group_ids"`
	WorkerId string  `bson:"worker_id"`
	Date string `bson:"date"`
	OnlineSecs int64 `bson:"online_secs"`
	BusySecs int64 `bson:"busy_secs"`
	InvalidSessionNum int  `bson:"invalid_session_num"`
	ReceiveSessionNUm int  `bson:"receive_session_num"`
	TransInSessionNum int  `bson:"trans_in_session_num"`
	TransOutSessionNum int  `bson:"trans_out_session_num"`
	ReplyNewsNum int  `bson:"reply_news_num"`
	ReceiveMsgTimes int `bson:"receive_msg_times"`
	ArRatio float64 `bson:"ar_ratio"`
	InSessionNum int `bson:"in_session_num"`
	RequireEvalTimes int `bson:"require_eval_times"`
	TotalSessionNum int  `bson:"total_session_num"`
	InvalidRatio float64  `bson:"invalid_ratio"`
	AvgResponseSecs int64  `bson:"avg_response_secs"`
	AvgSessionSecs int64  `bson:"avg_session_secs"`
	ServeUserNum int  `bson:"serv_user_num"`
	OneServeClientNum int  `bson:"one_serv_client_num"`
	ReceiveEvalTimes int  `bson:"receive_eval_times"`
	EvalRatio float64  `bson:"eval_ratio"`
	FirstRespSecs int64 `bson:"first_response_secs"`
	SessionKeepSecs int64 `bson:"session_keep_secs"`
}

type AgentWork struct {
	GroupId string  `bson:"group_ids"`
	WorkerId string  `bson:"worker_id"`
	OnlineSecs int64 `bson:"online_secs"`
	BusySecs int64 `bson:"busy_secs"`
	InSessionNum int `bson:"in_session_num"`
	InvalidSessionNum int  `bson:"invalid_session_num"`
	ReceiveSessionNUm int  `bson:"receive_session_num"`
	TransInSessionNum int  `bson:"trans_in_session_num"`
	TransOutSessionNum int  `bson:"trans_out_session_num"`
	ReplyNewsNum int  `bson:"reply_news_num"`
}

type AgentAchievements struct {
	GroupId string  `bson:"group_ids"`
	WorkerId string  `bson:"worker_id"`
	ReceiveMsgTimes int64 `bson:"receive_msg_times"`
	ReplyNewsNum int64 `bson:"reply_news_num"`
	ArRatio float64 `bson:"ar_ratio"`
	InSessionNum int `bson:"in_session_num"`
	RequireEvalTimes int `bson:"require_eval_times"`
	TotalSessionNum int  `bson:"total_session_num"`
	InvalidSessionNum int  `bson:"invalid_session_num"`
	InvalidRatio float64  `bson:"invalid_ratio"`
	AvgResponseSecs int64  `bson:"avg_response_secs"`
	AvgSessionSecs int64  `bson:"avg_session_secs"`
	ServeUserNum int  `bson:"serv_user_num"`
	OneServeClientNum int  `bson:"one_serv_client_num"`
	ReceiveEvalTimes int  `bson:"receive_eval_times"`
	EvalRatio float64  `bson:"eval_ratio"`
	FirstRespSecs int64 `bson:"first_response_secs"`
	SessionKeepSecs int64 `bson:"session_keep_secs"`
}

type AgentEval struct {
	VccId string `bson:"vcc_id"`
	AgentId string `bson:"ag_id"`
	DepId string `bson:"dep_id"`
	GroupId string  `bson:"group_ids"`
	WorkerId string  `bson:"worker_id"`
	IndepSessionNum int  `bson:"indep_session_num"`
	RequireEvalTimes int  `bson:"require_eval_times"`
	EvaluateNum int  `bson:"evaluate_num"`
	RateEvaluate float64 `bson:"rate_evaluate"`
	OneStarNum int `bson:"one_star_num"`
	TwoStarNum int `bson:"two_star_num"`
	ThreeStarNum int  `bson:"three_star_num"`
	FourStarNum  int `bson:"four_star_num"`
	FiveStarNum int `bson:"five_star_num"`
	SixStarNum int `bson:"six_star_num"`
	Date string `bson:date`
}

type ChannelEval struct {
	VccId string `bson:"vcc_id"`
	ChannelId string `bson:"channel_id"`
	SourceName string `bson:"source_name"`
	SourceType string `bson:"source_type"`
	EndSessionNum int  `bson:"end_session_num"`
	EvaluateNum int  `bson:"evaluate_num"`
	RateEvaluate float64 `bson:"rate_evaluate"`
	OneStarNum int `bson:"one_star_num"`
	TwoStarNum int `bson:"two_star_num"`
	ThreeStarNum int  `bson:"three_star_num"`
	FourStarNum  int `bson:"four_star_num"`
	FiveStarNum int `bson:"five_star_num"`
	SixStarNum int `bson:"six_star_num"`
	Date string `bson:date`
}