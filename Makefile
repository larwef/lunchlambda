TARGET=target

all: test build

test:
	go fmt ./...
	go test ./...

build:
	GOOS=linux go build -o $(TARGET)/main
	zip -j $(TARGET)/deployment.zip $(TARGET)/main

clean:
	rm -rf $(TARGET)

rebuild:
	clean all
