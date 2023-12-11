package snowflakenode

import (
	"github.com/bwmarrin/snowflake"
)

var Node *snowflake.Node

//
// NewNode
//  @Description: 创建唯一的雪花节点
//
func NewNode(nodeID int64) {
	Node, _ = snowflake.NewNode(nodeID)
}
