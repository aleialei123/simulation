package fabric_prepare

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"zhanghefan123/security_topology/configs"
)

const (
	InitializePathMap             = "InitializePathMap"
	GenerateGenesisBlockConfigYml = "GenerateGenesisBlockConfigYml"
	GenerateOrdererOrgCryptoYml   = "GenerateOrdererOrgCryptoYml"
	GeneratePeerOrgCryptoYml      = "GeneratePeerOrgCryptoYml"
	InvokeCryptogenTool           = "InvokeCryptogenTool"
)

const (
	FabricBin           = "FabricBin"
	FabricBinCryptogen  = "FabricBinCryptogen"
	GenesisBlockBasic   = "GenesisBlockBasic"
	GenesisBlockNew     = "GenesisBlockNew"
	OrdererOrgCryptoNew = "OrdererOrgCryptoNew"
	PeerOrgCryptoNew    = "PeerOrgCryptoNew"
	organizationsPath   = "organizationsPath"
)

type GenerateFunction func() error

func (p *FabricPrepare) Generate() error {
	generateSteps := []map[string]GenerateFunction{
		{InitializePathMap: p.InitializePathMap},
		{GenerateGenesisBlockConfigYml: p.GenerateGenesisBlockConfigYml},
		{GenerateOrdererOrgCryptoYml: p.GenerateOrdererOrgCryptoYml},
		{GeneratePeerOrgCryptoYml: p.GeneratePeerOrgCryptoYml},
		{InvokeCryptogenTool: p.InvokeCryptogenTool},
	}
	err := p.generatePrepareSteps(generateSteps)
	if err != nil {
		return fmt.Errorf("generate prepare failed %w", err)
	}
	return nil
}

func (p *FabricPrepare) InitializePathMap() error {
	if _, ok := p.generateSteps[InitializePathMap]; ok {
		prepareWorkLogger.Infof("already initialize path mapping")
		return nil
	}
	fabricConfig := configs.TopConfiguration.FabricConfig
	fabricProjectPath := fabricConfig.FabricProjectPath
	fabricNetworkPath := fabricConfig.FabricNetworkPath

	p.pathMapping[FabricBin] = path.Join(fabricProjectPath, "bin")
	p.pathMapping[FabricBinCryptogen] = path.Join(fabricProjectPath, "bin/cryptogen")
	p.pathMapping[GenesisBlockBasic] = path.Join(fabricNetworkPath, "bft_config/configtx_basic.yaml")
	p.pathMapping[GenesisBlockNew] = path.Join(fabricNetworkPath, "bft_config/configtx.yaml")
	p.pathMapping[OrdererOrgCryptoNew] = path.Join(fabricNetworkPath, "organizations/cryptogen/crypto-config-orderer.yaml")
	p.pathMapping[PeerOrgCryptoNew] = path.Join(fabricNetworkPath, "organizations/cryptogen")
	p.pathMapping[organizationsPath] = path.Join(fabricNetworkPath, "organizations")
	prepareWorkLogger.Infof("successfully initialize path mapping")
	p.generateSteps[InitializePathMap] = struct{}{}
	return nil
}

// generateSteps 按步骤进行初始化
func (p *FabricPrepare) generatePrepareSteps(generateSteps []map[string]GenerateFunction) (err error) {
	fmt.Println("generateStepsHello")
	moduleNum := len(generateSteps)
	for idx, initStep := range generateSteps {
		for name, generateFunc := range initStep {
			if err = generateFunc(); err != nil {
				return fmt.Errorf("generate step [%s] failed, %w", name, err)
			}
			prepareWorkLogger.Infof("Generate STEP (%d/%d) => init step [%s] success)", idx+1, moduleNum, name)
		}
	}
	fmt.Println()
	return
}

