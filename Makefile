IMAGE_VERSION=0.1
IMAGE_NAME=cwa-verification-fake
IMAGE_ID=$(IMAGE_NAME):$(IMAGE_VERSION)

ACCOUNT?="please-set-env-var"
AUTHOR?="please-set-env-var"

all-push-to-github-registry: build tag push

build:
	docker build -t $(IMAGE_ID) .

tag:
	docker tag $(IMAGE_ID) docker.pkg.github.com/$(ACCOUNT)/$(IMAGE_NAME)/$(IMAGE_ID)
	docker tag $(IMAGE_ID) docker.pkg.github.com/$(ACCOUNT)/$(IMAGE_NAME)/$(IMAGE_NAME):latest

login:
	test -f "./GH_TOKEN.txt" || { echo "add ./GH_TOKEN.txt with a personal access token"; exit 1;}
	cat ./GH_TOKEN.txt | docker login docker.pkg.github.com -u $(AUTHOR) --password-stdin

push: login
	docker push docker.pkg.github.com/$(ACCOUNT)/$(IMAGE_NAME)/$(IMAGE_ID)
	docker push docker.pkg.github.com/$(ACCOUNT)/$(IMAGE_NAME)/$(IMAGE_NAME):latest
