//
// Created by zhf on 2024/11/10.
//
#include "utils.h"
#include <random>

int get_random_number(){
    std::random_device rd;
    std::mt19937 mt(rd());
    std::uniform_int_distribution<> dist(1, 65000);  // choosees port between 4k-9k
    return dist(mt);
}