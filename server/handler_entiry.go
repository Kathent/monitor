package server

import (
	"fmt"

	"paas/icsoc-monitor/constants"
)

type MonitorRequest struct {
	VccId string `json:"vccId"`
	Start string `json:"start"`
	End string `json:"end"`
	AgentId string `json:"agentId"`
	Page int `json:"page"`
	PageSize int `json:"pageSize"`
	SortField string `json:"sortField"`
	SortOrder string `json:"sortOrder"`
	DepId string `json:"dep_id"`
	Status string `json:"status"`
}

type MonitorResponse struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Data struct{
		Total int `json:"total"`
		Rows []interface{} `json:"rows"`
	} `json:"data"`
}

type AgentRequest struct {
	VccId string
	Start string
	End string
	Page int
	PageSize int
	AgentId string
	StartIndex int
	EndIndex int
	SortField string
	SortOrder string
	DepId string
	Status string
}

func (ag AgentRequest) GetSortString() string{
	if ag.SortOrder == constants.SORT_ORDER_DEFAULT{
		return ag.SortField
	}
	return fmt.Sprintf("-%s", ag.SortField)
}

func (ag AgentRequest) GetSortSym() int {
	if ag.SortOrder == constants.SORT_ORDER_DEFAULT{
		return constants.SORT_ORDER_ASCEND_SYM
	}
	return constants.SORT_ORDER_DESCEND_SYM
}
