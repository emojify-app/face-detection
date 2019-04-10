build:
	go build -o facedetection .

build_docker:
	docker build -t nicholasjackson/emojify-facedetection:latest .

run_docker:
	docker run -it -p 9090:9090 nicholasjackson/emojify-facedetection:latest

push_docker:
	docker push nicholasjackson/emojify-facedetection:latest

test_docker:
	curl -XPOST --data-binary @test_fixtures/group.jpg localhost:9090
