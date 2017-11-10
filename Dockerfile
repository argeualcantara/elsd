#
# (C) Copyright 2016-2017 HP Development Company, L.P.
# Confidential computer software. Valid license from HP required for possession, use or copying.
# Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
# Computer Software Documentation, and Technical Data for Commercial Items are licensed
# to the U.S. Government under vendor's standard commercial license.
#
FROM golang:alpine as builder
ADD . /go/src/github.com/hpcwp/elsd
RUN go install github.com/hpcwp/elsd/cmd/elsd

FROM alpine 
WORKDIR /root
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/elsd /root/elsd

ENTRYPOINT /root/elsd -dynamodb.addr
EXPOSE 8082 8080