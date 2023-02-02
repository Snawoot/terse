FROM golang:1.19 AS build

ARG GIT_DESC=undefined

WORKDIR /go/src/github.com/Snawoot/terse
COPY . .
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-s -w -extldflags "-static" -X main.version='"$GIT_DESC" ./cmd/terse

FROM scratch AS arrange

FROM scratch
COPY --from=build /go/src/github.com/Snawoot/terse/terse /
USER 9999:9999
ENTRYPOINT ["/terse"]
