package api

import (
	"errors"
	"testing"
	"time"
)

type happyMock struct {
}

func (m *happyMock) GetAccount(id int64) (Account, error) {
	return Account{
		Id:             int64(1),
		Name:           "acc1",
		AccountType:    "savings",
		AccountSubType: "sub",
		OrgId:          int64(1),
		CurrentBalance: 50.00,
		IsDeleted:      false,
		CreatedTs:      time.Now(),
		UpdatedTs:      time.Now(),
	}, nil
}

func (m *happyMock) InsertAccount(request AccountCreationRequest) (int64, error) {
	return 1, nil
}
func (m *happyMock) DeleteAccount(id int64) error {
	return nil
}

func (m *happyMock) UpdateAccount(id int64, request AccountUpdateRequest) error {
	return nil
}

func TestAccountsService_GetAccount(t *testing.T) {

	t.Run("testing happy flow of getting an account", func(t *testing.T) {
		dao := happyMock{}
		accountsService := NewAccountsService(&dao)
		account, err := accountsService.GetAccount(1)
		if err != nil {
			t.Errorf("unexpected error when retrieving account")
		}
		if account.Id != 1 {
			t.Errorf("unexpected account id %d", account.Id)
		}
	})

	t.Run("testing general error flow of getting an account", func(t *testing.T) {

		dao := errorMock{}
		accountsService := NewAccountsService(&dao)
		_, err := accountsService.GetAccount(1)
		expectedErrorMsg := "failed to retrieve account of id 1 with err:db timeout"
		if err.Error() != expectedErrorMsg {
			t.Errorf("expected error message %s but found %s", expectedErrorMsg, err.Error())
		}
	})
}

func TestAccountsService_DeleteAccount(t *testing.T) {

	t.Run("testing happy flow of deleting an account", func(t *testing.T) {
		dao := happyMock{}
		accountsService := NewAccountsService(&dao)
		err := accountsService.DeleteAccount(1)
		if err != nil {
			t.Errorf("unexpected error when deleting account")
		}
	})

	t.Run("testing general error flow of deleting an account", func(t *testing.T) {
		dao := errorMock{}
		accountsService := NewAccountsService(&dao)
		err := accountsService.DeleteAccount(1)
		expectedErrorMsg := "db timeout"
		if err.Error() != expectedErrorMsg {
			t.Errorf("expected error message %s but found %s", expectedErrorMsg, err.Error())
		}
	})
}

func TestAccountsService_UpdateAccount(t *testing.T) {

	accountName := "acc1"
	accountType := "savings"
	orgId := int64(1)
	currentBalance := 24.00
	t.Run("testing happy flow of updating an account", func(t *testing.T) {
		dao := happyMock{}
		accountsService := NewAccountsService(&dao)

		err := accountsService.UpdateAccount(1, AccountUpdateRequest{AccountType: &accountType, Name: &accountName,
			AccountSubType: nil, OrgId: &orgId, CurrentBalance: &currentBalance})
		if err != nil {
			t.Errorf("unexpected error when updating account")
		}
	})

	t.Run("testing general error flow of updating an account", func(t *testing.T) {

		dao := errorMock{}
		accountsService := NewAccountsService(&dao)
		err := accountsService.UpdateAccount(1, AccountUpdateRequest{AccountType: &accountType, Name: &accountName,
			AccountSubType: nil, OrgId: &orgId, CurrentBalance: &currentBalance})
		expectedErrorMsg := "db timeout"
		if err.Error() != expectedErrorMsg {
			t.Errorf("expected error message %s but found %s", expectedErrorMsg, err.Error())
		}
	})
}

func TestAccountsService_CreateAccount(t *testing.T) {

	accountName := "acc1"
	accountType := "savings"
	orgId := int64(1)
	currentBalance := 24.00
	t.Run("testing happy flow of creating an account", func(t *testing.T) {

		dao := happyMock{}
		accountsService := NewAccountsService(&dao)
		err := accountsService.CreateAccount(AccountCreationRequest{AccountType: &accountType, Name: &accountName,
			AccountSubType: nil, OrgId: &orgId, CurrentBalance: &currentBalance})
		if err != nil {
			t.Errorf("unexpected error when creating account")
		}
	})

	t.Run("testing general error flow of creating an account", func(t *testing.T) {

		dao := errorMock{}
		accountsService := NewAccountsService(&dao)
		err := accountsService.CreateAccount(AccountCreationRequest{AccountType: &accountType, Name: &accountName,
			AccountSubType: nil, OrgId: &orgId, CurrentBalance: &currentBalance})
		expectedErrorMsg := "db timeout"
		if err.Error() != expectedErrorMsg {
			t.Errorf("expected error message %s but found %s", expectedErrorMsg, err.Error())
		}
	})
}

type errorMock struct {
}

func (m *errorMock) GetAccount(id int64) (Account, error) {
	return Account{}, errors.New("db timeout")
}

func (m *errorMock) InsertAccount(request AccountCreationRequest) (int64, error) {
	return 0, errors.New("db timeout")
}
func (m *errorMock) DeleteAccount(id int64) error {
	return errors.New("db timeout")
}

func (m *errorMock) UpdateAccount(id int64, request AccountUpdateRequest) error {
	return errors.New("db timeout")
}
