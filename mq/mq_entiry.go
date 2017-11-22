package mq

//满意度评价消息体，发给MQ的
type EvaluateMsgMQ struct {
	EventType       string `json:"eventType"`
	ChannelID       string `json:"channelID"`
	SourceType      uint8  `json:"sourceType"`
	VccID           string `json:"vccID"`
	UserID          string `json:"userID"`
	AgentID         string `json:"agentID"`
	GroupID         string `json:"groupID"`
	ConnSuccessTime int64  `json:"connSuccessTime"`
	EvaluateTime    int64  `json:"evaluateTime"`
	SessionID       string `json:"sessionID"`
	OptionName      string `json:"optionName"`
	EvaluateExplain string `json:"evaluate_explain"`
}

//网站用户留言消息体，发给MQ的
type WebMessageMQ struct {
	EventType   string `json:"eventType"`
	ChannelID   string `json:"channelID"`
	SourceType  uint8  `json:"sourceType"`
	VccID       string `json:"vccID"`
	UserID      string `json:"userID"`
	WxNickName  string `json:"wxNickName"`
	SessionID   string `json:"sessionID"`
	MessageInfo map[string]string `json:"messageInfo"`
	MessageTime int64  `json:"messageTime"`
}

//微信用户留言消息体，发给MQ的
type WXMessageMQ struct {
	EventType   string `json:"eventType"`
	ChannelID   string `json:"channelID"`
	SourceType  uint8  `json:"sourceType"`
	VccID       string `json:"vccID"`
	UserID      string `json:"userID"`
	WxNickName  string `json:"wxNickName"`
	SessionID   string `json:"sessionID"`
	MessageInfo map[string]string `json:"messageInfo"`
	MessageTime int64  `json:"messageTime"`
}

//用户或坐席中断消息体，发给MQ的
type SessionBreakMQ struct {
	EventType  string `json:"eventType"`
	ChannelID  string `json:"channelID"`
	SourceType uint8  `json:"sourceType"`
	VccID      string `json:"vccID"`
	UserType   uint8  `json:"userType"`
	ExtID      string `json:"extID"`
	UserName   string `json:"userName"`
	SessionID  string `json:"sessionID"`
	BreakTime  int64  `json:"breakTime"`
}

//进入排队事件发送的消息体，发给MQ的
type InQueueMQ struct {
	EventType   string `json:"eventType"`
	ChannelID   string `json:"channelID"`
	SourceType  uint8  `json:"sourceType"`
	VccID       string `json:"vccID"`
	UserID      string `json:"userID"`
	GroupID     string `json:"groupID"`
	SessionID   string `json:"sessionID"`
	InQueueTime int64  `json:"inQueueTime"`
}

//连接成功消息体，发给MQ的
type ConnSuccessMQ struct {
	EventType       string `json:"eventType"`
	ChannelID       string `json:"channelID"`
	SourceType      uint8  `json:"sourceType"`
	VccID           string `json:"vccID"`
	UserID          string `json:"userID"`
	AgentID         string `json:"agentID"`
	GroupID         string `json:"groupID"`
	SessionID       string `json:"sessionID"`
	ConnSuccessTime int64  `json:"connSuccessTime"`
}

//出排队事件发送的消息体，发给MQ的
type OutQueueMQ struct {
	EventType    string `json:"eventType"`
	ChannelID    string `json:"channelID"`
	SourceType   uint8  `json:"sourceType"`
	VccID        string `json:"vccID"`
	UserID       string `json:"userID"`
	AgentID      string `json:"agentID"`
	GroupID      string `json:"groupID"`
	SessionID    string `json:"sessionID"`
	OutQueueTime int64  `json:"outQueueTime"`
}

//坐席状态变更
type AgentStatusMQ struct {
	EventType string `json:"eventType"`
	StampTime int64  `json:"timeStamp"`
	AgentId   string `json:"agentID"`
	VccId     string `json:"vccID"`
	Status    int    `json:"status"`   //0:离线；1：忙碌；2：在线
}

//会话结束
type SessionEndMq struct {
	SessionId string `json:"sessionId"`
	Cid string `json:"c_id"`
	ConSucTime string `json:"connSuccessTime"`
	VccId string `json:"vccId"`
	Name string `json:"userName"`
	AgentId string `json:"agentId"`
	WorkerId string `json:"agentWorkId"`
	UserSpeakNum string `json:"userSpeakNum"`
	AgentSpeakTime string `json:"agentSpeakNum"`
	SessionEndTime string `json:"endTime"`
	FirstRespTime string `json:"firstResponseSecs"`
	CreateType string `json:"createType"`
	EndType string `json:"endType"`
	SourceType string `json:"sourceType"`
	EvaluateStatus string `json:"evaluateStatus"`
	Next string `json:"next"`
	SessionFrom string `json:"sessionFrom"`
	UserId string `json:"userId"`
	GiveUpQueueing string `json:"giveupQueueing"`
	Source string `json:"source"`
}

type Source struct {
	SourceId string `json:"sourceId"`
	SourceName string `json:"sourceName"`
}

//会话结束的会话内容消息体，发送MQ
type SessionContentMQ struct {
	EventType  string   `json:"eventType"`
	ChannelID  string   `json:"channelId"`
	SourceType int      `json:"sourceType"`
	SessionID  string   `json:"sessionId"`
	VccID      string   `json:"vccId"`
	UserID     string   `json:"userId"`
	Index      int      `json:"index"`
	Content    []map[string]interface{} `json:"content"`
}

//网站渠道变更
type WebChannels struct {
	Id           string      `json:"id"`
	VccId        string      `json:"vccId"`
	Action       int         `json:"action"` //添加或者删除
	SourceName   string      `json:"sourceName"`
	ShortName    string      `json:"shortName"`
	ChannelStyle interface{} `json:"channelStyle"`
}

//微信渠道变更
type WebChatChannels struct {
	Id              string `json:"id"`
	UserName        string `json:"userName"`
	VccId           string `json:"vccId"`
	NickName        string `json:"nickName"`
	HeadImg         string `json:"headImg"`
	ServiceTypeInfo string `json:"serviceTypeInfo"`
	VerifyTypeInfo  string `json:"verifyTypeInfo"`
	AppId           string `json:"appId"`
	FuncInfo        string `json:"funcInfo"`
	PrincipalName   string `json:"principalName"`
}

//坐席规则配置变更
type AgentRule struct {
	VccID         string `json:"vccId"`   //企业ID
	AgentID       string `json:"agentId"` //坐席ID
	Belong        string `json:"belong"`  //
	TotalCapacity int    `json:"totalCapacity"`
	AgentAvatar   string `json:"agentAvatar"`
	NickName      string `json:"nickName"`
	AgentWorkId   string `json:"agentWorkId"` //工号
	Name          string `json:"name"`
	ImId          string `json:"imId"`
	DepId 		  string `json:"dep_id"`
}