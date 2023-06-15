FROM alpine:latest
RUN apk add --no-cache go
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go version
RUN go mod download
COPY *.go ./
RUN go build -o /forum-api
EXPOSE 1025/tcp
ENTRYPOINT [ "/forum-api" ]