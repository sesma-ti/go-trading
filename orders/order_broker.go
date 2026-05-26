package orders

type OrderBroker interface {
	Send(order Order)
}
