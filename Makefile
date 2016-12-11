#
# (C) Copyright 2016-2017 HP Development Company, L.P.
# Confidential computer software. Valid license from HP required for possession, use or copying.
# Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
# Computer Software Documentation, and Technical Data for Commercial Items are licensed
# to the U.S. Government under vendor's standard commercial license.
#
all: container

ENVVAR=GOARCH=amd64 CGO_ENABLED=0
TAG=1.5.0

build: clean
	$(ENVVAR) go build -o els

build-linux: clean
	$(ENVVAR) GOOS=linux go build -o els

container: build-linux
	docker build -t cwp/els:$(TAG) .

clean-container: build-linux
	docker build --no-cache --force-rm -t cwp/els:$(TAG) .

clean:
	rm -f els

.PHONY: all build container clean clean-container c
