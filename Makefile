BUILD_FOLDER=build
BINARY_NAME=fasthak
BIN_FOLDER=/usr/local/bin

build:
	./build.sh

run:
	./run.sh

tidy:
	go mod tidy

vendor:
	go mod vendor

deps:
	curl https://dl.filippo.io/mkcert/latest?for=linux/amd64 -o ${BIN_FOLDER}/mkcert
	chmod +x ${BIN_FOLDER}/mkcert
	mkcert -install

install:
	cp -f ${BUILD_FOLDER}/${BINARY_NAME} ${BIN_FOLDER}/

clean:
	go clean
	rm ${BUILD_FOLDER}/${BINARY_NAME}
