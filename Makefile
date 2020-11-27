HOST=gene

all: build upload run

build:
	GOARM=6 GOARCH=arm GOOS=linux go build -o ./bin/yupa

upload:
	scp ./bin/yupa pi@${HOST}:~/yupa

run:
	ssh -t pi@${HOST} "source ~/.profile; ~/yupa"