compile:
	cd typecho && go install && cd logger && go install && cd ../ziputil/ && go install && cd ../../
	cd controllers && go install && cd ../
	cd models && go install && cd ../
	cd routers && go install && cd ../
	go build -o bin/main
	cd command && go build -o ../bin/update && cd ../

update:
	git pull origin master

clean:
	rm -f bin/*
