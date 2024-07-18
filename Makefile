pwd=$(shell pwd)

.PHONY: all taskapi unit-test

all: taskapi 

taskapi:
	CONFIG_DIR=${pwd}/configs go run ./main.go taskapi

unit-test:
	go test -count=1 -v ./...

mockgen:
	mockgen -source internal/repository/interface.go -destination=internal/repository/mocks/mock_repository.go -package=mocks