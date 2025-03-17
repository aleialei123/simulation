#!/bin/bash

# 设置环境变量
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:5051

num=2

# Orderer TLS证书路径
ORDERER_TLS_CERT=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer$num.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# 通道名称和链码名称
CHANNEL_NAME="mychannel"
CHAINCODE_NAME="secured"

# 创建多个资产
read -p "请输入要创建的资产数量: " n

# 使用 peer channel getinfo 提取链信息
CHAIN_INFO=$(peer channel getinfo -c $CHANNEL_NAME 2>/dev/null)


# 提取块高度信息
BLOCK_PRE_HEIGHT=$(echo $CHAIN_INFO | grep -o '"height":[0-9]*' | grep -o '[0-9]*')
echo "先前链高度：$BLOCK_PRE_HEIGHT"

for ((size=1; size<=n; size++))
do
  # 创建资产的 JSON 数据
  ASSET_PROPERTIES=$(echo -n "{\"object_type\":\"asset_properties\",\"color\":\"blue\",\"size\":$size,\"salt\":\"a94a8fe5ccb19ba61c4c0873d391e987982fbbd3\"}" | base64 | tr -d \\n)

  # 调用链码创建资产
  echo "正在创建资产，size=$size..."
  peer chaincode invoke \
    -o localhost:705$num \
    --ordererTLSHostnameOverride orderer$num.example.com \
    --tls \
    --cafile "$ORDERER_TLS_CERT" \
    -C "$CHANNEL_NAME" \
    -n "$CHAINCODE_NAME" \
    -c '{"function":"CreateAsset","Args":["Asset'$size'"]}' \
    --transient "{\"asset_properties\":\"$ASSET_PROPERTIES\"}"

  # 检查返回状态
  if [ $? -ne 0 ]; then
    echo "创建资产失败，size=$size"
    exit 1
  fi

done
sleep 2
# 使用 peer channel getinfo 提取链信息
peer channel getinfo -c $CHANNEL_NAME
CHAIN_INFO=$(peer channel getinfo -c $CHANNEL_NAME 2>/dev/null)

# 提取块高度信息
BLOCK_CUR_HEIGHT=$(echo $CHAIN_INFO | grep -o '"height":[0-9]*' | grep -o '[0-9]*')

echo "先前链高度：$BLOCK_PRE_HEIGHT ;目前链高度：$BLOCK_CUR_HEIGHT"
