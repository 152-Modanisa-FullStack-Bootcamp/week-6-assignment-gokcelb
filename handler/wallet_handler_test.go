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
	mockWallets := map[string]*model.Wallet{
		"lacin": {
			Username: "lacin",
			Balance:  0,
		},
		"yuksel": {
			Username: "yuksel",
			Balance:  200,
		},
	}
	mockService.EXPECT().GetAll().Return(mockWallets).Times(1)

	h := handler.NewWallet(mockService)
	h.HandleWalletEndpoints(resW, req)

	assert.Equal(t, http.StatusOK, resW.Result().StatusCode)

}

func TestHandleWalletEndpoints_RedirectsToGet(t *testing.T) {
	testCases := []struct {
		desc               string
		target             string
		username           string
		expectedWallet     *model.Wallet
		expectedErr        error
		expectedStatusCode int
	}{
		{
			desc:               "given target /lacin, service get is called with lacin and returns nil and not exists error",
			target:             "/lacin",
			username:           "lacin",
			expectedWallet:     nil,
			expectedErr:        service.ErrWalletNotExists,
			expectedStatusCode: http.StatusNotFound,
		},
		{
			desc:     "given target /doga, service get is called with doga and returns doga's wallet and nil",
			target:   "/doga",
			username: "doga",
			expectedWallet: &model.Wallet{
				Username: "doga",
				Balance:  0,
			},
			expectedErr:        nil,
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			resW := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, tC.target, http.NoBody)

			mockService := createMockWalletService(t)
			mockService.EXPECT().Get(tC.username).Return(tC.expectedWallet, tC.expectedErr)

			h := handler.NewWallet(mockService)
			h.HandleWalletEndpoints(resW, req)

			assert.Equal(t, tC.expectedStatusCode, resW.Result().StatusCode)
		})
	}
}

func TestHandleWalletEndpoints_RedirectsToCreate(t *testing.T) {
	testCases := []struct {
		desc               string
		httpMethod         string
		target             string
		username           string
		expectedStatusCode int
	}{
		{
			desc:               "given method PUT and target /fatma, calls service create with fatma and returns fatma's wallet",
			httpMethod:         http.MethodPut,
			target:             "/fatma",
			username:           "fatma",
			expectedStatusCode: http.StatusCreated,
		},
		{
			desc:               "given method PUT and targete /lacin/, calls service create with lacin and returns lacin's wallet",
			httpMethod:         http.MethodPut,
			target:             "/lacin/",
			username:           "lacin",
			expectedStatusCode: http.StatusCreated,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			resW := httptest.NewRecorder()
			req := httptest.NewRequest(tC.httpMethod, tC.target, http.NoBody)

			mockService := createMockWalletService(t)
			mockService.EXPECT().Create(tC.username)

			h := handler.NewWallet(mockService)
			h.HandleWalletEndpoints(resW, req)

			assert.Equal(t, tC.expectedStatusCode, resW.Result().StatusCode)

		})
	}
}

type updateRequestBody struct {
	Balance int `json:"balance"`
}

type errorResponseBody struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
		updatedWallet := &model.Wallet{
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
		mockService.EXPECT().Update("lacin", 100).Return(nil, service.ErrWalletNotExists)

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
		mockService.EXPECT().Update("lacin", -500).Return(nil, service.ErrBalanceBelowLimit)

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

func createMockWalletService(t *testing.T) *mocks.MockWalletService {
	return mocks.NewMockWalletService(gomock.NewController(t))
}
