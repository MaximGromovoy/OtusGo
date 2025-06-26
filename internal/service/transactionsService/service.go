package transactionsservice

import (
	"OtusGo/internal/interfaces"
	"OtusGo/internal/transaction"
)

type TransactionsService struct {
	repository interfaces.TransactionRepositoryInterface
}

var _ interfaces.TransactionsServiceInterface = (*TransactionsService)(nil)

func NewTransactionsService(repository interfaces.TransactionRepositoryInterface) *TransactionsService {
	return &TransactionsService{
		repository: repository,
	}
}

func (s *TransactionsService) SaveTransaction(t *transaction.Transaction) error {
	if err := s.repository.Add(t); err != nil {
		t.Fail()
		return err
	}

	t.Complete()
	if err := s.repository.Update(t); err != nil {
		return err
	}

	return nil
}
