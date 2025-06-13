package data

import (
	"fmt"
	"zflow/internal/model"
)

// EchoOperation 是一个简单的回显操作，用于测试
type EchoOperation struct {
	Message string
}

func (op *EchoOperation) Execute(ctx model.Context, input []byte, vars map[string]interface{}) ([]byte, error) {
	ctx.Log(fmt.Sprintf("Echo: %s", op.Message))
	return []byte(op.Message), nil
}

// FileReadOperation 读取文件内容
type FileReadOperation struct {
	Encoding string // utf-8, gbk, etc.
}

func (op *FileReadOperation) Execute(ctx model.Context, input []byte, vars map[string]interface{}) ([]byte, error) {
	// TODO: 实现文件读取逻辑
	return nil, nil
}

// FileWriteOperation 写入文件内容
type FileWriteOperation struct {
	Encoding string // utf-8, gbk, etc.
}

func (op *FileWriteOperation) Execute(ctx model.Context, input []byte, vars map[string]interface{}) ([]byte, error) {
	// TODO: 实现文件写入逻辑
	return nil, nil
}

// TransformOperation 数据转换操作
type TransformOperation struct {
	Template string // 转换模板
}

func (op *TransformOperation) Execute(ctx model.Context, input []byte, vars map[string]interface{}) ([]byte, error) {
	// TODO: 实现数据转换逻辑
	return nil, nil
}
