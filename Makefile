.PHONY: build-frontend
build-frontend: 
	${MAKE} -C ui/desktop build

.PHONY: generate
generate:
	pkger

.PHONY: build-cli
build-cli: loophole
	go build -tags cli,skippkger -o loophole .

.PHONY: build-desktop generate
build-dektop: build-frontend loophole-desktop
	go build -tags desktop -o loophole-desktop .