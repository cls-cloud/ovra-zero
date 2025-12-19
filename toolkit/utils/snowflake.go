package utils

import (
	"strconv"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

// Init 初始化雪花算法节点
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(machineID)
	return
}

// GetID 生成ID
func GetID() string {
	return strconv.FormatInt(node.Generate().Int64(), 10)
}

func GetIDInt64() int64 {
	return node.Generate().Int64()
}
