OS=`uname -s`
MY_IP=`ifconfig | grep --color=none -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep --color=none -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1' | head -n 1`

setup:
	@go get github.com/DATA-DOG/godog/cmd/godog
	@go get github.com/onsi/ginkgo/ginkgo
	@glide install

test: unit deps int# acceptance

unit:
	@env MY_IP=${MY_IP} ginkgo -r --randomizeAllSpecs --randomizeSuites --cover --focus="\[Unit\].*" .
	@${MAKE} cov

integration int:
	@env MY_IP=${MY_IP} ginkgo -r --randomizeAllSpecs --randomizeSuites --cover --focus="\[Integration\].*" .

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

dependencies deps:
	@docker-compose -p deepstream-golang up -d

stop-deps:
	@docker-compose -p deepstream-golang stop
	@docker-compose -p deepstream-golang rm -f
