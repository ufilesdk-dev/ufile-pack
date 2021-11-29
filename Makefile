

EXEC_NAME = US3-PACK

all: build_main
	 
build_main:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(EXEC_NAME) main.go

clean:
	rm ${EXEC_NAME}
