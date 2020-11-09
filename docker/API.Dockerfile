FROM golang:1.15

ARG GITHUB_TOKEN

COPY . /go/src/github.com/KUclap/KUclap-review
WORKDIR /go/src/github.com/KUclap/KUclap-review

RUN curl -H 'Authorization: token ${GITHUB_TOKEN}' -o config.toml https://raw.githubusercontent.com/KUclap/_ENV/main/config/config.toml
RUN cat config.toml
RUN mv config.toml ./config/config.toml
RUN mkdir builder
RUN go get ./...
RUN go build -o ./builder/kuclap-review-api .
RUN ls -al ./builder
RUN ls -al

CMD ["KIND=production ./builder/kuclap-review-api"]
EXPOSE 8000

# FROM golang:1.15

# ARG DB_SERVER
# ARG ORIGIN_ALLOWED
# ARG PORT

# COPY . /go/src/github.com/KUclap/KUclap-review
# WORKDIR /go/src/github.com/KUclap/KUclap-review

# RUN printf "DB_SERVER=$DB_SERVER\nORIGIN_ALLOWED=$ORIGIN_ALLOWED\nPORT=$PORT\n" > .env

# RUN go get ./...
# RUN go build -o ./builder/kuclap-review-api .
# # RUN go build 
# # RUN go get github.com/pilu/fresh

# CMD ["KIND=production ./builder/kuclap-review-api"]
# EXPOSE 8000