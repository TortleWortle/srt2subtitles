FROM golang:1.18.3-alpine3.16 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -v -o server ./cmd/server

FROM alpine

COPY --from=builder /app/server /app/server

EXPOSE 8080

CMD ["/app/server"]