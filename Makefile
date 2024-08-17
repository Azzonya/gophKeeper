cert:
	cd cert; ./gen.sh; cd ..

pkgs = $(shell go list ./... | grep -v /vendor | grep -v /tools | grep -v /testdata | xargs -I {} bash -c "if ls {}/**/*.go > /dev/null 2>&1; then echo {}; fi")

lint:
	@if [ -z "$(pkgs)" ]; then \
		echo "No Go files found."; \
		exit 0; \
	fi
	go install golang.org/x/lint/golint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/kisielk/errcheck@latest
	golint $(pkgs)
	go vet $(pkgs)
	staticcheck $(pkgs)
	errcheck $(pkgs)

