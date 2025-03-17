cmake_minimum_required(VERSION 3.16)
project(connection_exhausted_attack)

set(CMAKE_CXX_STANDARD 11)

add_executable(connection_exhausted_attack
        src/main.cpp
        src/common/common.cpp
        src/connection/connection.cpp
        src/utils/utils.cpp)

set(HEADERS_DIR
        /usr/local/include
        src/headers)

# CLion IDE will find symbols from <linux/*>
target_include_directories("connection_exhausted_attack" PRIVATE ${HEADERS_DIR}) # clion -> tools -> resync with remote hosts
target_link_libraries("connection_exhausted_attack" PRIVATE tins pthread)

