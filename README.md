

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



首先，前端调接口 **GET /node_types** 得到所有的 **节点类型**。

```json
{
  "builtin.add.v1": {
    "node_type": "builtin.add.v1",
    "category": "math",
    "note": "两个数字相加，输出结果",
    "operation": {

    },
    "properties": {
      "inputs": [
        {
          "name": "a",
          "label": "加数A",
          "port_type": "connection"
        },
        {
          "name": "b",
          "label": "加数B",
          "port_type": "connection"
        }
      ],
      "outputs": [
        {
          "name": "sum",
          "label": "和",
          "port_type": "connection"
        }
      ]
    }
  },
  "builtin.echo.v1": {
    "node_type": "builtin.echo.v1",
    "category": "util",
    "note": "回显输入内容，常用于调试或展示节点计算结果",
    "operation": {

    },
    "properties": {
      "inputs": [
        {
          "name": "input",
          "label": "输入内容",
          "port_type": "connection"
        }
      ],
      "outputs": [
        {
          "name": "output",
          "label": "输出内容",
          "port_type": "connection"
        }
      ]
    }
  },
  "builtin.mul.v1": {
    "node_type": "builtin.mul.v1",
    "category": "math",
    "note": "两个数字相乘，输出结果",
    "operation": {

    },
    "properties": {
      "inputs": [
        {
          "name": "a",
          "label": "乘数A",
          "port_type": "connection"
        },
        {
          "name": "b",
          "label": "乘数B",
          "port_type": "connection"
        }
      ],
      "outputs": [
        {
          "name": "product",
          "label": "积",
          "port_type": "connection"
        }
      ]
    }
  }
}
```



再通过 **GET /connection_types** 得到所有的 **连接类型**。

```json
{
  "1": {
    "connection_type": "1",
    "name": "data_flow",
    "description": "数据流连接，用于传递普通数据",
    "color": "#4CAF50",
    "allowed_port_types": [
      "connection"
    ]
  }
}
```



接着调用 **POST /run/work_flow** 传递搭建好的工作流





# 核心 API 接口







# 关键数据结构

### 设计工作流时可以用的数据结构

用户可以选择不同的 **节点类型** 与 **连线类型** 来组成工作流，下面数据结构就是最关键的两个 `struct` ，根据不同 节点实例 和 连线实例，就可以构造出不同的工作流。

> GET	/node_types	返回所有的节点类型
>
> GET	/connection_types	返回所有的连线类型

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



## 构建工作流时的数据实例

```JSON
{
  "nodes": [
    {
      "id": "add1",
      "node_type": "builtin.add.v1",
      "label": "加法节点"
    },
    {
      "id": "mul1",
      "node_type": "builtin.mul.v1",
      "label": "乘法节点"
    },
    {
      "id": "echo1",
      "node_type": "builtin.echo.v1",
      "label": "回显节点"
    }
  ],
  "connections": [
    {
      "connection_id": "c1",
      "connection_type": "1",
      "from": { "node_id": "add1", "port_name": "sum" },
      "to": { "node_id": "mul1", "port_name": "a" }
    },
    {
      "connection_id": "c2",
      "connection_type": "1",
      "from": { "node_id": "mul1", "port_name": "product" },
      "to": { "node_id": "echo1", "port_name": "input" }
    }
  ]
}
```





