#!/bin/bash
export PATH=${ROOTDIR}/../bin:${PWD}/../bin:$PATH



CHANNEL_NAME="$1"
DELAY="$2"
MAX_RETRY="$3"
VERBOSE="$4"
BFT="$5"
ORDERER_NUM="$6"
PEER_NUM="$7"
: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="3"}
: ${MAX_RETRY:="5"}
: ${VERBOSE:="false"}
: ${BFT:=0}

: ${CONTAINER_CLI:="docker"}

# imports  
. scripts/envVar.sh $ORDERER_NUM $PEER_NUM
if command -v ${CONTAINER_CLI}-compose > /dev/null 2>&1; then
    : ${CONTAINER_CLI_COMPOSE:="${CONTAINER_CLI}-compose"}
else
    : ${CONTAINER_CLI_COMPOSE:="${CONTAINER_CLI} compose"}
fi
infoln "Using ${CONTAINER_CLI} and ${CONTAINER_CLI_COMPOSE}"

if [ ! -d "channel-artifacts" ]; then
	mkdir channel-artifacts
fi

createChannelGenesisBlock() {
  setGlobals 1
	which configtxgen
	if [ "$?" -ne 0 ]; then
		fatalln "configtxgen tool not found."
	fi
	local bft_true=$1
	set -x

	if [ $bft_true -eq 1 ]; then
		configtxgen -profile ChannelUsingBFT -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
	else
		configtxgen -profile ChannelUsingRaft -outputBlock ./channel-artifacts/${CHANNEL_NAME}.block -channelID $CHANNEL_NAME
	fi
	res=$?
	{ set +x; } 2>/dev/null
  verifyResult $res "Failed to generate channel configuration transaction..."
}

createChannel() {
	# Poll in case the raft leader is not set yet
	local rc=1
	local COUNTER=1
	local bft_true=$1
	infoln "Adding orderers"
	while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
		sleep $DELAY
		set -x
    # . scripts/orderer1.sh ${CHANNEL_NAME}> /dev/null 2>&1
	# export PATH=${ROOTDIR}/../bin:${PWD}/../bin:$PATH
	export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer1.example.com/tls/server.crt
	export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer1.example.com/tls/server.key
	osnadmin channel join --channelID ${CHANNEL_NAME} --config-block ./channel-artifacts/${CHANNEL_NAME}.block -o localhost:8051 --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY" >> log.txt 2>&1
    if [ $bft_true -eq 1 ]; then
    #   . scripts/orderer2.sh ${CHANNEL_NAME}> /dev/null 2>&1
    #   . scripts/orderer3.sh ${CHANNEL_NAME}> /dev/null 2>&1
    #   . scripts/orderer4.sh ${CHANNEL_NAME}> /dev/null 2>&1
	#   . scripts/orderer5.sh ${CHANNEL_NAME}> /dev/null 2>&1
	#   . scripts/orderer6.sh ${CHANNEL_NAME}> /dev/null 2>&1
		for ((i=2;i<=ORDERER_NUM;i++));do	
			export ORDERER_ADMIN_TLS_SIGN_CERT=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer${i}.example.com/tls/server.crt
			export ORDERER_ADMIN_TLS_PRIVATE_KEY=${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer${i}.example.com/tls/server.key
			osnadmin channel join --channelID ${CHANNEL_NAME} --config-block ./channel-artifacts/${CHANNEL_NAME}.block -o localhost:$((8050+i)) --ca-file "$ORDERER_CA" --client-cert "$ORDERER_ADMIN_TLS_SIGN_CERT" --client-key "$ORDERER_ADMIN_TLS_PRIVATE_KEY" >> log.txt 2>&1
		done
    fi
		res=$?
		{ set +x; } 2>/dev/null
		let rc=$res
		COUNTER=$(expr $COUNTER + 1)
	done
	cat log.txt
	verifyResult $res "Channel creation failed"
}

# joinChannel ORG
joinChannel() {
  ORG=$1
  FABRIC_CFG_PATH=$PWD/../config/
  setGlobals $ORG
	local rc=1
	local COUNTER=1
	## Sometimes Join takes time, hence retry
	while [ $rc -ne 0 -a $COUNTER -lt $MAX_RETRY ] ; do
    sleep $DELAY
    set -x
    peer channel join -b $BLOCKFILE >&log.txt
    res=$?
    { set +x; } 2>/dev/null
		let rc=$res
		COUNTER=$(expr $COUNTER + 1)
	done
	cat log.txt
	verifyResult $res "After $MAX_RETRY attempts, peer0.org${ORG} has failed to join channel '$CHANNEL_NAME' "
}

setAnchorPeer() {
  ORG=$1
  . scripts/setAnchorPeer.sh $ORG $CHANNEL_NAME 
}


## User attempts to use BFT orderer in Fabric network with CA
if [ $BFT -eq 1 ] && [ -d "organizations/fabric-ca/ordererOrg/msp" ]; then
  fatalln "Fabric network seems to be using CA. This sample does not yet support the use of consensus type BFT and CA together."
fi

## Create channel genesis block
# FABRIC_CFG_PATH=$PWD/../config/
BLOCKFILE="./channel-artifacts/${CHANNEL_NAME}.block"

infoln "Generating channel genesis block '${CHANNEL_NAME}.block'"
FABRIC_CFG_PATH=${PWD}/configtx
if [ $BFT -eq 1 ]; then
  FABRIC_CFG_PATH=${PWD}/bft_config
fi
createChannelGenesisBlock $BFT


## Create channel
infoln "Creating channel ${CHANNEL_NAME}"
createChannel $BFT
successln "Channel '$CHANNEL_NAME' created"

## Join all the peers to the channel
# infoln "Joining org1 peer to the channel..."
# joinChannel 1
# infoln "Joining org2 peer to the channel..."
# joinChannel 2

for ((i=1;i<=PEER_NUM;i++));do
	infoln "Joining org${i} peer to the channel..."
	joinChannel ${i}
done



## Set the anchor peers for each org in the channel
# infoln "Setting anchor peer for org1..."
# setAnchorPeer 1
# infoln "Setting anchor peer for org2..."
# setAnchorPeer 2

for ((i=1;i<=PEER_NUM;i++));do
	infoln "Setting anchor peer for org${i}..."
	setAnchorPeer ${i}
done

successln "Channel '$CHANNEL_NAME' joined"
