
version = `git describe --tags`


.PHONY: build 

# Build the application 
# rm -rf .build/*
build:
	mkdir -p .build/
	go build -ldflags "-X main.Version=${version}" -o .build/dvc cli/main.go

# Build and install the application 
install: build 
	go install 

## Vet the code 
vet: 
	go vet ./...