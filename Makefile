VERSION=v0.10.0

build:
	go build -o facedetection .

build_docker:
	docker build -t nicholasjackson/emojify-facedetection:${VERSION} .

run_docker:
	docker run -it -p 9090:9090 nicholasjackson/emojify-facedetection:${VERSION}

push_docker:
	docker push nicholasjackson/emojify-facedetection:${VERSION}

tag:
	git tag ${VERSION} && git push origin ${VERSION}

test_docker:
	curl -s -XPOST --data-binary @test_fixtures/group.jpg localhost:9090 | jq
