#building
FROM golang:latest as builder
ADD . /app/
WORKDIR /app
ENV GOPROXY="http://192.168.88.35:8081/repository/Go/"
RUN go get github.com/gorilla/mux
RUN go get github.com/jackc/pgx/v4/pgxpool
CMD echo $GOPROXY
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /BooksApp .

#packaging
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /BooksApp ./
EXPOSE 8080
ENTRYPOINT ["./BooksApp"]
