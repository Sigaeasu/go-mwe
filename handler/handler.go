package handler

import (
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/Sigaeasu/go-mwe/models"
	"github.com/Sigaeasu/go-mwe/service"
	"github.com/Sigaeasu/go-mwe/repository"
	"github.com/Sigaeasu/go-mwe/utils/response"
	"github.com/golang-jwt/jwt/v4"
)

const (
	Limit = 10
)

type miniWalletHandler struct {
	miniWalletRepo repository.MiniWalletRepoInterface
}

type MiniWalletHandlerInterface interface {
	AuthMiniWallet(w http.ResponseWriter, r *http.Request)
	ViewMiniWalletBalance(w http.ResponseWriter, r *http.Request)
	ViewTransactions(w http.ResponseWriter, r *http.Request)
	EnableMiniWallet(w http.ResponseWriter, r *http.Request)
	DisableMiniWallet(w http.ResponseWriter, r *http.Request)
	DepositToMiniWallet(w http.ResponseWriter, r *http.Request)
	WithdrawFromMiniWallet(w http.ResponseWriter, r *http.Request)
}

func MiniWalletHandler(miniWalletRepo repository.MiniWalletRepoInterface) MiniWalletHandlerInterface {
	return &miniWalletHandler{
		miniWalletRepo: miniWalletRepo,
	}
}

func (h *miniWalletHandler) AuthMiniWallet(w http.ResponseWriter, r *http.Request) {
	customerXId := r.FormValue("customer_xid")
	if customerXId == "" {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: "Missing data for required field.",
			},
		}, http.StatusBadRequest)
		return
	}
	token, err := service.GenerateToken(customerXId)
	if err != nil {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}
	apiResponse(w, response.ResponseAPI{
		Status: "success",
		Data: map[string]string{
			"token": token,
		},
	}, http.StatusOK)
}

func (h *miniWalletHandler) ViewMiniWalletBalance(w http.ResponseWriter, r *http.Request) {
	cus := r.Context().Value(service.Customer).(jwt.MapClaims)
	custXId := cus["customer_xid"].(string)

	wallet, err := h.miniWalletRepo.FetchMiniWalletByID(custXId)
	if err != nil {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}
	if wallet.ID == "" {
		customerUnregistered(w)
		return
	}
	if !wallet.IsEnabled {
		walletIsDisabled(w)
		return
	}
	apiResponse(w, response.ResponseAPI{
		Status: "success",
		Data: ResponseWallet{
			ID:        wallet.ID,
			OwnedBy:   wallet.OwnedBy,
			Status:    "enabled",
			EnabledAt: wallet.EnabledAt.String(),
			Balance:   wallet.Balance,
		},
	}, http.StatusOK)
}

