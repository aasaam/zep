FROM golang:1-alpine AS builder

ADD . /src

RUN cd /src \
  && go mod tidy \
  && CGO_ENABLED=0 go build -o zep -ldflags '-w -s' . \
  && ls -lah /src/zep

FROM scratch

COPY --from=builder /src/zep /usr/bin/zep

ENTRYPOINT ["/usr/bin/zep"]
