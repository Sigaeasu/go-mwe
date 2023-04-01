package repository

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"github.com/go-pg/pg/v10"
	"github.com/Sigaeasu/go-mwe/models/entity"
	"github.com/Sigaeasu/go-mwe/models"
)

type MiniWalletRepoInterface interface {
	FetchMiniWalletByID(customerXId string) (*entity.Wallet, error)	
	FetchTransactionByID(customerId string) ([]entity.Transaction, error)	
	ChangeStatusOnMiniWallet(customerXId string, status bool) (*entity.Wallet, error)
	Deposit(params models.ParamsWallet) (*entity.Wallet, error)	
	Withdraw(params models.ParamsWallet) (*entity.Wallet, error)
	NewTransaction(params models.ParamsWallet, transactionType string) (*entity.Transaction, error)
	FailedTransaction(params models.ParamsWallet, transactionType string)
	CheckReferenceID(referenceID string) (*entity.Transaction, error)
}

type miniWalletDatabase struct {
	mutex  sync.RWMutex
	dbConn *pg.DB
}

func MiniWalletRepository(c *pg.DB) MiniWalletRepoInterface {
	return &miniWalletDatabase{dbConn: c}
}

func (pdb *miniWalletDatabase) FetchMiniWalletByID(customerXId string) (*entity.Wallet, error) {
	var wallet entity.Wallet
	pdb.mutex.Lock()
	err := pdb.dbConn.Model(&wallet).
		Where("id = ?", customerXId).
		Select()
	if err != nil {
		if err != pg.ErrNoRows {
			return nil, err
		} 
	}
	pdb.mutex.Unlock()
	return &wallet, nil
}

func (pdb *miniWalletDatabase) FetchTransactionByID(customerId string) ([]entity.Transaction, error) {
	var transaction []entity.Transaction
	pdb.mutex.Lock()
	err := pdb.dbConn.Model(&transaction).
		Where("created_by = ?", customerId).
		Select()
	if err != nil {
		if err != pg.ErrNoRows {
			return nil, err
		} 
	}
	pdb.mutex.Unlock()
	return transaction, nil
}

func (pdb *miniWalletDatabase) ChangeStatusOnMiniWallet(customerXId string, status bool) (*entity.Wallet, error) {
	if customerXId == "" {
		return nil, errors.New("customer_xid is empty")
	}

	resWallet, err := pdb.FetchMiniWalletByID(customerXId)
	if err != nil {
		return nil, err
	}

	if resWallet != nil {
		if status && resWallet.IsEnabled {
			return nil, fmt.Errorf("Already enabled")
		} else if !status && resWallet.IsEnabled == false {
			return nil, fmt.Errorf("Already disabled")
		} else {			
			wallet := entity.Wallet{}
			pdb.mutex.Lock()
			if resWallet.IsEnabled {
				res, err := pdb.dbConn.Model(&wallet).
					Where("id = ?", customerXId).
					Set("is_enabled = ?", status).
					Set("enabled_at = ?", time.Now()).
					Update()
				if err != nil {
					return nil, err
				}
				if res.RowsAffected() == 0 {
					return nil, fmt.Errorf("Fails to enable wallet")
				}
			} else {
				res, err := pdb.dbConn.Model(&wallet).
					Where("id = ?", customerXId).
					Set("is_enabled = ?", status).
					Set("disabled_at = ?", time.Now()).
					Update()
				if err != nil {
					return nil, err
				}
				if res.RowsAffected() == 0 {
					return nil, fmt.Errorf("Fails to disable wallet")
				}
			}			
			pdb.mutex.Unlock()
			return &wallet, nil
		}
	}
	return nil, errors.New("Customer not found")
}

func (pdb *miniWalletDatabase) Deposit(params models.ParamsWallet) (*entity.Wallet, error) {
	var wallet entity.Wallet
	pdb.mutex.Lock()
	res, err := pdb.dbConn.Model(&wallet).
		Where("owned_by = ? ", params.CreatedBy).
		Set("balance = ?", params.Balance+params.Amount).
		Update()
	if err != nil {
		return nil, err
	}
	if res.RowsAffected() == 0 {
		return nil, fmt.Errorf("Deposit Failed")
	}
	pdb.mutex.Unlock()
	return &wallet, nil
}

func (pdb *miniWalletDatabase) Withdraw(params models.ParamsWallet) (*entity.Wallet, error) {
	var wallet entity.Wallet
	pdb.mutex.Lock()
	res, err := pdb.dbConn.Model(&wallet).
		Where("id = ? ", params.CreatedBy).
		Set("balance = ?", params.Balance-params.Amount).
		Update()
	if err != nil {
		return nil, err
	}
	if res.RowsAffected() == 0 {
		return nil, fmt.Errorf("Withdraw Failed")
	}
	pdb.mutex.Unlock()
	return &wallet, nil
}

func (pdb *miniWalletDatabase) NewTransaction(params models.ParamsWallet, transactionType string) (*entity.Transaction, error) {
	var result entity.Transaction
	transaction := entity.Transaction{
		Amount: params.Amount,
		Type: transactionType,
		Status: "success",
		ReferenceID: params.ReferenceID,
		CreatedBy: params.CreatedBy,
		CreatedAt: time.Now(),
	}
	pdb.mutex.Lock()
	res, err := pdb.dbConn.Model(&transaction).Insert()
	if err != nil {
		return nil, err
	}
	if res.RowsAffected() == 0 {
		return nil, fmt.Errorf("Fail to create transaction log")
	}
	pdb.mutex.Unlock()
	return &result, nil
}

func (pdb *miniWalletDatabase) FailedTransaction(params models.ParamsWallet, transactionType string) {
	transaction := entity.Transaction{
		Amount: params.Amount,
		Type: transactionType,
		Status: "failed",
		ReferenceID: params.ReferenceID,
		CreatedBy: params.CreatedBy,
		CreatedAt: time.Now(),
	}
	pdb.mutex.Lock()
	pdb.dbConn.Model(&transaction).Insert()
	pdb.mutex.Unlock()
}

func (pdb *miniWalletDatabase) CheckReferenceID(referenceID string) (*entity.Transaction, error) {
	var result entity.Transaction
	pdb.mutex.Lock()
	err := pdb.dbConn.Model(&result).
		Where("reference_id = ?", referenceID).
		Select()
	pdb.mutex.Unlock()
	if err != nil && err == pg.ErrNoRows {
		return &result, nil
	}
	return nil, err
}