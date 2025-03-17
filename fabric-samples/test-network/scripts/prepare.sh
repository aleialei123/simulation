#!/bin/bash

# 定义输出的域名与 IP 对应文件
output_file="docker_hosts_mapping.txt"

# 清空文件内容
> "$output_file"

#读取orderer节点数量
read -p "orderer节点数量:" n

# 获取所有以 Fabric 开头的容器名称
containers=$(docker ps --format "{{.Names}}" | grep "^Fabric")

# 检查是否有匹配的容器
if [[ -z "$containers" ]]; then
    echo "No containers with names starting with 'Fabric' found."
    exit 1
fi

echo "Generating domain-to-IP mapping file: $output_file"
echo "-----------------------------------------------"

# 遍历每个匹配的容器，生成域名与 IP 的对应关系
for container in $containers; do
    # 获取容器中的 eth0 接口的 IP 地址
    ip=$(docker exec "$container" ip -o -4 addr show eth0 | awk '{print $4}' | cut -d'/' -f1)

    # 根据容器名称动态生成域名
    if [[ "$container" == FabricOrdererNode-* ]]; then
        number=$(echo "$container" | sed 's/^FabricOrdererNode-//')
        domain="orderer${number}.example.com"
    elif [[ "$container" == FabricPeerNode-* ]]; then
        number=$(echo "$container" | sed 's/^FabricPeerNode-//')
        domain="peer0.org${number}.example.com"
    else
        echo "Warning: Container $container does not match known patterns, skipping..."
        continue
    fi

    # 将域名和 IP 写入到文件中
    if [[ -n "$ip" ]]; then
        echo "$ip $domain" >> "$output_file"
        echo "Added: $ip $domain"
    else
        echo "Error: Could not retrieve IP for container $container"
    fi
done

echo "Mapping file generated successfully: $output_file"

# 将文件内容加入到每个容器的 /etc/hosts 文件中
echo "Updating /etc/hosts in containers:"
echo "----------------------------------"

while IFS= read -r line; do
    for container in $containers; do
        # 添加文件内容到每个容器的 /etc/hosts 中
        docker exec "$container" bash -c "grep -q '$line' /etc/hosts || echo '$line' >> /etc/hosts"
        echo "Updated /etc/hosts in container: $container with '$line'"
    done
done < "$output_file"

echo "All containers updated successfully."

# 将文件内容加入主机的 /etc/hosts 文件中
while IFS= read -r line
do
    if ! grep -q "$line" "/etc/hosts"; then
    echo "$line" | sudo tee -a "/etc/hosts" > /dev/null
    echo "添加: $line"
else
    echo "已存在: $line"
fi
done < "$output_file"


# 构建通道
./network.sh createChannel \
            -c mychannel \
            -bft \
            -orderer_num $n \
            -peer_num 1

# 在通道上部署智能和约
./network.sh deployCC \
            -ccn secured \
            -ccp ../asset-transfer-secured-agreement/chaincode-go/ \
            -ccl go \
            -ccep "OR('Org1MSP.peer')" \
            -bft \
            -orderer_num $n \
            -peer_num 1