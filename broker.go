package broker

type AddTask struct {
	Number1 int
	Number2 int
}

type Configuration struct {
	AMQPConnectionURL string
}

var Config = Configuration{
	AMQPConnectionURL: "amqp://guest:guest@broker:5672/",
}
