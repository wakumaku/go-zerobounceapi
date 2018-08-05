package zerobounceapi

import (
	"strconv"
)

func getIntField(value interface{}) int {

	switch val := value.(type) {
	case int:
		return val
	case float64:
		return int(val)
	case string:
		if val, err := strconv.Atoi(value.(string)); err == nil {
			return val
		}
	}

	return 0
}

func getStringField(value interface{}) string {

	switch val := value.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(value.(int))

	}

	return ""
}

func getBoolField(value interface{}) bool {

	switch val := value.(type) {
	case bool:
		return val
	case string:
		if val, err := strconv.ParseBool(val); err == nil {
			return val
		}
	case int:
		if val != 0 {
			return true
		}
	}

	return false
}
