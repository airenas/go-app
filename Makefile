-include version
#####################################################################################
## print usage information
help:
	@echo 'Usage:'
	@cat ${MAKEFILE_LIST} | grep -e "^## " -A 1 | grep -v '\-\-' | sed 's/^##//' | cut -f1 -d":" | \
		awk '{info=$$0; getline; print "  " $$0 ": " info;}' | column -t -s ':' | sort 
.PHONY: help
#####################################################################################
## call units tests
test/unit: 
	go test -v -race -count=1 ./...
.PHONY: test/unit
#####################################################################################
## code vet and lint
test/lint: 
	go vet ./...
	go get -u golang.org/x/lint/golint
	golint -set_exit_status ./...
	go mod tidy ## drop lint usage in mod
.PHONY: test/lint
## creates git tag for library version
git/tag:
	git tag $(version)	
.PHONY: git/tag
## push git tag version to github
git/push/tag:
	git push origin $(version)
.PHONY: git/push/tag
#####################################################################################
## cleans all temporary data
clean:
	go clean
	go mod tidy
.PHONY: clean	
