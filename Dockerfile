FROM golang:alpine AS BUILDER

RUN apk add --no-cache git

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -v -o main .

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

COPY --from=BUILDER /go/src/app/main .

EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=20s --retries=3 \
  CMD wget http://localhost:${PORT}/ -q -O - > /dev/null 2>&1

CMD ["./main"]
