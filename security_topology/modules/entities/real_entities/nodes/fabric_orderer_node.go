package nodes

import (
	"fmt"
	"zhanghefan123/security_topology/modules/entities/real_entities/normal_node"
	"zhanghefan123/security_topology/modules/entities/types"
)

type FabricOrdererNode struct {
	*normal_node.NormalNode
}

func NewFabricOrdererNode(nodeId int, X, Y float64) *FabricOrdererNode {
	fabricNodeType := types.NetworkNodeType_FabricOrdererNode
	normalNode := normal_node.NewNormalNode(
		types.NetworkNodeType_FabricOrdererNode,
		nodeId,
		fmt.Sprintf("%s-%d", fabricNodeType.String(), nodeId))
	normalNode.X = X
	normalNode.Y = Y
	return &FabricOrdererNode{normalNode}
}
