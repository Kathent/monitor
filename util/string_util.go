package util

import (
	"strings"
	"strconv"
)

func IsEmpty(val string) bool{
	return len(strings.TrimSpace(val)) <= 0
}

func GetString(str interface{}) string {
	if str == nil {
		return ""
	}else if tmp, ok := str.(string); ok {
		return tmp
	}else {
		return ""
	}
}

func GetInt(str interface{}) int {
	st := GetString(str)
	i, _ := strconv.Atoi(st)
	return i
}