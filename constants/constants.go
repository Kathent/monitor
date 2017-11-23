package constants

//Redis
//坐席监控表常量
const (
	//hash key  monitor:agent:{YMD}:{vccId}:{agentId}
	AGENT_MONITOR_HASH_KEY = "monitor:agent:%s:%s:%s"

	AGENT_MONITOR_FIELD_VCCID    = "vccId"
	AGENT_MONITOR_FIELD_AGENTID  = "ag_id"
	AGENT_MONITOR_FIELD_NAME     = "name"
	AGENT_MONITOR_FIELD_WORKER_ID     = "worker_id"
	AGENT_MONITOR_FIELD_DEP_ID     = "dep_id"
	AGENT_MONITOR_FIELD_STATUS     = "status"
	AGENT_MONITOR_FIELD_STATUS_START_TIME     = "status_start_time"
	AGENT_MONITOR_FIELD_CUR_SESSION_NUM     = "cur_session_num"
	AGENT_MONITOR_FIELD_MAX_SESSION_NUM     = "max_session_num"
	AGENT_MONITOR_FIELD_INVALID_SESSION_NUM     = "invalid_session_num"
	AGENT_MONITOR_FIELD_INDEP_SESSION_NUM     = "indep_session_num"
	AGENT_MONITOR_FIELD_DEP_SESSION_NUM     = "dep_session_num"
	AGENT_MONITOR_FIELD_TOTAL_SESSION_NUM     = "total_session_num"
	AGENT_MONITOR_FIELD_INVALID_RATIO     = "invalid_ratio"
	AGENT_MONITOR_FIELD_FIRST_RESP_TIME     = "first_resp_time"
	AGENT_MONITOR_FIELD_AVG_FIRST_RESP_TIME     = "avg_first_resp_time"
	AGENT_MONITOR_FIELD_SESSION_TIME_TOTAL     = "session_time_total"
	AGENT_MONITOR_FIELD_AVG_SESSION_TIME     = "avg_session_time"
	AGENT_MONITOR_FIELD_ONLINE_TIME_TOTAL     = "online_time_total"
	AGENT_MONITOR_FIELD_BUSY_TIME_TOTAL     = "busy_time_total"
	AGENT_MONITOR_FIELD_GROUP_IDS     = "group_ids"
	AGENT_MONITOR_FIELD_TRANSFER_IN_TIMES     = "transfer_in_times"
	AGENT_MONITOR_FIELD_TRANSFER_OUT_TIMES     = "transfer_out_times"
	AGENT_MONITOR_FIELD_REPLY_MSG_TIMES     = "reply_msg_times"
	AGENT_MONITOR_FIELD_RECEIVE_MSG_TIMES     = "receive_msg_times"
	AGENT_MONITOR_FIELD_AR_RATIO     = "ar_ratio"
	AGENT_MONITOR_FIELD_REQUIRE_EVAL_TIMES     = "require_eval_times"
	AGENT_MONITOR_FIELD_RECEIVE_EVAL_TIMES     = "receive_eval_times"
	AGENT_MONITOR_FIELD_EVAL_RATIO     = "eval_ratio"
	AGENT_MONITOR_FIELD_SERV_USER_NUM     = "serv_user_num"
	AGENT_MONITOR_FIELD_ONE_STAR_NUM     = "one_star_num"
	AGENT_MONITOR_FIELD_TWO_STAR_NUM     = "two_star_num"
	AGENT_MONITOR_FIELD_THREE_STAR_NUM     = "three_star_num"
	AGENT_MONITOR_FIELD_FOUR_STAR_NUM     = "four_star_num"
	AGENT_MONITOR_FIELD_FIVE_STAR_NUM     = "five_star_num"
	AGENT_MONITOR_FIELD_SIX_STAR_NUM     = "six_star_num"
	AGENT_MONITOR_FIELD_ONE_SERV_CLIENT_NUM     = "one_serv_client_num"
)

