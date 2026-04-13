.PHONY: build-frontend
build-frontend: 
	${MAKE} -C ui/desktop build

.PHONY: build-cli
build-cli:
	go build -tags cli -o loophole .

.PHONY: build-desktop
build-desktop: build-frontend
	go build -tags desktop -o loophole-desktop .
