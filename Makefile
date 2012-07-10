all:
	GOPATH=`pwd` go install com.mrd/types com.mrd/web gomedia

clean:
	GOPATH=`pwd` go clean -i -x com.mrd/types com.mrd/web gomedia