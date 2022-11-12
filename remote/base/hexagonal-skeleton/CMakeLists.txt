cmake_minimum_required (VERSION 3.16.3)
include(${CMAKE_BINARY_DIR}/conan_paths.cmake)

project (<<<NAME>>> VERSION <<<VERSION>>> LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED True)

set(CMAKE_PREFIX_PATH ${CMAKE_BINARY_DIR} ${CMAKE_PREFIX_PATH})

option(BUILD_TESTING "build the tests" OFF)

find_package(microservice-essentials REQUIRED)

add_subdirectory (src)

if(BUILD_TESTING)    
    add_subdirectory (tests)
endif()
