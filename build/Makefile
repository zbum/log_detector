GO = go

clean:
	rm -rf release
	mkdir -p release

build:
	env GOOS=darwin GOARCH=amd64  $(GO) build -o release/log_detector.darwin.amd64 ../main.go
	env GOOS=darwin GOARCH=arm64  $(GO) build -o release/log_detector.darwin.arm64 ../main.go
	env GOOS=windows GOARCH=amd64 $(GO) build -o release/log_detector.amd64.exe ../main.go
	env GOOS=linux GOARCH=amd64   $(GO) build -o release/log_detector.linux.amd64 ../main.go

all: build