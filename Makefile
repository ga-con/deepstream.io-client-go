setup:
	@go get github.com/DATA-DOG/godog/cmd/godog
	@go get github.com/onsi/ginkgo/ginkgo
	@glide install

test: unit acceptance

unit:
	@ginkgo -r --cover .
	@${MAKE} cov

test-coverage-run:
	@mkdir -p _build
	@-rm -rf _build/test-coverage-all.out
	@echo "mode: count" > _build/test-coverage-all.out
	@bash -c 'for f in $$(find . -name "*.coverprofile"); do tail -n +2 $$f >> _build/test-coverage-all.out; done'

cov: test-coverage-run
	@echo
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@echo "Functions Coverage"
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@echo
	@go tool cover -func=_build/test-coverage-all.out
	@echo
	@echo "=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-"
	@echo

cov-html: test-coverage-run
	@go tool cover -html=_build/test-coverage-all.out

acceptance:
	@godog
