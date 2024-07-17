install-test-utility:
	@if [ ! -e "${GOPATH}/bin/gotestsum" ]; then \
		echo "test tool gotestsum has not already been installed!"; \
		echo "installing gotestsum..."; \
		go install gotest.tools/gotestsum@latest; \
	else  \
		echo "test tool gotestsum has already been installed!"; \
	fi

test-critical: install-test-utility
	${GOPATH}/bin/gotestsum --format=testname -- -cover -v --tags="critical" ./...

test-all: test-critical

api-doc:
	@if [ ! -e "${GOPATH}/bin/swag" ]; then \
		echo "api doc generator "swag" has not already been installed!"; \
		echo "installing swag..."; \
		go get -u github.com/swaggo/swag/cmd/swag; \
	else  \
		echo "api doc generator "swag" has already been installed!"; \
	fi

	${GOPATH}/bin/swag init -g cmd/main.go --parseInternal --dir ./,internal/api/handler/  -o docs/api/


git-push: api-doc
	@git add --all
	
	# each time I make a change, I write it down on commit.txt file which is ignored by .gitignore file
	@git commit -F commit.txt 
	
	@git push