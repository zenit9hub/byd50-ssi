cmake_minimum_required(VERSION 3.22.1)
project("byd50demo")

add_library(native-bridge SHARED native_bridge.cpp)

add_library(foo SHARED IMPORTED)
set_target_properties(foo PROPERTIES
    IMPORTED_LOCATION "${CMAKE_SOURCE_DIR}/../jniLibs/${ANDROID_ABI}/libfoo.so"
    IMPORTED_SONAME "libfoo.so")

target_include_directories(native-bridge PRIVATE ${CMAKE_SOURCE_DIR})

target_link_libraries(native-bridge foo log)
