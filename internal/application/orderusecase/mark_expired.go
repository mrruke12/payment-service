package orderusecase

import "context"

func (s *Service) MarkExpired(ctx context.Context) error {
	const batchSize = 200

	return s.uow.Do(ctx, func(ctx context.Context, r RepositoryProvider) error {
		repo := r.ReservationRepository()

		return repo.TidyReservations(ctx, batchSize)
	})
}
