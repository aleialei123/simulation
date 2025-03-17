package fabric_prepare

import (
	"zhanghefan123/security_topology/modules/logger"
)

var (
	prepareWorkLogger = logger.GetLogger(logger.ModulePrepare)
)

type FabricPrepare struct {
	fabricOrdererNodeCount int
	fabricPeerNodeCount    int
	generateSteps          map[string]struct{}
	pathMapping            map[string]string
	ipv4AddressesOrderer   []string
	ipv4AddressesPeer      []string
}

func NewFabricPrepare(fabricOrdererNodeCount int, fabricPeerNodeCount int, ipv4AddressesOrderer []string, ipv4AddressesPeer []string) *FabricPrepare {
	return &FabricPrepare{
		fabricOrdererNodeCount: fabricOrdererNodeCount,
		fabricPeerNodeCount:    fabricPeerNodeCount,
		generateSteps:          make(map[string]struct{}),
		pathMapping:            make(map[string]string),
		ipv4AddressesOrderer:   ipv4AddressesOrderer,
		ipv4AddressesPeer:      ipv4AddressesPeer,
	}
}
