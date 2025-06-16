package cache

import (
	"sync"

	v1 "zflow/api/base"
)

// Cache 缓存
type Cache struct {
	mu        sync.RWMutex
	nodeTypes map[string]map[string]*v1.NodeType       // service -> nodeTypeID -> NodeType
	connTypes map[string]map[string]*v1.ConnectionType // service -> connTypeID -> ConnectionType
}

// NewCache 创建缓存
func NewCache() *Cache {
	return &Cache{
		nodeTypes: make(map[string]map[string]*v1.NodeType),
		connTypes: make(map[string]map[string]*v1.ConnectionType),
	}
}

// AddNodeType 添加节点类型
func (c *Cache) AddNodeType(service string, nodeType *v1.NodeType) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.nodeTypes[service]; !ok {
		c.nodeTypes[service] = make(map[string]*v1.NodeType)
	}
	c.nodeTypes[service][nodeType.Uid] = nodeType
}

// AddConnType 添加连接类型
func (c *Cache) AddConnType(service string, connType *v1.ConnectionType) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.connTypes[service]; !ok {
		c.connTypes[service] = make(map[string]*v1.ConnectionType)
	}
	c.connTypes[service][connType.Uid] = connType
}

// GetNodeTypes 获取所有节点类型
func (c *Cache) GetNodeTypes() map[string]map[string]*v1.NodeType {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]map[string]*v1.NodeType)
	for service, types := range c.nodeTypes {
		result[service] = make(map[string]*v1.NodeType)
		for id, nodeType := range types {
			result[service][id] = nodeType
		}
	}
	return result
}

// GetConnTypes 获取所有连接类型
func (c *Cache) GetConnTypes() map[string]map[string]*v1.ConnectionType {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]map[string]*v1.ConnectionType)
	for service, types := range c.connTypes {
		result[service] = make(map[string]*v1.ConnectionType)
		for id, connType := range types {
			result[service][id] = connType
		}
	}
	return result
}
