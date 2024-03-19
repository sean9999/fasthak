FROM golang:1.22

WORKDIR /go/src/app
COPY . .
COPY public /srv/public

RUN go get -d -v ./...
RUN go install -v ./...

VOLUME /srv
EXPOSE 9001

CMD ["fasthak", "--dir=/srv/public", "--port=9001"]