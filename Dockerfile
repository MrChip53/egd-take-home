# syntax=docker/dockerfile:1

FROM golang:1.21
ADD data /app/data
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-s -w' -o /app/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-s -w' -o /app/download ./cmd/downloadutility
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-s -w' -o /app/upload ./cmd/uploadutility
WORKDIR /app
CMD ["/app/server"]