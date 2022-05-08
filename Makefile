BUILD_FOLDER=build
BINARY_NAME=fasthak
BIN_FOLDER=/usr/local/bin

build:
	go build -o ${BUILD_FOLDER}/

run:
	go run *.go --dir=public --port=9443 

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
