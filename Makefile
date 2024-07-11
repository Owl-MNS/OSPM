install-test-utility:
	@if [ ! -e "${GOPATH}/bin/gotestsum" ]; then \
		echo "test tool gotestsum has not already been installed!"; \
		echo "installing gotestsum..."; \
		go install gotest.tools/gotestsum@latest; \
	else  \
		echo "test tool gotestsum has already been installed!"; \
	fi

test-critical: install-test-utility
	${GOPATH}/bin/gotestsum --format=testname -- -v --tags="critical" --parallel 30 ./...

test-all: test-critical