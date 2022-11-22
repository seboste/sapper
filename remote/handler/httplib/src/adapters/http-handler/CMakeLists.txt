find_package(httplib)
find_package(nlohmann_json)
find_package(OpenSSL)

target_sources(adapters
    PRIVATE
        http-handler.h
        http-handler.cpp
    )

target_link_libraries(adapters
    PRIVATE        
        httplib::httplib 
        nlohmann_json::nlohmann_json
        OpenSSL::OpenSSL
    )
