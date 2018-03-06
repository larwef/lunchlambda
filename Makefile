TARGET=target/lunchlambda

all: test build

test:
	go test ./...

build:
	GOOS=linux go build -o $(TARGET)
	zip $(TARGET).zip $(TARGET)

clean:
	rm -f $(TARGET) $(TARGET).zip

rebuild:
	clean all
