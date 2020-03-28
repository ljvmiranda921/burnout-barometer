.PHONY: clean-build clean-test
	
build: clean-build ## install deps and build executable
	go get -v
	go build -o burnout-barometer . 
clean-build:  ## remove build artifacts
	rm -rf burnout-barometer
clean-test:  ## remove test and coverage artifacts
	rm -rf coverage.{txt,xml,json}
	rm -rf report.xml 
	rm -rf coverage/
test: clean-test ## run tests with coverage report
	go get github.com/jstemmer/go-junit-report
	go get github.com/axw/gocov/gocov
	go get github.com/AlekSi/gocov-xml
	go get -u github.com/matm/gocov-html
	go test -v ./... -coverprofile=coverage.txt -covermode count 2>&1 | go-junit-report > report.xml
	gocov convert coverage.txt > coverage.json
	gocov-xml < coverage.json > coverage.xml
	mkdir coverage
	gocov-html < coverage.json > coverage/index.html
