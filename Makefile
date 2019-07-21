all: build run

build:
	docker build -t rest .

rebuild: clean build

clean:
	docker rmi rest

run:
	docker run -it --rm -p 80:80 rest