//渠道监控常量
const (
	//渠道监控hash key  monitor:channel:{YMD}:{vccId}:{channelId}
	CHANNEL_MONITOR_HASH_KEY = "monitor:channel:%s:%s:%s"

	CHANNEL_MONITOR_FIELD_CHANNEL_ID    = 		"channelId"
	CHANNEL_MONITOR_FIELD_VCCID 		=   	"vccId"
	CHANNEL_MONITOR_FIELD_SOURCE_TYPE   =   	"source_type"
	CHANNEL_MONITOR_FIELD_SOURCE_NAME   =   	"source_name"
	CHANNEL_MONITOR_FIELD_SESSION_NUM   =   	"session_num"
	CHANNEL_MONITOR_FIELD_QUEUE_NUM   =   	"queue_num"
	CHANNEL_MONITOR_FIELD_INVALID_SESSION_NUM   =   	"invalid_session_num"
	CHANNEL_MONITOR_FIELD_END_SESSION_NUM   =   	"end_session_num"
	CHANNEL_MONITOR_FIELD_SERVE_USER_NUM   =   	"serve_user_num"
	CHANNEL_MONITOR_FIELD_GIVEUP_QUEUE_NUM   =   	"giveup_queue_num"
	CHANNEL_MONITOR_FIELD_ADD_MSG_NUM   =   	"add_msg_num"
	CHANNEL_MONITOR_FIELD_DEAL_MSG_NUM   =   	"deal_msg_num"
	CHANNEL_MONITOR_FIELD_FIRST_RESP_TOTAL   =   	"first_resp_total"
	CHANNEL_MONITOR_FIELD_AVG_FIRST_RESP_TOTAL   =   	"avg_first_resp_total"
	CHANNEL_MONITOR_FIELD_SESSION_TIME_TOTAL   =   	"session_time_total"
	CHANNEL_MONITOR_FIELD_AVG_SESSION_TIME_TOTAL   =   	"avg_session_time_total"
	CHANNEL_MONITOR_FIELD_C_SESSION_NUM   =   	"c_session_num"
	CHANNEL_MONITOR_FIELD_EVALUATE_NUM   =   	"evaluate_num"
	CHANNEL_MONITOR_FIELD_RATE_EVALUATE   =   	"rate_evaluate"
	CHANNEL_MONITOR_FIELD_ONE_STAR_NUM   =   	"one_star_num"
	CHANNEL_MONITOR_FIELD_TWO_STAR_NUM   =   	"two_star_num"
	CHANNEL_MONITOR_FIELD_THREE_STAR_NUM   =   	"three_star_num"
	CHANNEL_MONITOR_FIELD_FOUR_STAR_NUM   =   	"four_star_num"
	CHANNEL_MONITOR_FIELD_FIVE_STAR_NUM   =   	"five_star_num"
	CHANNEL_MONITOR_FIELD_SIX_STAR_NUM   =   	"six_star_num"
)

//坐席工号有序集合常量
const(
	//坐席工号有序集合 monitor:agentSSet:{vccId}
	AGENT_WORK_ID_SORT_SET = "monitor:agentSSet:%s"
)

//已评价父会话
const(
	//已评价父会话集合 monitor:evalCid:{vccId}
	EVALUATED_CID_SET = "monitor:evalCid:%s"
)

//坐席服务客户set
const(
	//坐席服务客户set monitor:serve:{YMD}::{vccId}:{agentId}
	AGENT_SERVE_CLIENT_SET = "monitor:serve:%s::%s:%s"
)

//客户对应父会话set
const(
	//客户对应父会话set monitor:client_session:{YMD}:{vccId}:{userId}
	CLIENT_CID_SET = "monitor:client_session:%s:%s:%s"
)


