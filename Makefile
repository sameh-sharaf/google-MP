default: build

privateinit:
	mkdir -p .libs/bitbucket.org/ascensionlab/
	rm -rf .libs/bitbucket.org/ascensionlab/*
	git clone --recursive git@bitbucket.org:ascensionlab/internalhttp.git .libs/bitbucket.org/ascensionlab/internalhttp
	git clone --recursive git@bitbucket.org:ascensionlab/internalhelpers.git .libs/bitbucket.org/ascensionlab/internalhelpers

privatedeps:
	mkdir -p $(GOPATH)/src/bitbucket.org/ascensionlab
	make -f .libs/bitbucket.org/ascensionlab/internalhttp/Makefile deps
	make -f .libs/bitbucket.org/ascensionlab/internalhelpers/Makefile deps
	mv .libs/bitbucket.org/ascensionlab/* $(GOPATH)/src/bitbucket.org/ascensionlab/
	rm -rf .libs/
	go install bitbucket.org/ascensionlab/internalhttp
	go install bitbucket.org/ascensionlab/internalhelpers

deps:
	go get -u github.com/aws/aws-sdk-go/aws
	go get -u github.com/aws/aws-sdk-go/aws/session
	go get -u github.com/aws/aws-sdk-go/service/dynamodb
	go get -u github.com/aws/aws-sdk-go/service/sqs
	go get -u github.com/bugsnag/bugsnag-go

bin/cron-google-mp: *.go
	go build -o bin/cron-google-mp $^

build: bin/cron-google-mp

run: build
	bin/cron-google-mp

clean:
	rm -f bin/cron-google-mp
	rmdir bin
