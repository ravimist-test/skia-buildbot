.PHONY: build docs release push
build:
	go install  ./go/jsdoc

docs:
	npm install
	./node_modules/.bin/jsdoc -c jsdoc.config.js

release:
	CGO_ENABLED=0 GOOS=linux go install -a ./go/jsdoc
	./build_docker_release

push: release
	pushk jsdoc

run-local:
	docker run -p 8000:8000 -ti jsdoc
