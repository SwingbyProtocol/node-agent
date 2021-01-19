build-linux-amd64:
		GOOS=linux GOARCH=amd64 go build -o bin/app_linux_amd64 .
build:
		go build -o bin/server .

docker:
		docker build -t swingbylabs/node-agent .
push:
		docker push swingbylabs/node-agent