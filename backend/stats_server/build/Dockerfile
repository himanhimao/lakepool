FROM golang:1.11 AS builder

WORKDIR /tools
# Install glide
RUN wget https://github.com/Masterminds/glide/releases/download/v0.13.2/glide-v0.13.2-linux-amd64.tar.gz
RUN tar -zxvf glide-v0.13.2-linux-amd64.tar.gz
RUN mv linux-amd64/ glide/
RUN mv  glide/glide /usr/local/bin/
RUN rm -rf /tools

WORKDIR $GOPATH/src/github.com/himanhimao/lakepool/backend/stats_server
COPY glide.yaml ./
COPY glide.lock ./
ENV GLIDE_HOME $GOPATH/src/github.com/himanhimao/lakepool/backend/stats_server
RUN glide install -v
COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /go/bin/stats-server cmd/server/server.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /go/bin/stats-pusher cmd/pusher/pusher.go

FROM alpine:latest
RUN apk add --update ca-certificates
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /go/bin/stats-server /usr/local/bin
COPY --from=builder /go/bin/stats-pusher /usr/local/bin