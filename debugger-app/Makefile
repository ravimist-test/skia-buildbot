KGO := CGO_ENABLED=0 GOOS=linux go build

node_modules: package.json
	npm ci

.PHONY: test
test: go-test js-test
	true

.PHONY: go-test
go-test:
	go test ./go/...

.PHONY: js-test
js-test:
	xvfb-run --auto-servernum --server-args "-screen 0 1280x1024x24" npx karma start --single-run

.PHONY: js-test-debug
js-test-debug:
	npx karma start --no-single-run

.PHONY: puppeteer-test
puppeteer-test:
	NODE_OPTIONS="--max-old-space-size=$(NODE_MEMORY_LIMIT)" npx mocha -r ts-node/register ./**/*_puppeteer_test.ts

.PHONY: puppeteer-test-debug
puppeteer-test-debug:
	NODE_OPTIONS="--max-old-space-size=$(NODE_MEMORY_LIMIT)" npx mocha -r ts-node/register ./**/*_puppeteer_test.ts --inspect-brk

.PHONY: serve
serve:
	npx webpack-dev-server --watch-poll --mode=development

.PHONY: wasm_libs
wasm_libs:
	rm -rf build
	mkdir build
	docker run -u $$(id -u ${USER}):$$(id -g ${USER}) --rm \
	       --volume `pwd`/build:/OUT gcr.io/skia-public/skia-wasm-release:prod \
	       sh -c "cp -r /tmp/* /OUT"
	echo "export const SKIA_VERSION = '`cat build/VERSION`';" > build/version.js

.PHONY: build-frontend-ci
build-frontend-ci: wasm_libs
	npm ci
	npx webpack --mode=production

.PHONY: test-frontend-ci
test-frontend-ci: wasm_libs
	npm ci
	xvfb-run --auto-servernum --server-args "-screen 0 1280x1024x24" npx karma start --single-run

release: wasm_libs
	npx webpack --mode=production
	CGO_ENABLED=0 GOOS=linux go install -a ./go/debugger-app
	./build_release

.PHONY: release_ci
release_ci:
	npm ci
	npx webpack --mode=production
	CGO_ENABLED=0 GOOS=linux go install -a ./go/debugger-assets
	./build_docker_release