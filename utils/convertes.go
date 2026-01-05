package utils

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)


func SafeInt64(v interface{}) int64 {
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int:
		return int64(val)
	case int64:
		return val
	case nil:
		return 0
	default:
		log.Println("Cannot convert to int64:", v)
		return 0
	}
}

func SafeString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}


func NowUnixMilli() int64 {
	return int64(1e3) * int64((float64)(websocket.DefaultDialer.HandshakeTimeout.Seconds()))
}