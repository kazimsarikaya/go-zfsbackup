.PHONY: docker build

all: build

build:
	./build.sh $(filter-out $@,$(MAKECMDGOALS))

%:
	@:
