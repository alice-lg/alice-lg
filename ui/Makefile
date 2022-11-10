
.PHONY: clean

all: deps
	yarn build

deps:
	yarn install

test: deps
	yarn test --ci --all --watchAll=false

clean:
	rm -rf ./build

