# taskapi

A RESTful task API application

## Project layout

* cmd: for CLI commands
* configs: config files
* docs: swagger documents
* internal: application specific features/funcitons
* internal/delivery: network application layer
* internal/model: the place for core business logic data models
* internal/repository: data storage layer
* internal/service: core business logic functions
* pkg: general purpose features/funcitons
* pkg/config: read config
* pkg/echorouter: initialize echo router, including series middlerwares, such as set request id to header, set request id to ctx, set logger with ctx, recover. Beside, setting a centralized error handle function to process api returned error and unify error control. Also, registering pprof and swagger UI to echo router.
* pkg/errors: define series of http errors
* pkg/sqlite: initialize sqlite connection; we use the in memory mechanism in sqlite to store data
* pkg/zerolog: initialize zerolog
* build/docker: docker build application image file

## Core concept

* internal/model is the place for core business data model, and we use the models to exchange data between each layer. 
* we can notice that delivery <-> service <-> repository, when they receive input or return data, all of them use the structures defined in internal/model
* and each layer they should communicate each other by interface

## Packages

Here are the mainly used packages to build the task api server

* cobra: CLI application
* viper: read configuration
* echo: web framework
* zerolog: log
* fx: dependency injection
* sqlite: data storge
* gorm: orm
* echo-swagger: generate swagger document & UI

## Makefile

### Local run

`$ make taskapi`

### Run unit test

`$ make unit-test`

### Build docker image

`$ make docker-build`

### Docker run

`$ make docker-run`

### Generate swagger document

`$ make swagger`

## API Usage

After running the task api server, a quick way to play the APIs is to open the swagger UI: http://127.0.0.1:9090/swagger/index.html
