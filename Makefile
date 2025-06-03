.PHONY: clean build

build: build-wetware

build-wetware:
	# phala docker build -i fbx-demo -t latest .
	docker build -t wetware/fbx-demo wetware

run:
	docker-compose up
