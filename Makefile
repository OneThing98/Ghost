build:
	go build -o bin/docker-mocker ./cmd/docker-mocker/

run:
	sudo ./bin/docker-mocker run /bin/sh
	
clean:
	rm -rf bin/
