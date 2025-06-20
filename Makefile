all: build run

build: build-dstack build-app build-wetware build-tiktok

build-app:
	cd app && \
		make gen
	docker build -t wetware/fbx-demo-app app

build-llm:
	cd llm && \
		make gen
	docker build -t wetware/fbx-demo-llm llm

build-tiktok:
	docker build -t tikapi/tikapi tiktok

build-wetware:
	# phala docker build -i fbx-demo -t latest .
	cd wetware && \
		make
	docker build -t wetware/fbx-demo-wetware wetware

build-dstack:
	alias git="git -c url.\"https://github.com/".insteadOf=git@github.com:\" && \
    cd dstack/sdk/simulator && \
		./build.sh

clean: clean-app clean-llm stop-dstack

clean-app:
	cd app && \
	make clean

clean-llm:
	cd llm && \
	make clean

start-dstack:
	cd dstack/sdk/simulator && \
		./dstack-simulator > /dev/null &

stop-dstack:
	-pkill -f "dstack-simulator" 2>/dev/null || true

run-dstack: stop-dstack start-dstack

run-docker:
	docker-compose up

run: run-dstack run-docker
