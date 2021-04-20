FROM golang:1.16-buster as builder
WORKDIR /go/src/github.com/KUclap/KUclap-review

ARG GIT_ACCESS_TOKEN_CURL_CONFIG=${GIT_ACCESS_TOKEN_CURL_CONFIG}

COPY . .

RUN go mod download
RUN curl -o config.toml https://${GIT_ACCESS_TOKEN_CURL_CONFIG}@raw.githubusercontent.com/KUclap/_ENV/main/env/kuclap-review-api/config.toml
RUN mv config.toml config/

# RUN env GOOS=darwin GOARCH=arm64 go build -mod=readonly -o ./kuclap-review-api -v .
# RUN env GOOS=darwin GOARCH=arm go build -mod=readonly -o ./kuclap-review-api -v .
RUN go build -mod=readonly -v -o ./kuclap-review-api

FROM debian:buster-slim
# FROM golang:1.16-buster
WORKDIR /go/src/github.com/KUclap/KUclap-review
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/github.com/KUclap/KUclap-review/kuclap-review-api .
COPY --from=builder /go/src/github.com/KUclap/KUclap-review/config config/

ENV GO111MODULE=on
ENV PORT=${PORT}
ENV KIND=staging

EXPOSE 8000

CMD ./kuclap-review-api