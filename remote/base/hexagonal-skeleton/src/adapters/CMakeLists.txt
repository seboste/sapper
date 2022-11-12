add_library(adapters)

set_target_properties(adapters PROPERTIES LINKER_LANGUAGE CXX)

target_sources(adapters
    PRIVATE
    )

target_link_libraries(adapters
    PUBLIC
        ports
    )