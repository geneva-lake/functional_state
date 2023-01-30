package service

import (
	"time"

	"github.com/geneva-lake/functional_state/general"
	"github.com/google/uuid"
)

//   - -------------------------------------------------------------------------------------------------------------------
//     Error status for initiatory request processing
//   - -------------------------------------------------------------------------------------------------------------------
type UserError string

const (
	NoUserID      UserError = "no_user_id"
	WrongUserID   UserError = "wrong_user_id"
	UserNotFound  UserError = "user_not_found"
	InternalError UserError = "internal_error"
)

//   - -------------------------------------------------------------------------------------------------------------------
//      General api response
//   - -------------------------------------------------------------------------------------------------------------------
type UserResponse struct {
	Status general.ResponseStatus
	Error  UserError
	Result *UserDTO
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Dto for user
//   - -------------------------------------------------------------------------------------------------------------------
type UserDTO struct {
	ID                uuid.UUID
	Name              string
	Email             string
	Balance           uint64
	TransactionStatus UserDTOTransactionsStatus
	Transactions      []UserTransaction
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Status for user transactions
//   - -------------------------------------------------------------------------------------------------------------------
type UserDTOTransactionsStatus string

const (
	UserDTOTransactionsReceived     UserDTOTransactionsStatus = "transactions_received"
	UserDTOTransactionsNotFound     UserDTOTransactionsStatus = "transactions_not_found"
	UserDTOTransactionsNotAvailable UserDTOTransactionsStatus = "transactions_not_available"
)

func NewUserDTO(f *UserFlow) *UserDTO {
	dto := UserDTO{
		ID:    f.ID,
		Name:  f.Name,
		Email: f.Email,
	}
	switch f.Status {
	case TransactionsReceived:
		dto.Balance = f.Balance
		dto.Transactions = f.Transactions
		dto.TransactionStatus = UserDTOTransactionsReceived
	case TransactionsNotFound:
		dto.TransactionStatus = UserDTOTransactionsNotFound
	case TransactionsNotAvailable:
		dto.TransactionStatus = UserDTOTransactionsNotAvailable
	}
	return &dto
}

//   - -------------------------------------------------------------------------------------------------------------------
//     State of user processing
//   - -------------------------------------------------------------------------------------------------------------------
type UserFlowStatus uint

const (
	StoredUserFound UserFlowStatus = iota + 1 // user is found in database
	StoredUserNotFound // user not found in database
	StoredUserInternalError // error when reading from database
	TransactionsReceived // financial transactions successfully received
	TransactionsNotFound // there is no transaction of this user
	TransactionsNotAvailable // error occured when interaction with financials service 
)

//   - -------------------------------------------------------------------------------------------------------------------
//     Struct containing flow state and user information 
//   - -------------------------------------------------------------------------------------------------------------------
type UserFlow struct {
	Status       UserFlowStatus
	ID           uuid.UUID
	Name         string
	Email        string
	Balance      uint64
	Transactions []UserTransaction
	Err          error
}

//   - -------------------------------------------------------------------------------------------------------------------
//     User transactions
//   - -------------------------------------------------------------------------------------------------------------------
type UserTransaction struct {
	CreatedDT time.Time
	Amount    uint64
	Reason    string
}

//   - -------------------------------------------------------------------------------------------------------------------
//     User information stored in database
//   - -------------------------------------------------------------------------------------------------------------------
type StoredUser struct {
	ID    uuid.UUID
	Name  string
	Email string
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Response from financials service
//   - -------------------------------------------------------------------------------------------------------------------
type UserFinancialsResponse struct {
	Status general.ResponseStatus
	Error  FinancialsError
	Result *UserFinancials
}

type UserFinancials struct {
	Balance      uint64
	Transactions []UserFinancialsTransaction
}

type UserFinancialsTransaction struct {
	CreatedDT time.Time
	Amount    uint64
	Reason    string
}

//   - -------------------------------------------------------------------------------------------------------------------
//     Error status for financials service
//   - -------------------------------------------------------------------------------------------------------------------
type FinancialsError string

const (
	FinancialsUserNotFound  FinancialsError = "user_not_found"
	FinancialsInternalerror FinancialsError = "internal_error"
)
