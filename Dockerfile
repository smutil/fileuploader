FROM golang:1.14 as builder

WORKDIR /src/
COPY fileuploader.go /src/
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN CGO_ENABLED=0 go build -o /bin/fileuploader

FROM alpine:3.12.1

EXPOSE 3000

COPY --from=builder /bin/fileuploader /bin/fileuploader

USER nobody

ENTRYPOINT [ "/bin/fileuploader" ]
