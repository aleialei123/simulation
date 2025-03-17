package bft_config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"test-network/config_const"
)

// const order_port_start = 7050

func Generate_configtx(order_num, org_num int) {
	// var peer_num, order_num, org_num int
	// fmt.Println("Please enter the number of peer nodes")
	// fmt.Scanf("%d", &peer_num)
	// fmt.Println("Please enter the number of order nodes")
	// fmt.Scanf("%d", &order_num)
	// fmt.Println("Please enter the number of organization nodes")
	// fmt.Scanf("%d", &org_num)
	// 输入文件路径
	inputFile := "./bft_config/configtx_basic.yaml"
	// 输出文件路径
	outputFile := "./bft_config/configtx.yaml"
	order_port_start := config_const.Order_port_start

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

	// 标志位，指示是否找到 "OrdererEndpoints:" 行
	// foundOrdererEndpoints := false

	for scanner.Scan() {
		line := scanner.Text()

		// 检查是否是 "OrdererEndpoints:" 行
		if strings.Contains(line, "need_to_fill_1") {
			for i := 1; i <= order_num; i++ {
				fmt.Fprintf(newFile, "      - orderer%d.example.com:%d\n", i, order_port_start+i)
			}
		} else if strings.Contains(line, "need_to_fill_2") {
			for i := 1; i <= org_num; i++ {
				fmt.Fprintf(newFile, "  - &Org%d\n", i)
				fmt.Fprintf(newFile, "    Name: Org%dMSP\n", i)
				fmt.Fprintf(newFile, "    ID: Org%dMSP\n", i)
				fmt.Fprintf(newFile, "    MSPDir: ../organizations/peerOrganizations/org%d.example.com/msp\n", i)
				fmt.Fprintf(newFile, "    Policies:\n")
				fmt.Fprintf(newFile, "      Readers:\n")
				fmt.Fprintf(newFile, "        Type: Signature\n")
				fmt.Fprintf(newFile, "        Rule: \"OR('Org%dMSP.admin', 'Org%dMSP.peer', 'Org%dMSP.client')\"\n", i, i, i)
				fmt.Fprintf(newFile, "      Writers:\n")
				fmt.Fprintf(newFile, "        Type: Signature\n")
				fmt.Fprintf(newFile, "        Rule: \"OR('Org%dMSP.admin', 'Org%dMSP.client')\"\n", i, i)
				fmt.Fprintf(newFile, "      Admins:\n")
				fmt.Fprintf(newFile, "        Type: Signature\n")
				fmt.Fprintf(newFile, "        Rule: \"OR('Org%dMSP.admin')\"\n", i)
				fmt.Fprintf(newFile, "      Endorsement:\n")
				fmt.Fprintf(newFile, "        Type: Signature\n")
				fmt.Fprintf(newFile, "        Rule: \"OR('Org%dMSP.peer')\"\n", i)
			}
		} else if strings.Contains(line, "need_to_fill_3") {
			for i := 1; i <= order_num; i++ {
				fmt.Fprintf(newFile, "        - ID: %d\n", i)
				fmt.Fprintf(newFile, "          Host: orderer%d.example.com\n", i)
				fmt.Fprintf(newFile, "          Port: %d\n", order_port_start+i)
				fmt.Fprintf(newFile, "          MSPID: OrdererMSP\n")
				fmt.Fprintf(newFile, "          Identity: ../organizations/ordererOrganizations/example.com/orderers/orderer%d.example.com/msp/signcerts/orderer%d.example.com-cert.pem\n", i, i)
				fmt.Fprintf(newFile, "          ClientTLSCert: ../organizations/ordererOrganizations/example.com/orderers/orderer%d.example.com/tls/server.crt\n", i)
				fmt.Fprintf(newFile, "          ServerTLSCert: ../organizations/ordererOrganizations/example.com/orderers/orderer%d.example.com/tls/server.crt\n", i)
			}
		} else if strings.Contains(line, "need_to_fill_4") {
			for i := 1; i <= org_num; i++ {
				fmt.Fprintf(newFile, "        - *Org%d\n", i)
			}
		} else {
			_, err := fmt.Fprintln(newFile, line)
			if err != nil {
				log.Fatalf("failed to write line to new file: %v", err)
			}
		}
	}

	// 检查扫描时是否出错
	if err := scanner.Err(); err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	// 提示成功
	fmt.Printf("Successfully updated the YAML and saved to %s\n", outputFile)
}
