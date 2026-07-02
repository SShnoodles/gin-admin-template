package config

import (
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func InitIDGenerator() error {
	var err error
	node, err = snowflake.NewNode(1)
	return err
}

func IdGenerate() int64 {
	return node.Generate().Int64()
}