//mongodb
//工作量统计及绩效
const(
	//工作量统计及绩效mongodb表
	STATICS_AW_TABLE_NAME = "im_report_agent_work"

	STATICS_AW_FILED_VCC_ID     	= "vcc_id"
	STATICS_AW_FILED_AG_ID     	= "ag_id"
	STATICS_AW_FILED_GROUP_IDS     	= "group_ids"
	STATICS_AW_FILED_WORKER_ID     	= "worker_id"
	STATICS_AW_FILED_DEPT_ID     	= "dept_id"
	STATICS_AW_FILED_DATE     	= "date"
	STATICS_AW_FILED_ONLINE_SECS     	= "online_secs"
	STATICS_AW_FILED_BUSY_SECS     	= "busy_secs"
	STATICS_AW_FILED_IN_SESSION_NUM     	= "in_session_num"
	STATICS_AW_FILED_INVALID_SESSION_NUM     	= "invalid_session_num"
	STATICS_AW_FILED_RECEIVE_SESSION_NUM     	= "receive_session_num"
	STATICS_AW_FILED_TRANS_IN_SESSION_NUM     	= "trans_in_session_num"
	STATICS_AW_FILED_TRANS_OUT_SESSION_NUM     	= "trans_out_session_num"
	STATICS_AW_FILED_REPLY_NEWS_NUM     	= "reply_news_num"
	STATICS_AW_FILED_RECEIVE_MSG_TIMES     	= "receive_msg_times"
	STATICS_AW_FILED_AR_RATIO     	= "ar_ratio"
	STATICS_AW_FILED_REQUIRE_EVAL_TIMES     	= "require_eval_times"
	STATICS_AW_FILED_TOTAL_SESSION_NUM     	= "total_session_num"
	STATICS_AW_FILED_INVALID_RATIO     	= "invalid_ratio"
	STATICS_AW_FILED_FIRST_RESPONSE_SECS     	= "first_response_secs"
	STATICS_AW_FILED_AVG_RESPONSE_SECS     	= "avg_response_secs"
	STATICS_AW_FILED_TOTAL_SESSION_SECS     	= "total_session_secs"
	STATICS_AW_FILED_AVG_SESSION_SECS     	= "avg_session_secs"
	STATICS_AW_FILED_SERV_USER_NUM     	= "serv_user_num"
	STATICS_AW_FILED_ONE_SERV_CLIENT_NUM     	= "one_serv_client_num"
	STATICS_AW_FILED_RECEIVE_EVAL_TIMES     	= "receive_eval_times"
	STATICS_AW_FILED_EVAL_RATIO     	= "eval_ratio"
)

//坐席满意度mongodb表
const(
	STATICS_AE_TABLE_NAME 					= "im_report_evaluate_agent"

	STATICS_AE_FIELD_VCC_ID 					= "vcc_id"
	STATICS_AE_FIELD_AG_ID 					= "ag_id"
	STATICS_AE_FIELD_WORKER_ID 					= "worker_id"
	STATICS_AE_FIELD_DEP_ID 					= "dep_id"
	STATICS_AE_FIELD_GROUP_IDS 					= "group_ids"
	STATICS_AE_FIELD_INDEP_SESSION_NUM 					= "indep_session_num"
	STATICS_AE_FIELD_REQUIRE_EVAL_TIMES 					= "require_eval_times"
	STATICS_AE_FIELD_EVALUATE_NUM 					= "evaluate_num"
	STATICS_AE_FIELD_RATE_EVALUATE 					= "rate_evaluate"
	STATICS_AE_FIELD_ONE_STAR_NUM 					= "one_star_num"
	STATICS_AE_FIELD_TWO_STAR_NUM 					= "two_star_num"
	STATICS_AE_FIELD_THREE_STAR_NUM 					= "three_star_num"
	STATICS_AE_FIELD_FOUR_STAR_NUM 					= "four_star_num"
	STATICS_AE_FIELD_FIVE_STAR_NUM 					= "five_star_num"
	STATICS_AE_FIELD_SIX_STAR_NUM 					= "six_star_num"
)

//坐席状态记录
const(
	STATICS_AS_TABLE_NAME 				= "im_report_agent_status"

	STATICS_AS_FIELD_VCC_ID				= "vcc_id"
	STATICS_AS_FIELD_AG_ID				= "ag_id"
	STATICS_AS_FIELD_WORKER_ID				= "worker_id"
	STATICS_AS_FIELD_DEPT_ID				= "dept_id"
	STATICS_AS_FIELD_DATE				= "date"
	STATICS_AS_FIELD_PRE_STATUS				= "pre_status"
	STATICS_AS_FIELD_STATUS				= "status"
	STATICS_AS_FIELD_OP_TYPE				= "op_type"
	STATICS_AS_FIELD_PRE_STATUS_SECS				= "pre_status_secs"
	STATICS_AS_FIELD_TIME				= "time"
)

