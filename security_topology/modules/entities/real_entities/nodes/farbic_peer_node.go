package nodes

import (
	"fmt"
	"zhanghefan123/security_topology/modules/entities/real_entities/normal_node"
	"zhanghefan123/security_topology/modules/entities/types"
)

type FabricPeerNode struct {
	*normal_node.NormalNode
}

func NewFabricPeerNode(nodeId int, X, Y float64) *FabricPeerNode {
	fabricNodeType := types.NetworkNodeType_FabricPeerNode
	normalNode := normal_node.NewNormalNode(
		types.NetworkNodeType_FabricPeerNode,
		nodeId,
		fmt.Sprintf("%s-%d", fabricNodeType.String(), nodeId))
	normalNode.X = X
	normalNode.Y = Y
	return &FabricPeerNode{normalNode}
}
