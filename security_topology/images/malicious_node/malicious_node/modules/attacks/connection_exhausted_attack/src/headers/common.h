//
// Created by zhf on 2024/11/7.
//

#ifndef TCP_OFF_PATH_COMMON_H
#define TCP_OFF_PATH_COMMON_H

#include <random>
#include <tins/tins.h>

typedef std::chrono::high_resolution_clock::time_point Clock;
typedef Tins::HWAddress<6> MACAddress;

class CommandArgs {
public:
    int attack_thread_count{};
    char *interface_name{};
    char *server_ip{};
    int server_port{};
    int attack_duration{};
    CommandArgs();
    CommandArgs(char* interface_name, int attack_thread_count, char* server_ip, int server_port, int attack_duration);
};

class GlobalArgs {
public:
    CommandArgs commandArgs;
    Tins::PacketSender packetSender;
    Tins::NetworkInterface interface;
    GlobalArgs();
    void initialize(CommandArgs commandArgs);
};

extern GlobalArgs globalArgs;



#endif //TCP_OFF_PATH_COMMON_H
