.PHONY: all clean build run run-docker start-dstack stop-dstack

all: build run

build: build-dstack build-tikapi build-wetware

build-tikapi:
	docker build -t tikapi/tikapi tiktok

build-wetware:
	# phala docker build -i fbx-demo -t latest .
	cd wetware && \
		make
	docker build -t wetware/fbx-demo wetware

build-dstack:
	cd dstack/sdk/simulator && \
		./build.sh

clean: stop-dstack

start-dstack:
	cd dstack/sdk/simulator && \
		./dstack-simulator > /dev/null &

stop-dstack:
	-pkill -f "dstack-simulator" 2>/dev/null || true

run-dstack: stop-dstack start-dstack

run-docker:
	docker-compose up

run: run-dstack run-docker
