FROM golang:1.15

ARG DB_SERVER
ARG ORIGIN_ALLOWED
ARG PORT

COPY . /go/src/github.com/KUclap/KUclap-review/api
WORKDIR /go/src/github.com/KUclap/KUclap-review/api

RUN printf "DB_SERVER=$DB_SERVER\nORIGIN_ALLOWED=$ORIGIN_ALLOWED\nPORT=$PORT\n" > .env
RUN echo $GOPATH
RUN go get ./...
RUN go build 
RUN go get github.com/pilu/fresh


# CMD if [ ${APP_ENV} = production ]; \
#     then \
#     app; \
#     else \
#     go get github.com/pilu/fresh && \
#     fresh; \
#     fi

CMD ["fresh"]
EXPOSE 8000