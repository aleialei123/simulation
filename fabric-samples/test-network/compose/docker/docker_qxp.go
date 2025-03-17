package docker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func Docker_qxp(peer_num int) {
	// 输入文件路径
	inputFile := "./compose/docker/docker-compose-bft-test-net-basic.yaml"
	// 输出文件路径
	outputFile := "./compose/docker/docker-compose-bft-test-net.yaml"
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

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "need_to_fill_1") {
			for i := 1; i <= peer_num; i++ {
				fmt.Fprintf(newFile, "  peer0.org%d.example.com:\n", i)
				fmt.Fprintf(newFile, "    container_name: peer0.org%d.example.com\n", i)
				fmt.Fprintf(newFile, "    image: hyperledger/fabric-peer:latest\n")
				fmt.Fprintf(newFile, "    labels:\n")
				fmt.Fprintf(newFile, "      service: hyperledger-fabric\n")
				fmt.Fprintf(newFile, "    environment:\n")
				fmt.Fprintf(newFile, "      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock\n")
				fmt.Fprintf(newFile, "      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=fabric_test\n")
				fmt.Fprintf(newFile, "    volumes:\n")
				fmt.Fprintf(newFile, "      - ./docker/peercfg:/etc/hyperledger/peercfg\n")
				fmt.Fprintf(newFile, "      - ${DOCKER_SOCK}:/host/var/run/docker.sock\n")
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
