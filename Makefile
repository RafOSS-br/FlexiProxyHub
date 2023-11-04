build:
		go build -o bin/reportIntermediate cmd/reportIntermediate/main.go
run:
		go run cmd/reportIntermediate/main.go
test:
		go test ./...
clean:
		rm -rf bin/