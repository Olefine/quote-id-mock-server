FROM golang:1.8.3 AS build
WORKDIR /go/src/github.com/olefine/quote-id-mock/
RUN go get -d -v github.com/gin-gonic/gin
RUN go get -d -v github.com/mattn/go-sqlite3
COPY . .
RUN CC=$(which gcc) go build --ldflags '-w -linkmode external -extldflags "-static"' -o app .

FROM alpine
WORKDIR /root/
COPY --from=build /go/src/github.com/olefine/quote-id-mock/app .
COPY --from=build /go/src/github.com/olefine/quote-id-mock/store/ ./store/
CMD ["./app"]
