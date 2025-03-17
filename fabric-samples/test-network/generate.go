package main

import (
	"fmt"
	"test-network/bft_config"
	"test-network/compose"
	"test-network/compose/docker"
	"test-network/organizations/cryptogen"
)

func main() {
	var peer_num, order_num, org_num int
	fmt.Println("Please enter the number of peer nodes")
	fmt.Scanf("%d", &peer_num)
	fmt.Println("Please enter the number of order nodes")
	fmt.Scanf("%d", &order_num)
	fmt.Println("Please enter the number of organization nodes")
	fmt.Scanf("%d", &org_num)
	bft_config.Generate_configtx(order_num, org_num)
	compose.Compose_bft_test_net_qxp(peer_num, order_num, org_num)
	docker.Docker_qxp(peer_num)
	cryptogen.Crytogen_qxp(order_num, org_num)
	return
}
