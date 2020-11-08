FROM golang:1.15

ARG DB_SERVER
ARG ORIGIN_ALLOWED
ARG PORT

COPY . /go/src/github.com/KUclap/KUclap-review
WORKDIR /go/src/github.com/KUclap/KUclap-review

RUN printf "DB_SERVER=$DB_SERVER\nORIGIN_ALLOWED=$ORIGIN_ALLOWED\nPORT=$PORT\n" > .env

RUN go get ./...
RUN go build 
RUN go get github.com/pilu/fresh

CMD ["fresh"]
EXPOSE 8000