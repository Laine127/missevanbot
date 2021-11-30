ifeq ($(OS),Windows_NT)
 	EXECUTABLE="missevan.exe"
else
	ifeq ($(shell uname),Darwin)
		EXECUTABLE="missevan"
	else
		EXECUTABLE="missevan"
	endif
endif

build:
	go build -ldflags "-s -w" -o $(EXECUTABLE) cmd/main.go
run:
	screen -S missevan ./missevan
resume:
	screen -R missevan
stop:
	screen -S missevan -X quit
restart:
	make stop && make build && make run