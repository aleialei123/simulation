#!/bin/bash

# 设置环境变量
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:5051

#选择提交交易的orderer节点序号
num=3

# Orderer TLS证书路径
ORDERER_TLS_CERT=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer$num.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# 通道名称和链码名称
CHANNEL_NAME="mychannel"
CHAINCODE_NAME="secured"

read -p "请输入程序运行的时间:" STOP_TIME
SECONDS=0



# 创建资产数量
# read -p "请输入初始创建的资产数量: " n
n=1
# 使用 peer channel getinfo 提取链信息
CHAIN_INFO=$(peer channel getinfo -c $CHANNEL_NAME 2>/dev/null)

# 提取块高度信息
BLOCK_PRE_HEIGHT=$(echo $CHAIN_INFO | grep -o '"height":[0-9]*' | grep -o '[0-9]*')
echo "先前链高度：$BLOCK_PRE_HEIGHT"

# 无限循环提交交易
while [ $SECONDS -lt $STOP_TIME ]; do
  # 创建资产的 JSON 数据
  ASSET_PROPERTIES=$(echo -n "{\"object_type\":\"asset_properties\",\"color\":\"blue\",\"size\":$n,\"salt\":\"a94a8fe5ccb19ba61c4c0873d391e987982fbbd3\"}" | base64 | tr -d \\n)

  # 调用链码创建资产
  # echo "正在创建资产，size=$n..."
  peer chaincode invoke \
    -o localhost:705$num \
    --ordererTLSHostnameOverride orderer$num.example.com \
    --tls \
    --cafile "$ORDERER_TLS_CERT" \
    -C "$CHANNEL_NAME" \
    -n "$CHAINCODE_NAME" \
    -c '{"function":"CreateAsset","Args":["Asset'$n'"]}' \
    --transient "{\"asset_properties\":\"$ASSET_PROPERTIES\"}" \
    2>/dev/null

  # 检查返回状态
  if [ $? -ne 0 ]; then
    echo "创建资产失败，num=$num"
    num_pre=$num
    num=$(((num+1)%4))
    if [ $num -eq 0 ]; then
      num=4
    fi
    echo "num:$num_pre->$num"
    ORDERER_TLS_CERT=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer$num.example.com/msp/tlscacerts/tlsca.example.com-cert.pem 
    ((n--))
    # exit 1
  fi

  # 每次交易后递增资产的size                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          
  ((n++))


done
sleep 2

# 使用 peer channel getinfo 提取链信息
peer channel getinfo -c $CHANNEL_NAME > /dev/null 2>&1
CHAIN_INFO=$(peer channel getinfo -c $CHANNEL_NAME 2>/dev/null) > /dev/null 2>&1

# 提取块高度信息
BLOCK_CUR_HEIGHT=$(echo $CHAIN_INFO | grep -o '"height":[0-9]*' | grep -o '[0-9]*') > /dev/null 2>&1
BLOCK_ADD_HEIGHT=$(($BLOCK_CUR_HEIGHT-$BLOCK_PRE_HEIGHT))
echo "先前链高度：$BLOCK_PRE_HEIGHT;目前链高度：$BLOCK_CUR_HEIGHT;高度增长:$BLOCK_ADD_HEIGHT"