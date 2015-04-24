IMAGE = jacoelho/docker-notifier

.PHONY: run build container

all: container

ca-certificates.crt:
	docker run --rm debian /bin/sh -c "apt-get update && apt-get install -y ca-certificates && cat /etc/ssl/certs/ca-certificates.crt" > ca-certificates.crt

build:
	docker build -t build -f Dockerfile.build .
	docker run --rm build /bin/sh -c "cat /go/bin/app" > app
	chmod +x app

container: build ca-certificates.crt
	docker build -t $(IMAGE) .

run:
	docker run --rm -it -v $$(pwd):/go/src -v /var/run/docker.sock:/var/run/docker.sock golang
