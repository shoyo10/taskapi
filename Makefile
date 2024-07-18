pwd=$(shell pwd)

.PHONY: all taskapi unit-test

all: taskapi 

taskapi:
	CONFIG_DIR=${pwd}/configs go run ./main.go taskapi

unit-test:
	go test -count=1 -v ./...

mockgen:
	mockgen -source internal/repository/interface.go -destination=internal/repository/mocks/mock_repository.go -package=mocks
	mockgen -source internal/service/interface.go -destination=internal/service/mocks/mock_service.go -package=mocks

docker-build:
	docker build -f build/docker/Dockerfile -t taskapiserver .

docker-run:
	docker run --name taskserver -d --rm -p 9090:9090 taskapiserver

docker-stop:
	docker stop taskserver

swagger:
	swag init -g internal/delivery/http/route.go
