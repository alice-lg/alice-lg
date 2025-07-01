
.PHONY: clean

all: deps
	yarn build

dev:
	yarn start

deps:
	yarn install

test: deps
	yarn test --ci --all --watchAll=false

clean:
	rm -rf ./build

