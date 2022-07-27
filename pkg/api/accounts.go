package api

import (
	"fmt"
	"time"
)

type Account struct {
	Id             int64     `json:"id"`
	Name           string    `json:"name"`
	AccountType    string    `json:"account_type"`
	AccountSubType string    `json:"account_sub_type"`
	OrgId          int64     `json:"org_id"`
	CurrentBalance float64   `json:"current_balance"`
	IsDeleted      bool      `json:"is_deleted"`
	CreatedTs      time.Time `json:"created_ts"`
	UpdatedTs      time.Time `json:"updated_ts"`
}

type AccountCreationRequest struct {
	Name           *string  `json:"name"`
	AccountType    *string  `json:"account_type"`
	AccountSubType *string  `json:"account_sub_type"`
	OrgId          *int64   `json:"org_id"`
	CurrentBalance *float64 `json:"current_balance"`
}

type AccountUpdateRequest struct {
	Name           *string  `json:"name"`
	AccountType    *string  `json:"account_type"`
	AccountSubType *string  `json:"account_sub_type"`
	OrgId          *int64   `json:"org_id"`
	CurrentBalance *float64 `json:"current_balance"`
}

type AccountsService struct {
	dao AccountsDataAccessor
}

func NewAccountsService(dao AccountsDataAccessor) *AccountsService {
	return &AccountsService{dao}
}

func (service *AccountsService) GetAccount(id int64) (Account, error) {
	account, err := service.dao.GetAccount(id)
	if err != nil {
		return account, fmt.Errorf("failed to retrieve account of id %d with err:%v", id, err)
	}
	return account, nil
}

func (service *AccountsService) CreateAccount(request AccountCreationRequest) error {
	_, err := service.dao.InsertAccount(request)
	if err != nil {
		return err
	}
	return nil
}

func (service *AccountsService) DeleteAccount(id int64) error {
	err := service.dao.DeleteAccount(id)
	if err != nil {
		return err
	}
	return nil
}

func (service *AccountsService) UpdateAccount(id int64, request AccountUpdateRequest) error {
	err := service.dao.UpdateAccount(id, request)
	if err != nil {
		return err
	}
	return nil
}

type AccountsDataAccessor interface {
	GetAccount(id int64) (Account, error)
	//ListAccounts(id int64) []Account
	InsertAccount(request AccountCreationRequest) (int64, error)
	DeleteAccount(id int64) error
	UpdateAccount(id int64, request AccountUpdateRequest) error
}
