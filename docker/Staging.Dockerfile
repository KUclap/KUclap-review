FROM golang:1.18-buster as builder
WORKDIR /go/src/github.com/KUclap/KUclap-review

ARG GIT_ACCESS_TOKEN_CURL_CONFIG_KUCLAP_API_REVIEW
ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY
ARG AWS_DEFAULT_REGION

COPY . .

RUN go mod download

RUN curl -o config.toml https://${GIT_ACCESS_TOKEN_CURL_CONFIG_KUCLAP_API_REVIEW}@raw.githubusercontent.com/KUclap/_ENV/main/config/kuclap-review-api/config.toml
RUN mv config.toml ./config/config.toml

RUN go build -mod=readonly -v -o ./kuclap-review-api

FROM debian:buster-slim
WORKDIR /go/src/github.com/KUclap/KUclap-review

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/github.com/KUclap/KUclap-review/kuclap-review-api .
COPY --from=builder /go/src/github.com/KUclap/KUclap-review/config config/

ENV GO111MODULE=on

ENV PORT=${PORT}
ENV KIND=staging

ENV AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
ENV AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
ENV AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION}

EXPOSE 8000

CMD ./kuclap-review-api
