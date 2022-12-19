#include <core/core.h>
#include <microservice-essentials/context.h>
#include <microservice-essentials/utilities/environment.h>
#include <microservice-essentials/cross-cutting-concerns/graceful-shutdown.h>
//<<<SAPPER SECTION BEGIN MAIN-INCLUDES>>>
//<<<SAPPER SECTION END MAIN-INCLUDES>>>

int main()
{
    mse::Context::GetGlobalContext().Insert({   
//<<<SAPPER SECTION BEGIN MAIN-ENVIRONMENT>>>
//<<<SAPPER SECTION END MAIN-ENVIRONMENT>>>
            {"app", mse::getenv_or("APP", "<<<NAME>>>") },
            {"version", mse::getenv_or("VERSION", "<<<VERSION>>>") }
        });
    
    Core core;

//<<<SAPPER SECTION BEGIN MAIN-ADAPTER-INSTANTIATION>>>
//<<<SAPPER SECTION END MAIN-ADAPTER-INSTANTIATION>>>

    mse::GracefulShutdownOnSignal gracefulShutdown(mse::Signal::SIG_SHUTDOWN);

//<<<SAPPER SECTION BEGIN MAIN-HANDLE>>>
//<<<SAPPER SECTION END MAIN-HANDLE>>>

    return 0;
}
