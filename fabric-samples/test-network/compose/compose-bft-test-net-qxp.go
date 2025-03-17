package compose

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"test-network/config_const"
)

func Compose_bft_test_net_qxp(peer_num, order_num, org_num int) {
	// 输入文件路径
	inputFile := "./compose/compose-bft-test-net-basic.yaml"
	// 输出文件路径
	outputFile := "./compose/compose-bft-test-net.yaml"
	// 打开输入文件
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// 创建输出文件
	newFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer newFile.Close()

	// 使用 bufio.Scanner 逐行读取文件内容
	scanner := bufio.NewScanner(file)
	order_general_listenport_start := config_const.Order_general_listenport_start
	order_admin_listenport_start := config_const.Order_admin_listenport_start
	order_operation_listenport_start := config_const.Order_operation_listenport_start
	peer_listenport_start := config_const.Peer_listenport_start
	peer_chaincodeport_start := config_const.Peer_chaincodeport_start
	peer_operationport_start := config_const.Peer_operationport_start

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "need_to_fill_1") {
			for i := 1; i <= order_num; i++ {
				fmt.Fprintf(newFile, "  orderer%d.example.com:\n", i)
			}
			for i := 1; i <= peer_num; i++ {
				fmt.Fprintf(newFile, "  peer0.org%d.example.com:\n", i)
			}
		} else if strings.Contains(line, "need_to_fill_2") {
			for i := 1; i <= order_num; i++ {
				fmt.Fprintf(newFile, "  orderer%d.example.com:\n", i)
				fmt.Fprintf(newFile, "    container_name: orderer%d.example.com\n", i)
				fmt.Fprintf(newFile, "    image: hyperledger/fabric-orderer:latest\n")
				fmt.Fprintf(newFile, "    labels:\n")
				fmt.Fprintf(newFile, "      service: hyperledger-fabric\n")
				fmt.Fprintf(newFile, "    environment:\n")
				fmt.Fprintf(newFile, "      - FABRIC_LOGGING_SPEC=INFO\n")
				fmt.Fprintf(newFile, "      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0\n")
				fmt.Fprintf(newFile, "      - ORDERER_GENERAL_LISTENPORT=%d\n", order_general_listenport_start+i)
				fmt.Fprintf(newFile, "      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP\n")
				fmt.Fprintf(newFile, "      - ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp\n")
				fmt.Fprintf(newFile, "      - ORDERER_GENERAL_TLS_ENABLED=true\n")
				fmt.Fprintf(newFile, "      - ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key\n")
				fmt.Fprintf(newFile, "      - ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt\n")
				fmt.Fprintf(newFile, "      - ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]\n")
				fmt.Fprintf(newFile, "      - ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/var/hyperledger/orderer/tls/server.crt\n")
				fmt.Fprintf(newFile, "      - ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/var/hyperledger/orderer/tls/server.key\n")
				fmt.Fprintf(newFile, "      - ORDERER_GENERAL_CLUSTER_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]\n")
				fmt.Fprintf(newFile, "      - ORDERER_GENERAL_BOOTSTRAPMETHOD=none\n")
				fmt.Fprintf(newFile, "      - ORDERER_CHANNELPARTICIPATION_ENABLED=true\n")
				fmt.Fprintf(newFile, "      - ORDERER_ADMIN_TLS_ENABLED=true\n")
				fmt.Fprintf(newFile, "      - ORDERER_ADMIN_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt\n")
				fmt.Fprintf(newFile, "      - ORDERER_ADMIN_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key\n")
				fmt.Fprintf(newFile, "      - ORDERER_ADMIN_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]\n")
				fmt.Fprintf(newFile, "      - ORDERER_ADMIN_TLS_CLIENTROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]\n")
				fmt.Fprintf(newFile, "      - ORDERER_ADMIN_LISTENADDRESS=0.0.0.0:%d\n", order_admin_listenport_start+i)
				fmt.Fprintf(newFile, "      - ORDERER_OPERATIONS_LISTENADDRESS=orderer%d.example.com:%d\n", i, order_operation_listenport_start+i)
				fmt.Fprintf(newFile, "      - ORDERER_METRICS_PROVIDER=prometheus\n")
				fmt.Fprintf(newFile, "    working_dir: /root\n")
				fmt.Fprintf(newFile, "    command: orderer\n")
				fmt.Fprintf(newFile, "    volumes:\n")
				fmt.Fprintf(newFile, "      - ../organizations/ordererOrganizations/example.com/orderers/orderer%d.example.com/msp:/var/hyperledger/orderer/msp\n", i)
				fmt.Fprintf(newFile, "      - ../organizations/ordererOrganizations/example.com/orderers/orderer%d.example.com/tls/:/var/hyperledger/orderer/tls\n", i)
				fmt.Fprintf(newFile, "      - orderer%d.example.com:/var/hyperledger/production/orderer\n", i)
				fmt.Fprintf(newFile, "    ports:\n")
				fmt.Fprintf(newFile, "      - %d:%d\n", order_general_listenport_start+i, order_general_listenport_start+i)
				fmt.Fprintf(newFile, "      - %d:%d\n", order_admin_listenport_start+i, order_admin_listenport_start+i)
				fmt.Fprintf(newFile, "      - %d:%d\n", order_operation_listenport_start+i, order_operation_listenport_start+i)
				fmt.Fprintf(newFile, "    networks:\n")
				fmt.Fprintf(newFile, "      - test\n")
			}
			for i := 1; i <= org_num; i++ {
				fmt.Fprintf(newFile, "  peer0.org%d.example.com:\n", i)
				fmt.Fprintf(newFile, "    container_name: peer0.org%d.example.com\n", i)
				fmt.Fprintf(newFile, "    image: hyperledger/fabric-peer:latest\n")
				fmt.Fprintf(newFile, "    labels:\n")
				fmt.Fprintf(newFile, "      service: hyperledger-fabric\n")
				fmt.Fprintf(newFile, "    environment:\n")
				fmt.Fprintf(newFile, "      - FABRIC_CFG_PATH=/etc/hyperledger/peercfg\n")
				fmt.Fprintf(newFile, "      - FABRIC_LOGGING_SPEC=INFO\n")
				fmt.Fprintf(newFile, "      - CORE_PEER_TLS_ENABLED=true\n")
				fmt.Fprintf(newFile, "      - CORE_PEER_PROFILE_ENABLED=false\n")
				fmt.Fprintf(newFile, "      - CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt\n")
				fmt.Fprintf(newFile, "      - CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key\n")
				fmt.Fprintf(newFile, "      - CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt\n")
				fmt.Fprintf(newFile, "      - CORE_PEER_ID=peer0.org%d.example.com\n", i)
				fmt.Fprintf(newFile, "      - CORE_PEER_ADDRESS=peer0.org%d.example.com:%d\n", i, peer_listenport_start+i)
				fmt.Fprintf(newFile, "      - CORE_PEER_LISTENADDRESS=0.0.0.0:%d\n", peer_listenport_start+i)
				fmt.Fprintf(newFile, "      - CORE_PEER_CHAINCODEADDRESS=peer0.org%d.example.com:%d\n", i, peer_chaincodeport_start+i)
				fmt.Fprintf(newFile, "      - CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:%d\n", peer_chaincodeport_start+i)
				fmt.Fprintf(newFile, "      - CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.org%d.example.com:%d\n", i, peer_listenport_start+i)
				fmt.Fprintf(newFile, "      - CORE_PEER_GOSSIP_BOOTSTRAP=peer0.org%d.example.com:%d\n", i, peer_listenport_start+i)
				fmt.Fprintf(newFile, "      - CORE_PEER_LOCALMSPID=Org%dMSP\n", i)
				fmt.Fprintf(newFile, "      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/msp\n")
				fmt.Fprintf(newFile, "      - CORE_OPERATIONS_LISTENADDRESS=peer0.org%d.example.com:%d\n", i, peer_operationport_start+i)
				fmt.Fprintf(newFile, "      - CORE_METRICS_PROVIDER=prometheus\n")
				fmt.Fprintf(newFile, "      - CHAINCODE_AS_A_SERVICE_BUILDER_CONFIG={\"peername\":\"peer0org%d\"}\n", i)
				fmt.Fprintf(newFile, "      - CORE_CHAINCODE_EXECUTETIMEOUT=300s\n")
				fmt.Fprintf(newFile, "    volumes:\n")
				fmt.Fprintf(newFile, "      - ../organizations/peerOrganizations/org%d.example.com/peers/peer0.org%d.example.com:/etc/hyperledger/fabric\n", i, i)
				fmt.Fprintf(newFile, "      - peer0.org%d.example.com:/var/hyperledger/production\n", i)
				fmt.Fprintf(newFile, "    working_dir: /root\n")
				fmt.Fprintf(newFile, "    command: peer node start\n")
				fmt.Fprintf(newFile, "    ports:\n")
				fmt.Fprintf(newFile, "      - %d:%d\n", peer_listenport_start+i, peer_listenport_start+i)
				fmt.Fprintf(newFile, "      - %d:%d\n", peer_operationport_start+i, peer_operationport_start+i)
				fmt.Fprintf(newFile, "    networks:\n")
				fmt.Fprintf(newFile, "      - test\n")
			}
		} else {
			_, err := fmt.Fprintln(newFile, line)
			if err != nil {
				log.Fatalf("failed to write line to new file: %v", err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	// 提示成功
	fmt.Printf("Successfully updated the YAML and saved to %s\n", outputFile)
}
