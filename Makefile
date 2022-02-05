build:
	go build -o bin/AuthorisationServer.exe Command/AuthorisationServer/main.go
	go build -o bin/WorldServer.exe Command/WorldServer/main.go

install:
	cd Shared && go install

clean:
	rm -rf ./bin

run:
	go run main.go