package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"go-personal-finance/pkg/api"
	"testing"
	"time"
)

func TestAccountsDAO_GetAccount(t *testing.T) {

	db, mock := setupMockDB(t)
	defer db.Close()
	dao := accountsDAO{db: db}

	expectedSelectQuery := "SELECT id, name, type, subtype, org_id, current_balance, is_deleted, created_ts, " +
		"updated_ts FROM accounts WHERE id = ? AND is_deleted = 0"

	t.Run("testing happy flow", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"id", "name", "type", "subtype", "org_id", "current_balance", "is_deleted",
			"created_ts", "updated_ts"}).AddRow(1, "acc1", "investment", "life insurance", int64(1), 30.00, false, time.Now(), time.Now())
		mock.ExpectQuery(expectedSelectQuery).WillReturnRows(rows)

		account, err := dao.GetAccount(1)
		if err != nil {
			t.Errorf("Unexpected error when retrieving account: %v", err)
			return
		}

		actualId := account.Id
		expectedId := int64(1)
		if actualId != expectedId {
			t.Errorf("unexpected id returned")
		}
		checkingMockExpectations(t, mock)
	})

	t.Run("testing situation where we fail to retrieve account", func(t *testing.T) {

		mock.ExpectQuery(expectedSelectQuery).WillReturnError(errors.New("db timeout"))

		_, err := dao.GetAccount(1)
		expectedErrorMsg := "failed to retrieve account of id 1 with err: db timeout"

		if expectedErrorMsg != err.Error() {
			t.Errorf("expected error message %s but found %s", expectedErrorMsg, err.Error())
		}
		checkingMockExpectations(t, mock)
	})
}

func TestAccountsDAO_InsertAccount(t *testing.T) {

	db, mock := setupMockDB(t)
	defer db.Close()
	dao := accountsDAO{db: db}

	accountName := "acc1"
	accountType := "savings"
	orgId := int64(1)
	currentBalance := 24.00
	createQuery := "INSERT INTO accounts (name, type, subtype, org_id, current_balance) VALUES(?,?,?,?,?)"
	t.Run("testing happy flow", func(t *testing.T) {
		mock.ExpectExec(createQuery).
			WithArgs(accountName, accountType, nil, orgId, currentBalance).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		actualId, err := dao.InsertAccount(api.AccountCreationRequest{Name: &accountName, AccountSubType: nil,
			AccountType: &accountType, OrgId: &orgId, CurrentBalance: &currentBalance})
		if err != nil {
			t.Errorf("Unexpected error when trying to insert account: %v", err)
			return
		}

		expectedId := int64(1)
		if actualId != expectedId {
			t.Errorf("unexpected id returned")
		}

		checkingMockExpectations(t, mock)
	})

	t.Run("testing general error flow", func(t *testing.T) {
		mock.ExpectExec(createQuery).
			WithArgs(accountName, accountType, nil, orgId, currentBalance).
			WillReturnError(errors.New("db timeout"))

		_, err := dao.InsertAccount(api.AccountCreationRequest{Name: &accountName, AccountSubType: nil,
			AccountType: &accountType, OrgId: &orgId, CurrentBalance: &currentBalance})
		expectedErrorMsg := "failed to insert new account due to error db timeout"

		if expectedErrorMsg != err.Error() {
			t.Errorf("expected error message %s but found %s", expectedErrorMsg, err.Error())
		}
		checkingMockExpectations(t, mock)
	})
}

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error occurred when starting mock sql")
	}
	return db, mock
}

