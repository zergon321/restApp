all: build run

build: assets
	docker build -t rest .

assets:
	go-bindata -o assets/assets.go --pkg assets sql/...

rebuild: clean build

clean:
	docker rmi rest

run:
	docker run -it --rm -p 80:80 rest