package topology

import (
	"fmt"
	"net"
	"zhanghefan123/security_topology/configs"
)

type ChainmakerPorts struct {
	p2pPort int
	rpcPort int
}

type FabricordererPorts struct {
	generalPort   int
	adminPort     int
	operationPort int
	p2pPort int
	rpcPort int
}

type FabricPeerPorts struct {
	listenPort    int
	operationPort int
	p2pPort int
	rpcPort int
}

// GetChainMakerNodeListenAddresses 获取长安链节点的监听地址
func (t *Topology) GetChainMakerNodeListenAddresses() []string {
	listeningAddresses := make([]string, 0)
	for _, node := range t.ChainmakerNodes {
		ip, _, _ := net.ParseCIDR(node.Interfaces[0].Ipv4Addr)
		listeningAddresses = append(listeningAddresses, ip.String())
	}
	return listeningAddresses
}

// GetFabricPeerNodeListenAddresses 获取Fabric Peer节点的监听地址
func (t *Topology) GetFabricPeerNodeListenAddresses() []string {
	listeningAddresses := make([]string, 0)
	for _, node := range t.FabricPeerNodes {
		ip, _, _ := net.ParseCIDR(node.Interfaces[0].Ipv4Addr)
		listeningAddresses = append(listeningAddresses, ip.String())
	}
	return listeningAddresses
}

// GetFabricOrdererNodeListenAddresses 获取Fabric Orderer节点的监听地址
func (t *Topology) GetFabricOrdererNodeListenAddresses() []string {
	listeningAddresses := make([]string, 0)
	for _, node := range t.FabricOrdererNodes {
		ip, _, _ := net.ParseCIDR(node.Interfaces[0].Ipv4Addr)
		listeningAddresses = append(listeningAddresses, ip.String())
	}
	return listeningAddresses
}

// GetChainMakerNodeContainerNames 获取所有长安链容器的名称
func (t *Topology) GetChainMakerNodeContainerNames() []string {
	chainMakerNodeNames := make([]string, 0)
	for _, node := range t.ChainmakerNodes {
		chainMakerNodeNames = append(chainMakerNodeNames, node.ContainerName)
	}
	return chainMakerNodeNames
}

// GetContainerNameToAddressMapping 获取所有节点的从容器名称到地址的一个映射
func (t *Topology) GetContainerNameToAddressMapping() (map[string]string, error) {
	addressMapping := make(map[string]string)
	for _, abstractNode := range t.AllAbstractNodes {
		normalNode, err := abstractNode.GetNormalNodeFromAbstractNode()
		if err != nil {
			return nil, fmt.Errorf("GetContainerNameToAddressMapping abstract node error: %w", err)
		}
		ip, _, _ := net.ParseCIDR(normalNode.Interfaces[0].Ipv4Addr)
		addressMapping[normalNode.ContainerName] = ip.String()
	}
	return addressMapping, nil
}

// GetContainerNameToGraphIdMapping 获取从所有节点的容器名称到图节点id的一个映射
func (t *Topology) GetContainerNameToGraphIdMapping() (map[string]int, error) {
	idMapping := make(map[string]int)
	for _, abstractNode := range t.AllAbstractNodes {
		normalNode, err := abstractNode.GetNormalNodeFromAbstractNode()
		if err != nil {
			return nil, fmt.Errorf("GetContainerNameToGraphIdMapping abstract node error: %w", err)
		}
		idMapping[normalNode.ContainerName] = int(abstractNode.Node.ID() + 1)
	}
	return idMapping, nil
}

// GetContainerNameToPortMapping 获取所有共识节点从容器名到
func (t *Topology) GetContainerNameToPortMapping() (map[string]*ChainmakerPorts, map[string]*FabricordererPorts, map[string]*FabricPeerPorts, error) {
	portMapping := make(map[string]*ChainmakerPorts)
	portFabricOrdererMapping := make(map[string]*FabricordererPorts)
	portFabricPeerMapping := make(map[string]*FabricPeerPorts)
	p2pStartPort := configs.TopConfiguration.ChainMakerConfig.P2pStartPort
	rpcStartPort := configs.TopConfiguration.ChainMakerConfig.RpcStartPort
	for _, chainMakerNode := range t.ChainmakerNodes {
		p2pPort := chainMakerNode.Id + p2pStartPort - 1
		rpcPort := chainMakerNode.Id + rpcStartPort - 1
		portMapping[chainMakerNode.ContainerName] = &ChainmakerPorts{p2pPort, rpcPort}
	}
	orderGeneralListenStartPort := configs.TopConfiguration.FabricConfig.OrderGeneralListenStartPort
	orderAdminListenStartPort := configs.TopConfiguration.FabricConfig.OrderAdminListenStartPort
	orderOperationListenStartPort := configs.TopConfiguration.FabricConfig.OrderOperationListenStartPort
	orderP2pStartPort := configs.TopConfiguration.FabricConfig.OrderP2pStartPort
	orderRpcStartPort := configs.TopConfiguration.FabricConfig.OrderRpcStartPort
	for _, fabricOrdererNode := range t.FabricOrdererNodes {
		generalPort := fabricOrdererNode.Id + orderGeneralListenStartPort
		adminPort := fabricOrdererNode.Id + orderAdminListenStartPort
		operationPort := fabricOrdererNode.Id + orderOperationListenStartPort
		p2pPort := fabricOrdererNode.Id + orderP2pStartPort
		rpcPort := fabricOrdererNode.Id + orderRpcStartPort
		portFabricOrdererMapping[fabricOrdererNode.ContainerName] = &FabricordererPorts{generalPort, adminPort, operationPort,p2pPort,rpcPort}

	}
	peerListenStartPort := configs.TopConfiguration.FabricConfig.PeerListenStartPort
	peerOperationStartPort := configs.TopConfiguration.FabricConfig.PeerOperationStartPort
	peerP2pStartPort := configs.TopConfiguration.FabricConfig.PeerP2pStartPort
	peerRpcStartPort := configs.TopConfiguration.FabricConfig.PeerRpcStartPort
	for _, fabricPeerNode := range t.FabricPeerNodes {
		listenPort := fabricPeerNode.Id + peerListenStartPort
		operationPort := fabricPeerNode.Id + peerOperationStartPort
		p2pPort := fabricPeerNode.Id + peerP2pStartPort
		rpcPort := fabricPeerNode.Id + peerRpcStartPort
		portFabricPeerMapping[fabricPeerNode.ContainerName] = &FabricPeerPorts{listenPort, operationPort,p2pPort,rpcPort}

	}
	return portMapping, portFabricOrdererMapping, portFabricPeerMapping, nil
}
