FROM golang:latest as builder
ENV CGO_ENABLED=0
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .

FROM alpine
WORKDIR /app
RUN apk add --no-cache curl iproute2 wget
COPY --from=builder /app/main .
ENTRYPOINT ["./main"]
EXPOSE 80
