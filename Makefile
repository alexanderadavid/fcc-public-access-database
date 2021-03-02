build:
	go build -o bin/seed seed.go

test:
	go test .

dockerRunDB:
	docker run -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=root -e MYSQL_ROOT_USER=root mysql:latest

download:
	sh download.sh

seed:
	make download
	./bin/seed