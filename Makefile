DOCKER_IMAGE="registry.met.no/modellprod/roadlabels:v0.7.5"

build:
	go build -ldflags "-X main.h2s3accessKey=$(S3AccessKey) -X main.h2s3secretKey=$(S3SecretKey) -X main.buildTime=$$(date +'%Y-%m-%dT%H:%MZ') -X main.version=$$(git log --pretty=format:'%h' -n 1)" -o roadlabels roadlabels.go

dockerimg:
	docker build --no-cache  -t="$(DOCKER_IMAGE)" --build-arg S3SecretKey=$(S3SecretKey) --build-arg S3AccessKey=$(S3AccessKey) .

dockerrun:
	docker run --user 1000 -v /lustre:/lustre -v $(PWD)/var/lib/roadlabels:/var/lib/roadlabels -p 25260:25260 -i -t $(DOCKER_IMAGE)

dockerpush:
	docker image push $(DOCKER_IMAGE)

