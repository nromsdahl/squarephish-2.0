test-darwin-arm64:
	GO111MODULE=on GOOS=darwin GOARCH=arm64 go run cmd/squarephish/main.go --config config.json -v

lint:
	golangci-lint run ./...

build-darwin: build-darwin-arm64 build-darwin-amd64

build-linux: build-linux-amd64 build-linux-arm64

build-windows: build-windows-amd64 build-windows-arm64

build-darwin-arm64:
	GO111MODULE=on GOOS=darwin GOARCH=arm64 go build -o squarephish cmd/squarephish/main.go

build-darwin-amd64:
	GO111MODULE=on GOOS=darwin GOARCH=amd64 go build -o squarephish cmd/squarephish/main.go

build-linux-amd64:
	GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o squarephish cmd/squarephish/main.go

build-linux-arm64:
	GO111MODULE=on GOOS=linux GOARCH=arm64 go build -o squarephish cmd/squarephish/main.go

build-windows-amd64:
	GO111MODULE=on GOOS=windows GOARCH=amd64 go build -o squarephish cmd/squarephish/main.go

build-windows-arm64:
	GO111MODULE=on GOOS=windows GOARCH=arm64 go build -o squarephish cmd/squarephish/main.go

run:
	./squarephish --config config.json