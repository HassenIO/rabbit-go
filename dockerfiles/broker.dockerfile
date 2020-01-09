FROM rabbitmq:3.8

EXPOSE 5672 15672
RUN rabbitmq-plugins enable rabbitmq_management
