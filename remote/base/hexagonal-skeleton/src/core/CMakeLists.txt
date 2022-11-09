add_library(core)

target_sources(core
    PRIVATE
        core.h
        core.cpp
    )

target_link_libraries(core
    PUBLIC
        ports        
    )