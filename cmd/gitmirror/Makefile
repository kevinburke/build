# Copyright 2014 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

VERSION=v4
DOCKER_IMAGE_STAGE0=stage0/gitmirror:latest

stage0: *.go Dockerfile.0
	docker build --rm --force-rm -f Dockerfile.0 --tag=$(DOCKER_IMAGE_STAGE0) ../..

gitmirror: stage0
	docker run --rm $(DOCKER_IMAGE_STAGE0) base64 /go/bin/$@ | base64 -d > $@
	chmod a+x $@

docker-prod: Dockerfile gitmirror
	docker build --rm --force-rm --tag=gcr.io/symbolic-datum-552/gitmirror:$(VERSION) .
docker-dev: Dockerfile gitmirror
	docker build --rm --force-rm --tag=gcr.io/go-dashboard-dev/gitmirror:latest .

push-prod: docker-prod
	gcloud docker -- push gcr.io/symbolic-datum-552/gitmirror:$(VERSION)
push-dev: docker-dev
	gcloud docker -- push gcr.io/go-dashboard-dev/gitmirror:latest

.PHONY: clean
clean:
	$(RM) gitmirror
