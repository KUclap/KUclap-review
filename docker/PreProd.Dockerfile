FROM golang:1.16-buster as builder
WORKDIR /go/src/github.com/KUclap/KUclap-review

ARG GIT_ACCESS_TOKEN_CURL_CONFIG
ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY
ARG AWS_DEFAULT_REGION

COPY . .

RUN curl -o config.toml https://${GIT_ACCESS_TOKEN_CURL_CONFIG}@raw.githubusercontent.com/KUclap/_ENV/main/config/kuclap-review-api/config.toml
RUN mv config.toml ./config/config.toml

RUN go mod download
RUN go build  -mod=readonly -v -o ./kuclap-review-api

################
# BUILDER STAGE
################

FROM debian:buster-slim
WORKDIR /go/src/github.com/KUclap/KUclap-review

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/github.com/KUclap/KUclap-review/kuclap-review-api .
COPY --from=builder /go/src/github.com/KUclap/KUclap-review/config config/

ENV GO111MODULE=on
ENV KIND=preproduction
ENV AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
ENV AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
ENV AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION}

EXPOSE 8089

CMD ./kuclap-review-api

# FROM golang:1.15

# ARG GIT_ACCESS_TOKEN_CURL_CONFIG
# ARG AWS_ACCESS_KEY_ID
# ARG AWS_SECRET_ACCESS_KEY
# ARG AWS_DEFAULT_REGION

# COPY . /go/src/github.com/KUclap/KUclap-review
# WORKDIR /go/src/github.com/KUclap/KUclap-review

# RUN curl -o config.toml https://${GIT_ACCESS_TOKEN_CURL_CONFIG}@raw.githubusercontent.com/KUclap/_ENV/main/config/kuclap-review-api/config.toml
# RUN mv config.toml ./config/config.toml && mkdir builder
# RUN go get ./...
# RUN go build -o ./builder/kuclap-review-api .

# ENV KIND=preproduction 
# ENV AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
# ENV AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
# ENV AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION}

# CMD ["./builder/kuclap-review-api"]

# EXPOSE 8089