rabbit:
	docker run --hostname my-rabbit -p 5672:5672 -p 15672:15672 --name some-rabbit rabbitmq:3.8
start:
	docker-compose up --build