build:
	go build -o bin/ghost ./cmd/ghost/

run:
	sudo ./bin/docker-mocker run /bin/sh
	
clean:
	rm -rf bin/
