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