func (h *miniWalletHandler) ViewTransactions(w http.ResponseWriter, r *http.Request) {
	cus := r.Context().Value(service.Customer).(jwt.MapClaims)
	custXId := cus["customer_xid"].(string)

	wallet, err := h.miniWalletRepo.FetchMiniWalletByID(custXId)
	if err != nil {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}
	if wallet.ID == "" {
		customerUnregistered(w)
		return
	}
	if !wallet.IsEnabled {
		walletIsDisabled(w)
		return
	}

	transaction, err := h.miniWalletRepo.FetchTransactionByID(wallet.OwnedBy)
	if err != nil {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	apiResponse(w, response.ResponseAPI{
		Status: "success",
		Data: transaction,
	}, http.StatusOK)
}

func (h *miniWalletHandler) EnableMiniWallet(w http.ResponseWriter, r *http.Request) {
	cus := r.Context().Value(service.Customer).(jwt.MapClaims)
	custXId := cus["customer_xid"].(string)

	res, err := h.miniWalletRepo.ChangeStatusOnMiniWallet(custXId, true)
	if err != nil {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	status := "enabled"
	if !res.IsEnabled {
		status = "disabled"
	}
	apiResponse(w, response.ResponseAPI{
		Status: "success",
		Data: ResponseWallet{
			ID: res.ID,
			OwnedBy: res.OwnedBy,
			Status: status,
			EnabledAt: res.EnabledAt.String(),
			Balance: res.Balance,
		},
	}, http.StatusOK)
}

func (h *miniWalletHandler) DisableMiniWallet(w http.ResponseWriter, r *http.Request) {
	cus := r.Context().Value(service.Customer).(jwt.MapClaims)
	custXId := cus["customer_xid"].(string)

	res, err := h.miniWalletRepo.ChangeStatusOnMiniWallet(custXId, false)
	if err != nil {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	status := "enabled"
	if !res.IsEnabled {
		status = "disabled"
	}
	apiResponse(w, response.ResponseAPI{
		Status: "success",
		Data: ResponseWallet{
			ID: res.ID,
			OwnedBy: res.OwnedBy,
			Status: status,
			EnabledAt: res.EnabledAt.String(),
			Balance: res.Balance,
		},
	}, http.StatusOK)
}

func (h *miniWalletHandler) DepositToMiniWallet(w http.ResponseWriter, r *http.Request) {
	cus := r.Context().Value(service.Customer).(jwt.MapClaims)
	custXId := cus["customer_xid"].(string)

	rawAmount := r.FormValue("amount")
	amount, err := strconv.ParseFloat(rawAmount, 64)
	referenceId := r.FormValue("reference_id")

	reference, err := h.miniWalletRepo.CheckReferenceID(referenceId)
	if err != nil {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: "Duplicate Reference ID",
			},
		}, http.StatusInternalServerError)
		if reference != nil {
			return
		}
		return
	}

	wallet, err := h.miniWalletRepo.FetchMiniWalletByID(custXId)
	if err != nil {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}
	if wallet.ID == "" {
		customerUnregistered(w)
		return
	}
	if !wallet.IsEnabled {
		walletIsDisabled(w)
		return
	}
	
	params := models.ParamsWallet{
		Amount: amount,
		Balance: wallet.Balance,
		ReferenceID: referenceId,
		CreatedBy: wallet.OwnedBy,
	}
	deposit, err := h.miniWalletRepo.Deposit(params)
	if err != nil {
		h.miniWalletRepo.FailedTransaction(params, "deposit")
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}
	
	transaction, err := h.miniWalletRepo.NewTransaction(params, "deposit")
	if err != nil {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	apiResponse(w, response.ResponseAPI{
		Status: "success",
		Data: ResponseDepositWallet{
			ID: deposit.ID,
			DepositedBy: custXId,
			Status: "success",
			DepositAt: transaction.CreatedAt.String(),
			Amount: amount,
			ReferenceId: referenceId,
		},
	}, http.StatusOK)
}

func (h *miniWalletHandler) WithdrawFromMiniWallet(w http.ResponseWriter, r *http.Request) {
	cus := r.Context().Value(service.Customer).(jwt.MapClaims)
	custXId := cus["customer_xid"].(string)

	rawAmount := r.FormValue("amount")
	amount, err := strconv.ParseFloat(rawAmount, 64)
	referenceId := r.FormValue("reference_id")

	wallet, err := h.miniWalletRepo.FetchMiniWalletByID(custXId)
	if err != nil {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}
	if wallet.ID == "" {
		customerUnregistered(w)
		return
	}
	if !wallet.IsEnabled {
		walletIsDisabled(w)
		return
	}

	if wallet.Balance < amount {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: "Balance is insufficient",
			},
		}, http.StatusInternalServerError)
		return
	}

	if !wallet.IsEnabled {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: "Wallet disabled",
			},
		}, http.StatusInternalServerError)
		return
	}
	
	params := models.ParamsWallet{
		Amount: amount,
		Balance: wallet.Balance,
		ReferenceID: referenceId,
		CreatedBy: wallet.OwnedBy,
	}
	withdraw, err := h.miniWalletRepo.Withdraw(params)
	if err != nil {
		h.miniWalletRepo.FailedTransaction(params, "withdrawn")
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}
	
	transaction, err := h.miniWalletRepo.NewTransaction(params, "withdraw")
	if err != nil {
		apiResponse(w, response.ResponseAPI{
			Status: "fail",
			Data: &response.ApiError{
				Error: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	apiResponse(w, response.ResponseAPI{
		Status: "success",
		Data: ResponseWithdrawWallet{
			ID: withdraw.ID,
			WithdrawnBy: custXId,
			Status: "success",
			WithdrawnAt: transaction.CreatedAt.String(),
			Amount: amount,
			ReferenceId: referenceId,
		},
	}, http.StatusOK)
}

func apiResponse(ar http.ResponseWriter, data interface{}, statusCode int) {
	ar.Header().Set("Content-type", "application/json")
	ar.WriteHeader(statusCode)
	json.NewEncoder(ar).Encode(data)
}

func customerUnregistered(w http.ResponseWriter) {
	apiResponse(w, response.ResponseAPI{
		Status: "fail",
		Data: &response.ApiError{
			Error: "Customer is not registered",
		},
	}, http.StatusInternalServerError)
}

func walletIsDisabled(w http.ResponseWriter) {
	apiResponse(w, response.ResponseAPI{
		Status: "fail",
		Data: &response.ApiError{
			Error: "Wallet disabled",
		},
	}, http.StatusInternalServerError)
}