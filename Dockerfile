# Building binary
FROM alpine:edge AS go
WORKDIR /go
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk add --no-cache go git musl-dev
COPY . ./
RUN go mod vendor
RUN go build -o musebot *.go


## Create Release
FROM alpine:edge
WORKDIR /usr/local/bin/
RUN apk add ca-certificates && mkdir "/storage/"
COPY --from=go /go/musebot .
ENV WEB_LISTEN ":8080"
CMD ["./musebot"]
