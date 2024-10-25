
version = `git describe --tags`


.PHONY: build install

# Build the application 
# rm -rf .build/*
build:
	mkdir -p .build/
	go build -ldflags "-X main.Version=${version}" -o .build/dvc cmd/dvc/main.go

# Build and install the application 
install: 
	cd cmd/dvc; go install -ldflags "-X main.Version=${version}"

tag: 
	./scripts/tag.sh $(versionType)

version: tag


## Vet the code 
vet: 
	go vet ./...
