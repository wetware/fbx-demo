.PHONY: clean build

all: build run

build: build-tikapi build-wetware

build-tikapi:
	docker build -t tikapi/tikapi tiktok

build-wetware:
	# phala docker build -i fbx-demo -t latest .
	cd wetware && \
		make
	docker build -t wetware/fbx-demo wetware

run:
	docker-compose up
