
version = `git describe --tags`


.PHONY: build install

# Build the application 
# rm -rf .build/*
build:
	mkdir -p .build/
	go build -ldflags "-X main.Version=${version}" -o .build/dvc cli/dvc/main.go

# Build and install the application 
install: 
	cd cli/dvc; go install -ldflags "-X main.Version=${version}"

## Vet the code 
vet: 
	go vet ./...