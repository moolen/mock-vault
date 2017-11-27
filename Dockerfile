FROM golang:1.9.0 as builder
WORKDIR /go/src/github.com/moolen/mock-vault
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /app  
WORKDIR /app
COPY --from=builder /go/src/github.com/moolen/mock-vault/mock-vault .
CMD ["./mock-vault"]  
