package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/handler"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/mocks"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createMockWalletService(t *testing.T) *mocks.MockWalletService {
	return mocks.NewMockWalletService(gomock.NewController(t))
}

type updateRequestBody struct {
	Balance int `json:"balance"`
}

type errorResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func TestHandleWalletEndpoints_InvalidEndpoints(t *testing.T) {
	testCases := []struct {
		desc               string
		httpMethod         string
		target             string
		expectedStatusCode int
	}{
		{
			desc:               "when http method is get and target is /asd/fa, expect status not found",
			httpMethod:         http.MethodGet,
			target:             "/asd/fa",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			desc:               "when http method is put and target is /, expect status not found",
			httpMethod:         http.MethodPut,
			target:             "/",
			expectedStatusCode: http.StatusNotFound,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			resW := httptest.NewRecorder()
			req := httptest.NewRequest(tC.httpMethod, tC.target, http.NoBody)

			h := handler.NewWallet(nil)
			h.HandleWalletEndpoints(resW, req)

			assert.Equal(t, tC.expectedStatusCode, resW.Result().StatusCode)
		})
	}
}

func TestHandleWalletEndpoints_RedirectsToGetAll(t *testing.T) {
	resW := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

	mockService := createMockWalletService(t)
	mockWallets := []model.Wallet{
		{Username: "lacin", Balance: 0},
		{Username: "yuksel", Balance: 200},
	}
	mockService.EXPECT().GetAll().Return(mockWallets).Times(1)

	h := handler.NewWallet(mockService)
	h.HandleWalletEndpoints(resW, req)

	resBody, _ := json.Marshal(mockWallets)

	assert.Equal(t, http.StatusOK, resW.Result().StatusCode)
	assert.Equal(t, string(resBody), resW.Body.String())
}

func TestHandleWalletEndpoints_RedirectsToGet(t *testing.T) {
	t.Run("given target /lacin, service get is called with lacin and returns empty object and not exists error", func(t *testing.T) {
		resW := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/lacin", http.NoBody)

		mockService := createMockWalletService(t)
		mockService.EXPECT().Get("lacin").Return(model.Wallet{}, service.ErrWalletNotExists)

		h := handler.NewWallet(mockService)
		h.HandleWalletEndpoints(resW, req)

		resBody, _ := json.Marshal(errorResponseBody{
			Code:    http.StatusNotFound,
			Message: service.ErrWalletNotExists.Error(),
		})

		assert.Equal(t, http.StatusNotFound, resW.Result().StatusCode)
		assert.Equal(t, string(resBody), resW.Body.String())
	})

	t.Run("given target /doga, service get is called with doga and returns doga's wallet and nil", func(t *testing.T) {
		resW := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/doga", http.NoBody)

		mockService := createMockWalletService(t)
		mockWallet := model.Wallet{Username: "doga", Balance: 50}
		mockService.EXPECT().Get("doga").Return(mockWallet, nil)

		h := handler.NewWallet(mockService)
		h.HandleWalletEndpoints(resW, req)

		resBody, _ := json.Marshal(mockWallet)

		assert.Equal(t, http.StatusOK, resW.Result().StatusCode)
		assert.Equal(t, string(resBody), resW.Body.String())
	})
}

func TestHandleWalletEndpoints_RedirectsToCreate(t *testing.T) {
	t.Run("given method PUT and target /fatma, calls service create with fatma and returns fatma's wallet", func(t *testing.T) {
		resW := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/fatma", http.NoBody)

		mockService := createMockWalletService(t)
		mockWallet := model.Wallet{Username: "fatma", Balance: 0}
		mockService.EXPECT().Create("fatma").Return(mockWallet)

		h := handler.NewWallet(mockService)
		h.HandleWalletEndpoints(resW, req)

		resBody, _ := json.Marshal(mockWallet)

		assert.Equal(t, http.StatusCreated, resW.Result().StatusCode)
		assert.Equal(t, string(resBody), resW.Body.String())
	})

	t.Run("given method PUT and targete /lacin/, calls service create with lacin and returns lacin's wallet", func(t *testing.T) {
		resW := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/lacin/", http.NoBody)

		mockService := createMockWalletService(t)
		mockWallet := model.Wallet{Username: "lacin", Balance: 0}
		mockService.EXPECT().Create("lacin").Return(mockWallet)

		h := handler.NewWallet(mockService)
		h.HandleWalletEndpoints(resW, req)

		resBody, _ := json.Marshal(mockWallet)

		assert.Equal(t, http.StatusCreated, resW.Result().StatusCode)
		assert.Equal(t, string(resBody), resW.Body.String())
	})
}

func TestHandleWalletEndpoints_RedirectsToUpdate(t *testing.T) {
	t.Run("successfully update wallet", func(t *testing.T) {
		updateReqBody := &updateRequestBody{
			Balance: 100,
		}
		body, _ := json.Marshal(updateReqBody)

		resW := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/lacin", bytes.NewReader(body))

		mockService := createMockWalletService(t)
		updatedWallet := model.Wallet{
			Username: "lacin",
			Balance:  100,
		}
		mockService.EXPECT().Update("lacin", 100).Return(updatedWallet, nil)

		h := handler.NewWallet(mockService)
		h.HandleWalletEndpoints(resW, req)

		resBody, _ := json.Marshal(updatedWallet)

		assert.Equal(t, http.StatusCreated, resW.Result().StatusCode)
		assert.Equal(t, string(resBody), resW.Body.String())
	})

	t.Run("given user not exists error, unsuccessful update", func(t *testing.T) {
		updateReqBody := &updateRequestBody{
			Balance: 100,
		}
		body, _ := json.Marshal(updateReqBody)

		resW := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/lacin", bytes.NewReader(body))

		mockService := createMockWalletService(t)
		mockService.EXPECT().Update("lacin", 100).Return(model.Wallet{}, service.ErrWalletNotExists)

		h := handler.NewWallet(mockService)
		h.HandleWalletEndpoints(resW, req)

		errorResponse := &errorResponseBody{
			Code:    http.StatusNotFound,
			Message: service.ErrWalletNotExists.Error(),
		}
		resBody, _ := json.Marshal(errorResponse)

		assert.Equal(t, http.StatusNotFound, resW.Result().StatusCode)
		assert.Equal(t, string(resBody), resW.Body.String())
	})

	t.Run("given updated balance is below limit, unsuccessful update", func(t *testing.T) {
		updateReqBody := &updateRequestBody{
			Balance: -500,
		}
		body, _ := json.Marshal(updateReqBody)

		resW := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/lacin", bytes.NewReader(body))

		mockService := createMockWalletService(t)
		mockService.EXPECT().Update("lacin", -500).Return(model.Wallet{}, service.ErrBalanceBelowLimit)

		h := handler.NewWallet(mockService)
		h.HandleWalletEndpoints(resW, req)

		errorResponse := &errorResponseBody{
			Code:    http.StatusBadRequest,
			Message: service.ErrBalanceBelowLimit.Error(),
		}
		resBody, _ := json.Marshal(errorResponse)

		assert.Equal(t, http.StatusBadRequest, resW.Result().StatusCode)
		assert.Equal(t, string(resBody), resW.Body.String())
	})
}
