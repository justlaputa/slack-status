FROM golang:1.12-alpine as build
WORKDIR /go/src/app
COPY . .
RUN go build -v -o slack-status .
RUN ls -lt

FROM alpine
LABEL Author="laputa"
LABEL Email="<justlaputa@gmail.com>"

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /root
COPY --from=build /go/src/app/slack-status /root/slack-status

ENTRYPOINT ["/root/slack-status"]
