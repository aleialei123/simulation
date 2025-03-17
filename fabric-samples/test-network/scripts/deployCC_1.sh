read -p "orderer节点数量:" n


# ./network.sh createChannel \
#             -c mychannel \
#             -bft \
#             -orderer_num $n \
#             -peer_num 1
# ./network.sh deployCC \
#             -ccn secured \
#             -ccp ../asset-transfer-secured-agreement/chaincode-go/ \
#             -ccl go \
#             -ccep "OR('Org1MSP.peer')" \
#             -bft \
#             -orderer_num $n \
#             -peer_num 1

    echo $n