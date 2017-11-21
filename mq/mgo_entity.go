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
	OpType string	`bson:"op_type"`
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
	VccId string `json:"vcc_id"`
	AgentId string `json:"ag_id"`
	DeptId string `json:"dept_id"`
	GroupId string  `json:"group_ids"`
	WorkerId string  `json:"worker_id"`
	Date string `json:"date"`
	OnlineSecs int64 `json:"online_secs"`
	BusySecs int64 `json:"busy_secs"`
	InvalidSessionNum int  `json:"invalid_session_num"`
	ReceiveSessionNUm int  `json:"receive_session_num"`
	TransInSessionNum int  `json:"trans_in_session_num"`
	TransOutSessionNum int  `json:"trans_out_session_num"`
	ReplyNewsNum int  `json:"reply_news_num"`
	ReceiveMsgTimes int `json:"receive_msg_times"`
	ArRatio float64 `json:"ar_ratio"`
	InSessionNum int `json:"in_session_num"`
	RequireEvalTimes int `json:"require_eval_times"`
	TotalSessionNum int  `json:"total_session_num"`
	InvalidRatio float64  `json:"invalid_ratio"`
	AvgResponseSecs int64  `json:"avg_response_secs"`
	AvgSessionSecs int64  `json:"avg_session_secs"`
	ServeUserNum int  `json:"serv_user_num"`
	OneServeClientNum int  `json:"one_serv_client_num"`
	ReceiveEvalTimes int  `json:"receive_eval_times"`
	EvalRatio float64  `json:"eval_ratio"`
	FirstRespSecs int64 `json:"first_response_secs"`
	SessionKeepSecs int64 `json:"session_keep_secs"`
}

type AgentWork struct {
	GroupId string  `json:"group_ids"`
	WorkerId string  `json:"worker_id"`
	OnlineSecs int64 `json:"online_secs"`
	BusySecs int64 `json:"busy_secs"`
	InSessionNum int `json:"in_session_num"`
	InvalidSessionNum int  `json:"invalid_session_num"`
	ReceiveSessionNUm int  `json:"receive_session_num"`
	TransInSessionNum int  `json:"trans_in_session_num"`
	TransOutSessionNum int  `json:"trans_out_session_num"`
	ReplyNewsNum int  `json:"reply_news_num"`
}

type AgentAchievements struct {
	GroupId string  `json:"group_ids"`
	WorkerId string  `json:"worker_id"`
	ReceiveMsgTimes int64 `json:"receive_msg_times"`
	ReplyNewsNum int64 `json:"reply_news_num"`
	ArRatio float64 `json:"ar_ratio"`
	InSessionNum int `json:"in_session_num"`
	RequireEvalTimes int `json:"require_eval_times"`
	TotalSessionNum int  `json:"total_session_num"`
	InvalidSessionNum int  `json:"invalid_session_num"`
	InvalidRatio float64  `json:"invalid_ratio"`
	AvgResponseSecs int64  `json:"avg_response_secs"`
	AvgSessionSecs int64  `json:"avg_session_secs"`
	ServeUserNum int  `json:"serv_user_num"`
	OneServeClientNum int  `json:"one_serv_client_num"`
	ReceiveEvalTimes int  `json:"receive_eval_times"`
	EvalRatio float64  `json:"eval_ratio"`
	FirstRespSecs int64 `json:"first_response_secs"`
	SessionKeepSecs int64 `json:"session_keep_secs"`
}

type AgentEval struct {
	VccId string `json:"vcc_id"`
	AgentId string `json:"ag_id"`
	DepId string `json:"dep_id"`
	GroupId string  `json:"group_ids"`
	WorkerId string  `json:"worker_id"`
	IndepSessionNum int  `json:"indep_session_num"`
	RequireEvalTimes int  `json:"require_eval_times"`
	EvaluateNum int  `json:"evaluate_num"`
	RateEvaluate float64 `json:"rate_evaluate"`
	OneStarNum int `json:"one_star_num"`
	TwoStarNum int `json:"two_star_num"`
	ThreeStarNum int  `json:"three_star_num"`
	FourStarNum  int `json:"four_star_num"`
	FiveStarNum int `json:"five_star_num"`
	SixStarNum int `json:"six_star_num"`
}

type ChannelEval struct {
	VccId string `json:"vcc_id"`
	ChannelId string `json:"channel_id"`
	SourceName string `json:"source_name"`
	SourceType string `json:"source_type"`
	EndSessionNum int  `json:"end_session_num"`
	EvaluateNum int  `json:"evaluate_num"`
	RateEvaluate float64 `json:"rate_evaluate"`
	OneStarNum int `json:"one_star_num"`
	TwoStarNum int `json:"two_star_num"`
	ThreeStarNum int  `json:"three_star_num"`
	FourStarNum  int `json:"four_star_num"`
	FiveStarNum int `json:"five_star_num"`
	SixStarNum int `json:"six_star_num"`
}