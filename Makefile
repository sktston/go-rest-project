.PHONY: docs
docs:
	swag init

.PHONY: build
build: docs
	go build -o output/go_app

.PHONY: run
run:
	output/go_app

.PHONY: test
test:
	go test -v
