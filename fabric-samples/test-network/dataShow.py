import subprocess
import time
import os
import matplotlib.pyplot as plt
import json

# 设置环境变量（根据实际环境调整）
os.environ["PATH"] = os.path.join(os.getcwd(), "../bin") + ":" + os.getcwd() + ":" + os.environ["PATH"]
os.environ["FABRIC_CFG_PATH"] = os.path.join(os.getcwd(), "../config/")
os.environ["CORE_PEER_TLS_ENABLED"] = "true"
os.environ["CORE_PEER_LOCALMSPID"] = "Org1MSP"
os.environ["CORE_PEER_MSPCONFIGPATH"] = os.path.join(os.getcwd(), "organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp")
os.environ["CORE_PEER_TLS_ROOTCERT_FILE"] = os.path.join(os.getcwd(), "organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt")
os.environ["CORE_PEER_ADDRESS"] = "localhost:5051"

# 通道名称
CHANNEL_NAME = "mychannel"

# 获取初始区块高度
def get_block_height():
    result = subprocess.run(["peer", "channel", "getinfo", "-c", CHANNEL_NAME], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if result.returncode != 0:
        return None
    chain_info = result.stdout.decode("utf-8")
    print("chain_info:", chain_info)
    if chain_info.startswith('Blockchain info: '):
        chain_info = chain_info[len('Blockchain info: '):].strip()
    data = json.loads(chain_info)
    block_height = data['height']
    return block_height

BLOCK_PRE_HEIGHT = get_block_height()

# 记录启动时间
START_TIME = int(time.time())

# 设置图表
plt.ion()
fig, ax = plt.subplots()
x_vals = []
y_vals = []

# 显示图表窗口
plt.show()  # 重要：确保图表显示

# 无限循环计算出块速率
while True:
    # 获取当前区块高度
    time.sleep(1)
    BLOCK_CUR_HEIGHT = get_block_height()

    if BLOCK_CUR_HEIGHT is None:
        continue  # 如果获取区块高度失败，跳过本次循环

    CUR_TIME = int(time.time())

    # 计算区块增量（当前高度 - 上次高度）
    BLOCK_ADD_HEIGHT = BLOCK_CUR_HEIGHT - BLOCK_PRE_HEIGHT

    # 计算经过的时间（秒）
    ELAPSED_TIME = CUR_TIME - START_TIME

    # 输出每秒钟的出块速率
    print(f"经过时间: {ELAPSED_TIME} 秒, 当前区块高度: {BLOCK_CUR_HEIGHT}, 出块数量: {BLOCK_ADD_HEIGHT}, 出块速率: {BLOCK_ADD_HEIGHT} blocks/sec")

    # 更新数据
    x_vals.append(ELAPSED_TIME)
    y_vals.append(BLOCK_ADD_HEIGHT)

    # 限制数据点最多为100个
    if len(x_vals) > 100:
        x_vals.pop(0)
        y_vals.pop(0)

    # 更新图表
    ax.clear()
    ax.plot(x_vals, y_vals, label='Block Rate')
    ax.set_xlabel('Elapsed Time (seconds)')
    ax.set_ylabel('Block Rate (blocks/sec)')
    ax.set_ylim(bottom=0)  # 确保纵坐标从0开始
    ax.legend()

    plt.draw()
    plt.pause(1)  # 每秒刷新一次

    # 更新前一个区块高度
    BLOCK_PRE_HEIGHT = BLOCK_CUR_HEIGHT