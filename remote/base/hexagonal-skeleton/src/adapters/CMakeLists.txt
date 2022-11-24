add_library(adapters)

set_target_properties(adapters PROPERTIES LINKER_LANGUAGE CXX)

#<<<SAPPER SECTION BEGIN ADAPTERS-SUBDIRECTORIES>>>
#<<<SAPPER SECTION END ADAPTERS-SUBDIRECTORIES>>>

target_sources(adapters
    PRIVATE
    )

target_link_libraries(adapters
    PUBLIC
        ports
        microservice-essentials::microservice-essentials
    )