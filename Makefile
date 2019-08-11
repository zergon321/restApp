all: build run

build: assets
	docker build -t rest .

assets:
	go-bindata -o assets/assets.go --pkg assets sql/...

rebuild: clean build

clean:
	docker rmi rest

run:
	docker run -it --rm -p 80:80 rest \
	--dbdriver="postgres" \
	--dbprotocol="postgres" \
	--dbusername="postgres" \
	--dbpassword="XXXXXX" \
	--dbhost="0.0.0.0" \
	--dbname="accounting"