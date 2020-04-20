FROM golang:latest as builder
ENV CGO_ENABLED=0
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .

FROM scratch
WORKDIR /app
COPY --from=builder /app/main .
ENTRYPOINT ["./main"]
EXPOSE 80