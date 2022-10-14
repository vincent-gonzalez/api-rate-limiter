# variables
BUILD-DIR=bin
EXE-NAME=api-rate-limiter
SOURCE-DIR=src

build:
	cd ./$(SOURCE-DIR); go build -o "../$(BUILD-DIR)/$(EXE-NAME)"

buildRun: build run

run:
	./$(BUILD-DIR)/$(EXE-NAME)
