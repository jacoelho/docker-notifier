NAME = jacoelho/docker-notifier

build:
	docker build --no-cache -t $(NAME).dev -f Dockerfile.dev .
	docker run --rm -v $$(pwd):/build $(NAME).dev cp /go/bin/app /build
	docker build --no-cache -t $(NAME) .

run:
	docker run --rm -it -v $$(pwd):/go/src -v /var/run/docker.sock:/var/run/docker.sock golang
