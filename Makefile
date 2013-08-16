all:
	go get github.com/joshholt/web
	go get github.com/joshholt/types
	GOPATH=`pwd` go install gomedia

clean:
	GOPATH=`pwd` go clean -i -x github.com/joshholt/web github.com/joshholt/types gomedia