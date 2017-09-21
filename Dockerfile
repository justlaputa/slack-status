FROM octoblu/alpine-ca-certificates
MAINTAINER laputa <justlaputa@gmail.com>

COPY slack-status /slack-status

ENTRYPOINT ["/slack-status"]
