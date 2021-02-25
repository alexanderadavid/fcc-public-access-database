build:
	go build -o src/bin/seed src/seed.go

dockerRunDB:
	docker run -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=root -e MYSQL_ROOT_USER=root mysql:latest --local_infile=1  

seed:
	cd src/; sh download.sh;
	cd src/; ./bin/seed;