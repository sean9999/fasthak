BUILD_FOLDER=build
BINARY_NAME=fasthak
BIN_FOLDER=/usr/local/bin
REPO=github.com/sean9999/fasthak
SEMVER := $$(git tag --sort=-version:refname | head -n 1)

build:
	./build.sh

docker-build:
	docker build -t ${REPO}:${SEMVER} .

docker-run:
	docker run -p 9001:9001 -v $${PWD}/public:/srv/public ${REPO}:${SEMVER}

run:
	./run.sh

tidy:
	go mod tidy

vendor:
	go mod vendor

deps:
	curl --output-dir certs -O https://www.rec.la/rec.la-cert.crt
	curl --output-dir certs -O https://www.rec.la/rec.la-key.pem

install:
	cp -f ${BUILD_FOLDER}/${BINARY_NAME} ${BIN_FOLDER}/

clean:
	go clean
	rm ${BUILD_FOLDER}/${BINARY_NAME}
