build:
		go build -o bin/FlexiProxyHub cmd/FlexiProxyHub/main.go
run:
		go run cmd/FlexiProxyHub/main.go
test:
		go test ./...
clean:
		rm -rf bin/