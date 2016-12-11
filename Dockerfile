#
# (C) Copyright 2016-2017 HP Development Company, L.P.
# Confidential computer software. Valid license from HP required for possession, use or copying.
# Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
# Computer Software Documentation, and Technical Data for Commercial Items are licensed
# to the U.S. Government under vendor's standard commercial license.
#
FROM alpine:latest
MAINTAINER ELS Team <els-team@groups.hp.com>

# Install bash for sanity & profit
RUN apk add --update bash && rm -rf /var/cache/apk/*

# ELS binary
ADD ./els /els

# Container entrypoint scripts
ADD ./build/docker/*.sh /
RUN chmod 755 /*.sh /els

WORKDIR /
EXPOSE 7300
ENTRYPOINT ["/bin/bash", "/docker-entrypoint.sh"]