func (p *FabricPrepare) GenerateGenesisBlockConfigYml() error {
	if _, ok := p.generateSteps[GenerateGenesisBlockConfigYml]; ok {
		prepareWorkLogger.Errorf("already generate genesis block config file")
		return nil
	}
	inputFile := p.pathMapping[GenesisBlockBasic]
	outputFile := p.pathMapping[GenesisBlockNew]
	ordererNum := p.fabricOrdererNodeCount
	peerNum := p.fabricPeerNodeCount
	orderStartPort := configs.TopConfiguration.FabricConfig.OrderStartPort
	// 打开输入文件
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 创建输出文件
	newFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Errorf("failed to create file: %v", err)
	}
	defer newFile.Close()
	err = os.Chmod(outputFile, 0777)
	if err != nil {
		fmt.Errorf("failed to open Permission: %v", err)
	}

	// 使用 bufio.Scanner 逐行读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// 检查是否是 "OrdererEndpoints:" 行
		if strings.Contains(line, "need_to_fill_1") {
			for i := 1; i <= ordererNum; i++ {

				fmt.Fprintf(newFile, "      - orderer%d.example.com:%d\n", i, orderStartPort+i)
				// fmt.Fprintf(newFile, "      - %s:%d\n", p.ipv4AddressesOrderer[i-1], orderStartPort+i)
			}
		} else if strings.Contains(line, "need_to_fill_2") {
			for i := 1; i <= peerNum; i++ {
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
			for i := 1; i <= ordererNum; i++ {
				fmt.Fprintf(newFile, "        - ID: %d\n", i)
				fmt.Fprintf(newFile, "          Host: orderer%d.example.com\n", i)
				// fmt.Fprintf(newFile, "          Host: %s\n", p.ipv4AddressesOrderer[i-1])
				fmt.Fprintf(newFile, "          Port: %d\n", orderStartPort+i)
				fmt.Fprintf(newFile, "          MSPID: OrdererMSP\n")
				fmt.Fprintf(newFile, "          Identity: ../organizations/ordererOrganizations/example.com/orderers/orderer%d.example.com/msp/signcerts/orderer%d.example.com-cert.pem\n", i, i)
				fmt.Fprintf(newFile, "          ClientTLSCert: ../organizations/ordererOrganizations/example.com/orderers/orderer%d.example.com/tls/server.crt\n", i)
				fmt.Fprintf(newFile, "          ServerTLSCert: ../organizations/ordererOrganizations/example.com/orderers/orderer%d.example.com/tls/server.crt\n", i)
			}
		} else if strings.Contains(line, "need_to_fill_4") {
			for i := 1; i <= peerNum; i++ {
				fmt.Fprintf(newFile, "        - *Org%d\n", i)
			}
		} else {
			_, err := fmt.Fprintln(newFile, line)
			if err != nil {
				fmt.Errorf("failed to write line to new file: %v", err)
			}
		}
	}

	// 检查扫描时是否出错
	if err := scanner.Err(); err != nil {
		fmt.Errorf("failed to read file: %v", err)
	}
	newFile.Close()

	// 提示成功
	prepareWorkLogger.Infof("Successfully updated the YAML and saved to %s\n", outputFile)
	return nil
}

