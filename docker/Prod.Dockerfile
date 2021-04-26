FROM golang:1.15

ARG GIT_ACCESS_TOKEN_CURL_CONFIG
ARG AWS_ACCESS_KEY_ID
ARG AWS_SECRET_ACCESS_KEY

COPY . /go/src/github.com/KUclap/KUclap-review
WORKDIR /go/src/github.com/KUclap/KUclap-review

RUN curl -o config.toml https://${GIT_ACCESS_TOKEN_CURL_CONFIG}@raw.githubusercontent.com/KUclap/_ENV/main/config/kuclap-review-api/config.toml
RUN mv config.toml ./config/config.toml && mkdir builder
RUN go get ./...
RUN go build -o ./builder/kuclap-review-api .

ENV KIND=production 
ENV AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
ENV AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}

CMD ["./builder/kuclap-review-api"]

EXPOSE 8000




# FROM golang:1.15

# ARG GIT_ACCESS_TOKEN_CURL_CONFIG
# ARG AWS_ACCESS_KEY_ID
# ARG AWS_SECRET_ACCESS_KEY

# COPY . /go/src/github.com/KUclap/KUclap-review
# WORKDIR /go/src/github.com/KUclap/KUclap-review

# RUN curl -o config.toml https://${GIT_ACCESS_TOKEN_CURL_CONFIG}@raw.githubusercontent.com/KUclap/_ENV/main/config/kuclap-review-api/config.toml
# RUN mv config.toml ./config/config.toml && mkdir builder
# RUN go get ./...
# RUN go build -o ./builder/kuclap-review-api .

# ENV KIND=production 
# ENV AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
# ENV AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}

# CMD ["./builder/kuclap-review-api"]

# EXPOSE 8000


