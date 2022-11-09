#pragma once

#include <ports/api.h>

class Core : public Api
{
public:
    virtual void SetEntity(const Entity& entity) override;
    virtual Entity GetEntity(const std::string& id) const override;
};
