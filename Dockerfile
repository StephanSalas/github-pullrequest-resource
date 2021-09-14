ARG base_image=alpine:latest
ARG builder_image=concourse/golang-builder

FROM ${builder_image} AS builder
COPY . $GOPATH/src/github.com/stephansalas/github-pullrequest-resource
WORKDIR $GOPATH/src/github.com/stephansalas/github-pullrequest-resource
ENV CGO_ENABLED 0
RUN go mod vendor
RUN go build -o /assets/out github.com/stephansalas/github-pullrequest-resource/cmd/out
RUN go build -o /assets/in github.com/stephansalas/github-pullrequest-resource/cmd/in
RUN go build -o /assets/check github.com/stephansalas/github-pullrequest-resource/cmd/check

FROM ${base_image} AS resource
RUN apk update && apk upgrade
RUN apk add --update bash tzdata ca-certificates
COPY --from=builder /assets /opt/resource

FROM resource