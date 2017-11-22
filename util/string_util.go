package util

import (
	"strings"
	"strconv"
	"fmt"
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

func GetStringDefault(str interface{}, defaultStr string) string {
	if str == nil {
		return defaultStr
	}else if tmp, ok := str.(string); ok {
		if IsEmpty(tmp) {
			return defaultStr
		}
		return tmp
	}else {
		return defaultStr
	}
}

func GetInt(str interface{}) int {
	st := GetString(str)
	i, _ := strconv.Atoi(st)
	return i
}

func GetIntString(intVal interface{}) string{
	if tmp, ok := intVal.(int); ok {
		return fmt.Sprintf("%d", tmp)
	}else if tmp, ok := intVal.(int64); ok {
		return fmt.Sprintf("%d", tmp)
	}else if tmp, ok := intVal.(int8); ok {
		return fmt.Sprintf("%d", tmp)
	}
	return fmt.Sprintf("%v", intVal)
}

func GetDefaultInt(intVal, defaultVal int) int {
	if intVal <= 0 {
		return defaultVal
	}
	return intVal
}

func GetDefaultInt64(intVal, defaultVal int64) int64 {
	if intVal <= 0 {
		return defaultVal
	}
	return intVal
}