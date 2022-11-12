#pragma once

#include <ports/model.h>
#include <string>

class Api
{
public:
    Api() = default;
    virtual ~Api() = default;

    virtual void SetEntity(const Entity& entity) = 0;
    virtual Entity GetEntity(const std::string& id) const = 0;
};
