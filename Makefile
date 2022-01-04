NAME=missevanbot
ifeq ($(OS),Windows_NT)
 	EXECUTABLE="$(NAME).exe"
else
	ifeq ($(shell uname),Darwin)
		EXECUTABLE=$(NAME)
	else
		EXECUTABLE=$(NAME)
	endif
endif

build:
	go build -ldflags "-s -w" -o $(EXECUTABLE) cmd/main.go