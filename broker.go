package broker

// PersonMsg is the message that goes through the queue and represent a person
type PersonMsg struct {
	Name string
	Age  int
}

// Configuration is the configuration structure of the AMQP connection
type Configuration struct {
	AMQPConnectionURL string
}

// Config is an instance if the AMQP configuration
var Config = Configuration{
	AMQPConnectionURL: "amqp://guest:guest@broker:5672/",
}
