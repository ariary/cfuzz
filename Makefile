before.build:
	go mod download && go mod vendor

build.cfuzz:
	@echo "build in ${PWD}";go build cmd/cfuzz/cfuzz.go
