#include <core/core.h>
#include <microservice-essentials/context.h>
#include <microservice-essentials/utilities/environment.h>
#include <microservice-essentials/cross-cutting-concerns/graceful-shutdown.h>

int main()
{
    mse::Context::GetGlobalContext().Insert({ 
            {"app", mse::getenv_or("APP", "<<<NAME>>>") },
            {"version", mse::getenv_or("VERSION", "<<<VERSION>>>") }
        });
    
    Core core;

    mse::GracefulShutdownOnSignal gracefulShutdown(mse::Signal::SIG_SHUTDOWN);

    return 0;
}