func (p *FabricPrepare) GenerateOrdererOrgCryptoYml() error {
	if _, ok := p.generateSteps[GenerateOrdererOrgCryptoYml]; ok {
		prepareWorkLogger.Errorf("already generate orderer organization cryptogen file")
		return nil
	}
	outputFile := p.pathMapping[OrdererOrgCryptoNew]
	ordererNum := p.fabricOrdererNodeCount
	if ordererNum != 0 {
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()
		err = os.Chmod(outputFile, 0777)
		if err != nil {
			fmt.Errorf("failed to open Permission: %v", err)
		}
		fmt.Fprintf(file, "OrdererOrgs:\n")
		fmt.Fprintf(file, "  - Name: Orderer\n")
		fmt.Fprintf(file, "    Domain: example.com\n")
		fmt.Fprintf(file, "    EnableNodeOUs: true\n")
		fmt.Fprintf(file, "    Specs:\n")
		for i := 1; i <= ordererNum; i++ {
			fmt.Fprintf(file, "      - Hostname: orderer%d\n", i)
			fmt.Fprintf(file, "        SANS:\n")
			fmt.Fprintf(file, "          - localhost\n")
			// fmt.Fprintf(file, "          - %s\n", p.ipv4AddressesOrderer[i-1])
		}
		file.Close()
		prepareWorkLogger.Infof("Successfully updated the YAML and saved to %s\n", outputFile)
	}
	return nil
}
func (p *FabricPrepare) GeneratePeerOrgCryptoYml() error {
	if _, ok := p.generateSteps[GeneratePeerOrgCryptoYml]; ok {
		prepareWorkLogger.Errorf("already generate peer organization cryptogen file")
		return nil
	}
	peerNum := p.fabricPeerNodeCount
	for i := 1; i <= peerNum; i++ {
		outputFile := path.Join(p.pathMapping[PeerOrgCryptoNew], fmt.Sprintf("crypto-config-org%d.yaml", i))
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Errorf("failed to open file: %v", err)
		}
		err = os.Chmod(outputFile, 0777)
		if err != nil {
			fmt.Errorf("failed to open Permission: %v", err)
		}
		fmt.Fprintf(file, "PeerOrgs:\n")
		fmt.Fprintf(file, "  - Name: Org%d\n", i)
		fmt.Fprintf(file, "    Domain: org%d.example.com\n", i)
		fmt.Fprintf(file, "    EnableNodeOUs: true\n")
		fmt.Fprintf(file, "    Template:\n")
		fmt.Fprintf(file, "      Count: 1\n")
		fmt.Fprintf(file, "      SANS:\n")
		fmt.Fprintf(file, "        - localhost\n")
		// fmt.Fprintf(file, "        - %s\n", p.ipv4AddressesPeer[i-1])
		fmt.Fprintf(file, "    Users:\n")
		fmt.Fprintf(file, "      Count: 1\n")
		file.Close()
		prepareWorkLogger.Infof("Successfully updated the YAML and saved to %s\n", outputFile)
	}
	return nil
}
func (p *FabricPrepare) InvokeCryptogenTool() error {
	if _, ok := p.generateSteps[InvokeCryptogenTool]; ok {
		prepareWorkLogger.Errorf("already invoke crytogen tool")
		return nil
	}
	ordererNum := p.fabricOrdererNodeCount
	peerNum := p.fabricPeerNodeCount
	configFile := p.pathMapping[OrdererOrgCryptoNew]
	if ordererNum != 0 {
		cmd := exec.Command(p.pathMapping[FabricBinCryptogen], "generate", "--config", configFile, "--output", p.pathMapping[organizationsPath])
		_, err := cmd.CombinedOutput()
		if err != nil {
			prepareWorkLogger.Errorf("can not generate orderer msp")
		}
	}
	prepareWorkLogger.Infof("Successfully generate orderer msp")
	for i := 1; i <= peerNum; i++ {
		configFile := path.Join(p.pathMapping[PeerOrgCryptoNew], fmt.Sprintf("crypto-config-org%d.yaml", i))
		cmd := exec.Command(p.pathMapping[FabricBinCryptogen], "generate", "--config", configFile, "--output", p.pathMapping[organizationsPath])
		_, err := cmd.CombinedOutput()
		if err != nil {
			prepareWorkLogger.Errorf("can not generate peer msp")
		}
	}
	fmt.Println("p.pathMapping[organizationsPath]:%s", p.pathMapping[organizationsPath])
	err := setPermissions(p.pathMapping[organizationsPath], 0777)
	if err != nil {
		prepareWorkLogger.Errorf("Error changing permissions:", err)
	}
	return nil
}
func setPermissions(path string, mode os.FileMode) error {
	err := os.Chmod(path, mode)
	if err != nil {
		return err
	}

	// 如果是目录，递归修改目录下的所有文件和子目录
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if err := os.Chmod(path, mode); err != nil {
				return err
			}
			return nil
		})
		return err
	}
	return nil
}
