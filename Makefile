build-linux-amd64:
		GOOS=linux GOARCH=amd64 go build -o bin/app_linux_amd64 .
build:
		go build -o bin/server .
build-agent-linux:
		GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64/agent .

docker:
		docker build -t swingbylabs/node-agent .
push:
		docker push swingbylabs/node-agent