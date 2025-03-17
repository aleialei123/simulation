package cryptogen

import (
	"fmt"
	"log"
	"os"
)

func Crytogen_qxp(order_num int, org_num int) {
	outputFile := "./organizations/cryptogen/crypto-config-orderer.yaml"
	// 打开输入文件
	// file, err := os.Open(outputFile)
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()
	fmt.Fprintf(file, "OrdererOrgs:\n")
	fmt.Fprintf(file, "  - Name: Orderer\n")
	fmt.Fprintf(file, "    Domain: example.com\n")
	fmt.Fprintf(file, "    EnableNodeOUs: true\n")
	fmt.Fprintf(file, "    Specs:\n")
	for i := 1; i <= order_num; i++ {
		fmt.Fprintf(file, "      - Hostname: orderer%d\n", i)
		fmt.Fprintf(file, "        SANS:\n")
		fmt.Fprintf(file, "          - localhost\n")
	}
	file.Close()
	fmt.Printf("Successfully updated the YAML and saved to %s\n", outputFile)
	for i := 1; i <= org_num; i++ {
		outputFile := fmt.Sprintf("./organizations/cryptogen/crypto-config-org%d.yaml", i)
		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatalf("failed to open file: %v", err)
		}
		fmt.Fprintf(file, "PeerOrgs:\n")
		fmt.Fprintf(file, "  - Name: Org%d\n", i)
		fmt.Fprintf(file, "    Domain: org%d.example.com\n", i)
		fmt.Fprintf(file, "    EnableNodeOUs: true\n")
		fmt.Fprintf(file, "    Template:\n")
		fmt.Fprintf(file, "      Count: 1\n")
		fmt.Fprintf(file, "      SANS:\n")
		fmt.Fprintf(file, "        - localhost\n")
		fmt.Fprintf(file, "    Users:\n")
		fmt.Fprintf(file, "      Count: 1\n")
		file.Close()
		fmt.Printf("Successfully updated the YAML and saved to %s\n", outputFile)
	}
}
