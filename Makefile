
version = 0.0.0


.PHONY: build 

# Build the application 
build:
	rm -rf .build/*
	mkdir -p .build/
	go build -ldflags "-X main.Version=${version}" -o .build/dvc cli/main.go

# Build and install the application 
install: build 
	go install 

## Vet the code 
vet: 
	go vet ./...