//会话记录
const(
	STATICS_SR_TABLE_NAME 					= "im_report_session_record"

	STATICS_SR_FIELD_SID					= "sid"
	STATICS_SR_FIELD_C_ID					= "c_id"
	STATICS_SR_FIELD_DATE					= "date"
	STATICS_SR_FIELD_VCC_ID					= "vcc_id"
	STATICS_SR_FIELD_NAME					= "name"
	STATICS_SR_FIELD_AG_ID					= "ag_id"
	STATICS_SR_FIELD_WORKER_ID					= "worker_id"
	STATICS_SR_FIELD_CLIENT_NEWS_NUM					= "client_news_num"
	STATICS_SR_FIELD_AGENT_NEWS_NUM					= "agent_news_num"
	STATICS_SR_FIELD_SESSION_START_TIME					= "session_start_time"
	STATICS_SR_FIELD_SESSION_END_TIME					= "session_end_time"
	STATICS_SR_FIELD_FIRST_RESPONSE_SECS					= "first_response_secs"
	STATICS_SR_FIELD_SESSION_KEEP_SECS					= "session_keep_secs"
	STATICS_SR_FIELD_SESSION_CREATE_TYPE					= "session_create_type"
	STATICS_SR_FIELD_SESSION_END_TYPE					= "session_end_type"
	STATICS_SR_FIELD_SOURCE_TYPE					= "source_type"
	STATICS_SR_FIELD_IS_EVALUATE					= "is_evaluate"
	STATICS_SR_FIELD_EVALUATE					= "evaluate"
	STATICS_SR_FIELD_EVALUATE_EXPLAIN					= "evaluate_explain"
)

//会话内容记录
const(
	STATICS_SC_TABLE_NAME 	= "im_session_content"

	STATICS_SC_FIELD_SID   =   			"sid"
	STATICS_SC_FIELD_INDEX   =   			"index"
	STATICS_SC_FIELD_TYPE   =   			"type"
	STATICS_SC_FIELD_CONTENT   =   			"content"
)

//渠道满意度报表
const(
	STATICS_CE_TABLE_NAME 				= "im_report_evaluate_source"

	STATICS_CE_FIELD_VCC_ID				= "vcc_id"
	STATICS_CE_FIELD_CHANNEL_ID				= "channel_id"
	STATICS_CE_FIELD_SOURCE_NAME				= "source_name"
	STATICS_CE_FIELD_SOURCE_TYPE				= "source_type"
	STATICS_CE_FIELD_END_SESSION_NUM				= "end_session_num"
	STATICS_CE_FIELD_EVALUATE_NUM				= "evaluate_num"
	STATICS_CE_FIELD_RATE_EVALUATE				= "rate_evaluate"
	STATICS_CE_FIELD_ONE_STAR_NUM				= "one_star_num"
	STATICS_CE_FIELD_TWO_STAR_NUM				= "two_star_num"
	STATICS_CE_FIELD_THREE_STAR_NUM				= "three_star_num"
	STATICS_CE_FIELD_FOUR_STAR_NUM				= "four_star_num"
	STATICS_CE_FIELD_FIVE_STAR_NUM				= "five_star_num"
	STATICS_CE_FIELD_SIX_STAR_NUM				= "six_star_num"
)

const (
	DATE_FORMATE = "20060102"
	DATE_FORMATE_ALL = "2006-01-02 15:04:05"

	EVENT_TYPE_SESSION_END = "sessionEnd"
	EVENT_TYPE_SESSION_CONTENT = "sessionContent"

	SORT_ORDER_DEFAULT = "ascend"
	SORT_ORDER_DESCEND = "descend"

	SORT_ORDER_DESCEND_SYM = -1
	SORT_ORDER_ASCEND_SYM = 1

	SESSION_CREATE_TYPE_DIRECT = "direct"
	SESSION_CREATE_TYPE_DIRECT_INT = "0"
	SESSION_CREATE_TYPE_REVERT = "revert"
	SESSION_CREATE_TYPE_REVERT_INT = "1"
	SESSION_CREATE_TYPE_MESSAGE = "message"
	SESSION_CREATE_TYPE_MESSAGE_INT = "2"
)
