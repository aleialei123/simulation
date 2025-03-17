//
// Created by zhf on 2024/11/7.
//
#include "common.h"

CommandArgs::CommandArgs() = default;

CommandArgs::CommandArgs(char* interface_name, int attack_thread_count, char* server_ip, int server_port, int attack_duration) {
    this->interface_name = interface_name;
    this->attack_thread_count = attack_thread_count;
    this->server_ip = server_ip;
    this->server_port = server_port;
    this->attack_duration = attack_duration;
}

GlobalArgs::GlobalArgs() = default;

void GlobalArgs::initialize(CommandArgs commandArgsTmp){
    this->commandArgs = commandArgsTmp;
    this->packetSender = Tins::PacketSender(this->commandArgs.interface_name);
    this->interface = Tins::NetworkInterface(this->commandArgs.interface_name);
}
