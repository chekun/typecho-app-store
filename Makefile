compile:
	cd typecho && go install && cd ../
	go build -o bin/main
	cd command && go build -o ../bin/update && cd ../

update:
	git pull origin master

clean:
	rm -f bin/*
