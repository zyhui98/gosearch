export tag=v1.0
root:
	export ROOT=github.com/cncamp/golang

build:
	echo "building gosearch binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building gosearch container"
	docker build -t zyhui98/gosearch:${tag} .

push: release
	echo "pushing zyhui98/gosearch"
	docker push zyhui98/gosearch:${tag}
