setup:
	@go get github.com/DATA-DOG/godog/cmd/godog
	@go get github.com/onsi/ginkgo/ginkgo
	@glide install

test: unit acceptance

unit:
	@ginkgo -r --cover .

acceptance:
	@godog
