#
# (C) Copyright 2016-2017 HP Development Company, L.P.
# Confidential computer software. Valid license from HP required for possession, use or copying.
# Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
# Computer Software Documentation, and Technical Data for Commercial Items are licensed
# to the U.S. Government under vendor's standard commercial license.
#
FROM golang

ADD . /go/src/github.com/galo/els-go
RUN go install github.com/galo/els-go/cmd/elsd

WORKDIR  /go/src/github.com/galo/els-go
ENTRYPOINT /go/bin/elsd

EXPOSE 8082

