pwd=$(shell pwd)

all: taskapi 

taskapi:
	CONFIG_DIR=${pwd}/configs go run ./main.go taskapi
