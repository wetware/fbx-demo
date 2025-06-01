.PHONY: clean build

build:
	# phala docker build -i fbx-demo -t latest .
	docker build -t wetware/fbx-demo .

run:
	docker-compose up
