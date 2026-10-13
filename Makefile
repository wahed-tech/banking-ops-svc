.PHONY: list help \
		fmt \
		static-check fix \
		test test-codecov-gen test-cov-html test-cov-summary

help:
	@echo "\033[1;33mMakefile\033[0m"
	@echo "\033[0;33m\thelp\033[0m - Print help"
	@echo "\033[0;33m\tlist\033[0m - Print list of all possible commands"

list:
	@$(MAKE) -qp | awk -F':' '/^[a-zA-Z0-9][^$$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /);for(i in A)print A[i]}' | grep -v Makefile | sort -u

#############
### TOOLS ###
#############
fmt:
	go fmt ./...

#####################
### STATIC CHECKS ###
#####################
static-check:
	golangci-lint run

fix:
	golangci-lint run --fix
############
### TEST ###
############
test:
	ENV=test go test -tags="unit functional integration" ./...

test-codecov-gen:
	ENV=test go test -tags="unit functional integration" -coverprofile cov.out -race -coverpkg=./... ./...
	go tool cover -html=cov.out -o coverage.html

test-codecov-gen-integration:
	#export ES_INTEGRATION_TESTS=true; export ENV=test; \
	ENV=test go test -tags="unit functional integration" -coverprofile cov.out -race -coverpkg=./... ./...
	go tool cover -html=cov.out -o coverage.html

test-cov-html:
	go tool cover -html cov.out

test-cov-summary:
	go tool cover -func cov.out
