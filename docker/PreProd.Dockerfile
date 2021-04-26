################
# GLOBAL ARGS
################

ARG GIT_ACCESS_TOKEN_CURL_CONFIG
ARG ARG_AWS_ACCESS_KEY_ID
ARG ARG_AWS_SECRET_ACCESS_KEY
ARG ARG_AWS_DEFAULT_REGION

################
# BUILDER STAGE
################

FROM golang:1.16-buster as builder
WORKDIR /go/src/github.com/KUclap/KUclap-review

ARG GIT_ACCESS_TOKEN_CURL_CONFIG

COPY . .

RUN curl -o config.toml https://${GIT_ACCESS_TOKEN_CURL_CONFIG}@raw.githubusercontent.com/KUclap/_ENV/main/config/kuclap-review-api/config.toml
RUN mv config.toml ./config/config.toml

RUN go mod download
RUN go build  -mod=readonly -v -o ./kuclap-review-api

FROM debian:buster-slim
WORKDIR /go/src/github.com/KUclap/KUclap-review

ARG ARG_AWS_ACCESS_KEY_ID
ARG ARG_AWS_SECRET_ACCESS_KEY
ARG ARG_AWS_DEFAULT_REGION

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/github.com/KUclap/KUclap-review/kuclap-review-api .
COPY --from=builder /go/src/github.com/KUclap/KUclap-review/config config/

RUN echo "${ARG_AWS_ACCESS_KEY_ID} ${ARG_AWS_SECRET_ACCESS_KEY} ${ARG_AWS_DEFAULT_REGION}" >> ENVFILE.txt

ENV GO111MODULE=on
ENV KIND=preproduction
ENV AWS_ACCESS_KEY_ID $ARG_AWS_ACCESS_KEY_ID
ENV AWS_SECRET_ACCESS_KEY $ARG_AWS_SECRET_ACCESS_KEY
ENV AWS_DEFAULT_REGION $ARG_AWS_DEFAULT_REGION

EXPOSE 8089

CMD ./kuclap-review-api