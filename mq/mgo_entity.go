package mq

type AgentStatus struct {
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
