package initializations

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	node, _ = snowflake.NewNode(1)
}

func IdGenerate() int64 {
	return node.Generate().Int64()
}
