add_executable(<<<NAME>>>)

target_sources(<<<NAME>>>
    PRIVATE
        main.cpp
    )

target_link_libraries(<<<NAME>>>
    PUBLIC
        core
        adapters
        microservice-essentials::microservice-essentials
    )
