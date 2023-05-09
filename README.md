# Sapper

Sapper is CLI tool that enables you to rapidly create, extend, and update C++ microservices. It is based on an ever growing template repository --- the brick repo --- with a focus on cloud native technology from which you can compose your microservices, e.g.:
- http and gRPC handler
- postgreSQL database
- JWT based authentication
- docker and kubernetes
- support for AWS, Azure, and GCP
- ...

## Foundation

Besides using C++, the microservices that can be created with sapper are constrained by the following:

### Hexagonal architecture

The structure of each microservice follows the [hexagonal](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)) -- aka ports & adapters -- architecture:
![hexagonal architecture diagram]()

| Module | Desription| Depends on |
| :------ | :-----| :---- |
| ports | Domain entities and **interfaces** that the core needs to work with.| - |
| core | Contains all **business logic** without directly depending on any infrastructure (e.g. databases). | ports |
| adapters | | ports, external libraries|
| cmd | The **main** application that brings everything together: Reads configuration, instantiates the core, instantiates the adapters, injects the adapters into the core, calls the handler. | core, adapters |

### Linux

### Conan package manager

### CMake

### Microservice-Essentials library



## Getting started


## Best practices



## Reference

sapper brick add

sapper brick list

sapper brick search

sapper service add

sapper service update

sapper service build

sapper service test

sapper service deploy

sapper remote add

sapper remote list

## Extend



## Contribute

