package orderusecase

type Service struct {
	uow UnitOfWork
}

func NewService(uow UnitOfWork) *Service {
	return &Service{
		uow: uow,
	}
}
