package rabbitmq

const (
	AsyncTransfeEnable   = true
	RabbitURL            = "amqp://guest:guest@192.168.200.217:5672/"
	TransExchangeName    = "uploadserver.trans"
	TransCOSQueueName    = "uploadserver.trans.cos"
	TransCOSErrQueueName = "uploadserver.trans.cos.err"
	TransOSSRoutingKey   = "cos"
)
