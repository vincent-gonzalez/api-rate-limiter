# variables
BUILD-DIR=bin
EXE-NAME=api-rate-limiter
SOURCE-DIR=src

build:
	cd ./$(SOURCE-DIR); go build -o "../$(BUILD-DIR)/$(EXE-NAME)"

buildRun: build run

run:
	./$(BUILD-DIR)/$(EXE-NAME) -rateLimit=100 -requestLimit=10

loadTest:
	vegeta attack -duration=1s -rate=1000 -targets=vegeta.conf | vegeta report
