package tool

import (
	v1 "zflow/api/base"
	"zflow/app/zflow/model"
)

// ConvertNodeType 将 model.NodeType 转换为 v1.NodeType
func ConvertNodeType(nt model.NodeType) *v1.NodeType {
	properties := make(map[string]*v1.PortList)
	for k, ports := range nt.Properties {
		portList := &v1.PortList{
			Ports: make([]*v1.Port, len(ports)),
		}
		for i, port := range ports {
			portList.Ports[i] = &v1.Port{
				Name:     port.Name,
				Label:    port.Label,
				PortType: port.PortType,
			}
		}
		properties[k] = portList
	}

	return &v1.NodeType{
		Uid:        nt.UID,
		Category:   nt.Category,
		Note:       nt.Note,
		Properties: properties,
	}
}

// ConvertConnType 将 model.ConnectionType 转换为 v1.ConnectionType
func ConvertConnType(ct model.ConnectionType) *v1.ConnectionType {
	return &v1.ConnectionType{
		Uid:              ct.UID,
		Name:             ct.Name,
		Description:      ct.Description,
		Color:            ct.Color,
		AllowedPortTypes: ct.AllowedPortTypes,
	}
}
