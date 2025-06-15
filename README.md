# 工作流引擎运行流程

```text
             ┌──────────────┐
             │ 加载 JSON 配置 │
             └──────┬───────┘
                    │
        ┌───────────▼────────────┐
        │ 构建 Workflow（包含 Dag）│
        └──────┬────────────┬────┘
               │            │
   ┌───────────▼──┐     ┌───▼──────────────┐
   │ 校验节点与连接 │     │ 构建 NodeType / Op│
   └──────┬───────┘     └────────┬─────────┘
          │                      │
          └────┬─────────────────┘
               ▼
      ┌───────────────┐
      │ 拓扑排序（DAG） │
      └──────┬────────┘
             ▼
    ┌────────────────────┐
    │ 顺序执行 Node/Op    │
    │ - 输入收集          │
    │ - 执行 Operation   │
    │ - 输出转发          │
    └─────────┬──────────┘
              ▼
       ┌─────────────┐
       │ 完成 / 失败   │
       └─────────────┘


```



# 关键数据结构

用户可以选择不同的 **节点类型** 与 **连线类型** 来组成工作流，下面数据结构就是最关键的两个 `struct` ，根据不同 节点实例 和 连线实例，就可以构造出不同的工作流。

```go
// NodeType 定义节点模板
type NodeType struct {
    // UID 节点模板全局唯一标识
	UID        string            `json:"node_type"`
    // Category 节点模板类别
	Category   string            `json:"category"`
    // Note 节点模板的说明
	Note       string            `json:"note"`
    // Operation 节点模板要执行的操作
	Operation  Operation         `json:"operation"`
    // Properties 节点模板的 输入/输出
	Properties map[string][]Port `json:"properties"`
}

// ConnectionType 决定连线的语义与可连接端口类型
type ConnectionType struct {
	UID              string   `json:"connection_type"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Color            string   `json:"color"`
	AllowedPortTypes []string `json:"allowed_port_types"`
}

```

