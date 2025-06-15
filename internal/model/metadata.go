package model

// Port 描述"某种节点类型"暴露出的端口
type Port struct {
	Name     string `json:"name"`      // in / out1 / pdf_out ...
	Label    string `json:"label"`     // 可选，人类可读
	PortType string `json:"port_type"` // connection / file
}

// Endpoint 表示一条连接线上的"端点"
type Endpoint struct {
	NodeID   string `json:"node_id"`
	PortName string `json:"port_name"`
}

// NodeType 定义节点模板
type NodeType struct {
	UID        string            `json:"node_type"`
	Category   string            `json:"category"`
	Note       string            `json:"note"`
	Operation  Operation         `json:"operation"`
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
