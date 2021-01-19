build:
		go build -o bin/server .
build-linux-amd64:
		GOOS=linux GOARCH=amd64 go build -o bin/app_linux_amd64 .
build-agent-linux:
		GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64/agent agent/main.go
docker:
		docker build -t swingbylabs/node-agent .
push:
		docker push swingbylabs/node-agent