func TestAccountsDAO_DeleteAccount(t *testing.T) {

	db, mock := setupMockDB(t)
	defer db.Close()
	dao := accountsDAO{db: db}

	t.Run("testing happy flow", func(t *testing.T) {
		accountId := int64(1)
		mock.ExpectExec("UPDATE accounts SET is_deleted=1 WHERE id = ? AND is_deleted = 0").
			WithArgs(accountId).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := dao.DeleteAccount(accountId)
		if err != nil {
			t.Errorf("Unexpected error when trying to delete account: %v", err)
			return
		}
	})

	t.Run("testing general error flow", func(t *testing.T) {
		accountId := int64(1)
		mock.ExpectExec("UPDATE accounts SET is_deleted=1 WHERE id = ? AND is_deleted = 0").
			WithArgs(accountId).
			WillReturnError(errors.New("db timeout"))

		err := dao.DeleteAccount(accountId)
		expectedErrorMsg := "failed to delete account 1 due to error db timeout"
		if expectedErrorMsg != err.Error() {
			t.Errorf("expected error message %s but found %s", expectedErrorMsg, err.Error())
		}
	})

	t.Run("testing deletion of non existent account", func(t *testing.T) {
		accountId := int64(1)
		mock.ExpectExec("UPDATE accounts SET is_deleted=1 WHERE id = ? AND is_deleted = 0").
			WithArgs(accountId).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := dao.DeleteAccount(accountId)
		expectedErrorMsg := "WARN: detected request to delete non-existent account with id 1 \n"
		if expectedErrorMsg != err.Error() {
			t.Errorf("expected error message %s but found %s", expectedErrorMsg, err.Error())
		}
	})

	checkingMockExpectations(t, mock)
}

func TestAccountsDAO_UpdateAccount(t *testing.T) {

	db, mock := setupMockDB(t)
	defer db.Close()

	accountId := int64(1)
	accountName := "acc1"
	accountType := "savings"
	orgId := int64(1)
	currentBalance := 24.00

	expectedUpdateQuery := "UPDATE accounts SET name = ?, type = ?, subtype = ?, org_id = ?, current_balance = ? " +
		"WHERE id = ? AND is_deleted = 0"
	dao := accountsDAO{db: db}

	t.Run("testing happy flow of updating account", func(t *testing.T) {
		mock.ExpectExec(expectedUpdateQuery).
			WithArgs(accountName, accountType, nil, orgId, currentBalance, accountId).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := dao.UpdateAccount(accountId, api.AccountUpdateRequest{AccountType: &accountType, Name: &accountName,
			AccountSubType: nil, OrgId: &orgId, CurrentBalance: &currentBalance})
		if err != nil {
			t.Errorf("Unexpected error when trying to update account: %v", err)
			return
		}

		checkingMockExpectations(t, mock)
	})

	t.Run("testing generic error handling flows of updating an account", func(t *testing.T) {
		mock.ExpectExec(expectedUpdateQuery).
			WithArgs(accountName, accountType, nil, orgId, currentBalance, accountId).
			WillReturnError(errors.New("db timeout"))

		err := dao.UpdateAccount(accountId, api.AccountUpdateRequest{AccountType: &accountType, Name: &accountName,
			AccountSubType: nil, OrgId: &orgId, CurrentBalance: &currentBalance})
		expectedErrorMsg := "failed to update account 1 due to error db timeout"

		if expectedErrorMsg != err.Error() {
			t.Errorf("expected error message %s but found %s", expectedErrorMsg, err.Error())
		}

		checkingMockExpectations(t, mock)
	})

	t.Run("testing situation when account doesn't exist with that id", func(t *testing.T) {
		mock.ExpectExec(expectedUpdateQuery).
			WithArgs(accountName, accountType, nil, orgId, currentBalance, accountId).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := dao.UpdateAccount(accountId, api.AccountUpdateRequest{AccountType: &accountType, Name: &accountName,
			AccountSubType: nil, OrgId: &orgId, CurrentBalance: &currentBalance})
		expectedErrorMsg := "WARN: detected request to update non-existent account with id 1 \n"

		if expectedErrorMsg != err.Error() {
			t.Errorf("expected error message %s but found %s", expectedErrorMsg, err.Error())
		}

		checkingMockExpectations(t, mock)
	})

	t.Run("testing situation when account hasn't changed due to some anomaly", func(t *testing.T) {
		mock.ExpectExec(expectedUpdateQuery).
			WithArgs(accountName, accountType, nil, orgId, currentBalance, accountId).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := dao.UpdateAccount(accountId, api.AccountUpdateRequest{AccountType: &accountType, Name: &accountName,
			AccountSubType: nil, OrgId: &orgId, CurrentBalance: &currentBalance})
		expectedErrorMsg := "WARN: detected request to update non-existent account with id 1 \n"

		if expectedErrorMsg != err.Error() {
			t.Errorf("expected error message %s but found %s", expectedErrorMsg, err.Error())
		}

		checkingMockExpectations(t, mock)
	})
}

func checkingMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	err := mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("expected query not found with err: %v", err)
	}
}
