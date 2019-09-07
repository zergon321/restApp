all: build run

build: assets
	docker-compose build

assets:
	go-bindata -o assets/assets.go --pkg assets sql/...

rebuild: clean build

clean:
	docker rmi restapp_api

run:
	docker-compose up