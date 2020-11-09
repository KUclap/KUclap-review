FROM golang:1.15

ARG GIT_ACCESS_TOKEN_CURL_CONFIG

COPY . /go/src/github.com/KUclap/KUclap-review
WORKDIR /go/src/github.com/KUclap/KUclap-review

RUN curl -o config.toml https://${GIT_ACCESS_TOKEN_CURL_CONFIG}@raw.githubusercontent.com/KUclap/_ENV/main/config/kuclap-review-api/config.toml
RUN mv config.toml ./config/config.toml && mkdir builder
RUN go get ./...
RUN go build -o ./builder/kuclap-review-api .

ENV KIND=production 
CMD ["./builder/kuclap-review-api"]

EXPOSE 8000