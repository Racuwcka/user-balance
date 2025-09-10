package deposit

type provider interface {
	Deposit(userId uint32, serviceId uint16, amount uint) error
}

type Service struct {
	provider provider
}

func (s Service) Deposit(userId uint32, serviceId uint16, amount uint) error {
	return s.provider.Deposit(userId, serviceId, amount)
}

func New(p provider) *Service {
	return &Service{provider: p}
}
