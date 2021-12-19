NAME="missevan-pre"
ifeq ($(OS),Windows_NT)
 	EXECUTABLE=$(NAME)+".exe"
else
	ifeq ($(shell uname),Darwin)
		EXECUTABLE=$(NAME)
	else
		EXECUTABLE=$(NAME)
	endif
endif

build:
	go build -ldflags "-s -w" -o $(EXECUTABLE) cmd/main.go
start:
	screen -S $(EXECUTABLE) ./$(EXECUTABLE)
resume:
	screen -R $(EXECUTABLE)
stop:
	screen -S $(EXECUTABLE) -X quit
restart:
	make stop && make build && make start