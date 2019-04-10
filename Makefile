VERSION=v0.1.1

build:
	go build -o facedetection .

build_docker:
	docker build -t nicholasjackson/emojify-facedetection:${VERSION} .

run_docker:
	docker run -it -p 9090:9090 nicholasjackson/emojify-facedetection:${VERSION}

push_docker:
	docker push nicholasjackson/emojify-facedetection:${VERSION}

test_docker:
	curl -XPOST --data-binary @test_fixtures/group.jpg localhost:9090
