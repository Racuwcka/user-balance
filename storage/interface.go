package storage

type Storage interface {
	Deposit(userId uint32, serviceId uint16, amount uint) error
}
