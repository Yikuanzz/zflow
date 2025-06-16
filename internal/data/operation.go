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

// ... existing code ...
// AddOperation 实现
// 加法操作，输入 a、b，输出 sum
type AddOperation struct{}

func (op *AddOperation) Execute(ctx model.Context, inputs map[string][]byte, vars map[string]interface{}) (map[string][]byte, error) {
	aBytes, aOk := inputs["a"]
	bBytes, bOk := inputs["b"]
	if !aOk || !bOk {
		return nil, fmt.Errorf("加法节点缺少输入 a 或 b")
	}
	var a, b int
	_, err := fmt.Sscanf(string(aBytes), "%d", &a)
	if err != nil {
		return nil, fmt.Errorf("加法节点输入 a 解析失败: %v", err)
	}
	_, err = fmt.Sscanf(string(bBytes), "%d", &b)
	if err != nil {
		return nil, fmt.Errorf("加法节点输入 b 解析失败: %v", err)
	}
	sum := a + b
	return map[string][]byte{"sum": []byte(fmt.Sprintf("%d", sum))}, nil
}

var AddOperationInst = &AddOperation{}

// MulOperation 实现
// 乘法操作，输入 a、b，输出 product
type MulOperation struct{}

func (op *MulOperation) Execute(ctx model.Context, inputs map[string][]byte, vars map[string]interface{}) (map[string][]byte, error) {
	aBytes, aOk := inputs["a"]
	bBytes, bOk := inputs["b"]
	if !aOk || !bOk {
		return nil, fmt.Errorf("乘法节点缺少输入 a 或 b")
	}
	var a, b int
	_, err := fmt.Sscanf(string(aBytes), "%d", &a)
	if err != nil {
		return nil, fmt.Errorf("乘法节点输入 a 解析失败: %v", err)
	}
	_, err = fmt.Sscanf(string(bBytes), "%d", &b)
	if err != nil {
		return nil, fmt.Errorf("乘法节点输入 b 解析失败: %v", err)
	}
	product := a * b
	return map[string][]byte{"product": []byte(fmt.Sprintf("%d", product))}, nil
}

var MulOperationInst = &MulOperation{}

// EchoOperation 实现（新版，支持 input->output）
type EchoOperationV2 struct{}

func (op *EchoOperationV2) Execute(ctx model.Context, inputs map[string][]byte, vars map[string]interface{}) (map[string][]byte, error) {
	input, ok := inputs["input"]
	if !ok {
		return nil, fmt.Errorf("Echo 节点缺少 input 输入")
	}
	ctx.Log(fmt.Sprintf("Echo: %s", string(input)))
	return map[string][]byte{"output": input}, nil
}

var EchoOperationInst = &EchoOperationV2{}

// ... existing code ...
