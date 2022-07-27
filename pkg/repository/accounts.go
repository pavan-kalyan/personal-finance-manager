package repository

import (
	"database/sql"
	"fmt"
	"go-personal-finance/pkg/api"
)

const accountsTableName = "accounts"
const organizationsTableName = "organizations"

type accountsDAO struct {
	db *sql.DB
}

func NewAccountsDAO(db *sql.DB) *accountsDAO {
	return &accountsDAO{db}
}

func (dao *accountsDAO) GetAccount(id int64) (api.Account, error) {
	row := dao.db.QueryRow(fmt.Sprintf("SELECT id, name, type, subtype, org_id, current_balance, is_deleted, "+
		"created_ts, updated_ts FROM %s WHERE id = ? AND is_deleted = 0", accountsTableName), id)
	account := api.Account{}
	//TODO: too much dependence on order of columns. Find a better way.
	var accountSubType sql.NullString
	err := row.Scan(&account.Id, &account.Name, &account.AccountType, &accountSubType, &account.OrgId,
		&account.CurrentBalance, &account.IsDeleted, &account.CreatedTs, &account.UpdatedTs)
	if err != nil {
		return account, fmt.Errorf("failed to retrieve account of id %d with err: %v", id, err)
	}
	if accountSubType.Valid {
		account.AccountSubType = accountSubType.String
	}

	return account, nil
}
func (dao *accountsDAO) InsertAccount(request api.AccountCreationRequest) (int64, error) {

	var accountName = NewNullString(request.Name)
	var accountType = NewNullString(request.AccountType)
	var accountSubType = NewNullString(request.AccountSubType)
	var orgId = *request.OrgId
	var currentBalance float64
	if request.CurrentBalance != nil {
		currentBalance = *request.CurrentBalance
	}
	result, err := dao.db.Exec(fmt.Sprintf("INSERT INTO %s (name, type, subtype, org_id, current_balance) "+
		"VALUES(?,?,?,?,?)",
		accountsTableName), accountName, accountType, accountSubType, orgId, currentBalance)
	if err != nil {
		return -1, fmt.Errorf("failed to insert new account due to error %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("failed to retrieve last inserted id %v", err)
	}
	return id, err
}

func (dao *accountsDAO) DeleteAccount(id int64) error {
	result, err := dao.db.Exec(fmt.Sprintf("UPDATE %s SET is_deleted=1 WHERE id = ? AND is_deleted = 0", accountsTableName), id)
	if err != nil {
		return fmt.Errorf("failed to delete account %d due to error %v", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to verify if account %d has been deleted due to error %v", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("WARN: detected request to delete non-existent account with id %d \n", id)
	}
	return nil
}

// UpdateAccount assumes all values in request struct are non nil/**
func (dao *accountsDAO) UpdateAccount(id int64, request api.AccountUpdateRequest) error {
	var accountName = NewNullString(request.Name)
	var accountType = NewNullString(request.AccountType)
	var accountSubType = NewNullString(request.AccountSubType)
	var orgId = *request.OrgId
	var currentBalance float64
	if request.CurrentBalance != nil {
		currentBalance = *request.CurrentBalance
	}
	result, err := dao.db.Exec(fmt.Sprintf("UPDATE %s SET name = ?, type = ?, subtype = ?, org_id = ?, "+
		"current_balance = ?  "+
		"WHERE id = ? AND is_deleted = 0", accountsTableName), accountName, accountType, accountSubType, orgId, currentBalance, id)
	if err != nil {
		return fmt.Errorf("failed to update account %d due to error %v", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to verify if account %d has been updated due to error %v", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("WARN: detected request to update non-existent account with id %d \n", id)
	}
	return nil
}
