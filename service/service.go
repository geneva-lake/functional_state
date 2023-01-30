package service

import (
	"context"
	"fmt"
	"time"

	"github.com/geneva-lake/functional_state/general"
	"github.com/google/uuid"
)

//   - -------------------------------------------------------------------------------------------------------------------
//     User service
//   - -------------------------------------------------------------------------------------------------------------------
type Service struct {
	userFlow *UserFlow
	config   *Config
	repo     *Repository
	ctx      context.Context
}

func NewService(config *Config) *Service {
	repo := NewRepository(general.NewPgsql(config.DBConnectionString))
	s := Service{
		config:   config,
		userFlow: new(UserFlow),
		repo:     repo,
	}
	return &s
}

func (s *Service) WithContext(ctx context.Context) *Service {
	s.ctx = ctx
	return s
}

func (s *Service) WithUserID(id uuid.UUID) *Service {
	s.userFlow.ID = id
	return s
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Getting user from storage and processing errors
//   - -------------------------------------------------------------------------------------------------------------------
func (s *Service) FromStorage() *Service {
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()
	user, err := s.repo.UserGetByID(ctx, s.userFlow.ID)
	if err != nil {
		s.userFlow.Err = err
		s.userFlow.Status = StoredUserInternalError
		return s
	}
	if user == nil {
		s.userFlow.Status = StoredUserNotFound
		return s
	}
	s.userFlow.Status = StoredUserFound
	s.userFlow.Name = user.Name
	s.userFlow.Email = user.Email
	return s
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Getting financials information and statuses algebra
//   - -------------------------------------------------------------------------------------------------------------------
func (s *Service) FromFinancials() *Service {
	if s.userFlow.Status != StoredUserFound {
		return s
	}
	ctx, cancel := context.WithTimeout(s.ctx, 5*time.Second)
	defer cancel()
	url := fmt.Sprintf(s.config.FiancialsURL, s.userFlow.ID.String())
	resp, err := general.MakeHTTPRequest[interface{}, UserFinancialsResponse](ctx, "GET", url, nil)
	if err != nil {
		s.userFlow.Status = TransactionsNotAvailable
		s.userFlow.Err = err
		return s
	}
	if resp.Status == general.StatusError {
		switch resp.Error {
		case FinancialsUserNotFound:
			s.userFlow.Status = TransactionsNotFound
		case FinancialsInternalerror:
			s.userFlow.Status = TransactionsNotAvailable
		}
		return s
	}
	s.userFlow.Status = TransactionsReceived
	s.userFlow.Balance = resp.Result.Balance
	for _, tr := range resp.Result.Transactions {
		s.userFlow.Transactions = append(s.userFlow.Transactions, UserTransaction(tr))
	}
	return s
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Returning flow object
//   - -------------------------------------------------------------------------------------------------------------------
func (s *Service) Answer() *UserFlow {
	return s.userFlow
}
