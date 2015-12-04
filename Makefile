build:
	GOPATH=`pwd`/. go install chief-stats/consumer

build_for_production:
	GOPATH=`pwd`/. GOARCH=amd64 go install chief-stats/consumer

deploy:  build_for_production
	scp
