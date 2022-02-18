package service_test

import (
	"testing"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/config"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/mocks"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createMockRepository(t *testing.T) *mocks.MockWalletRepository {
	return mocks.NewMockWalletRepository(gomock.NewController(t))
}

func TestGetAll(t *testing.T) {

	mockRepo := createMockRepository(t)
	mockWallets := []model.Wallet{
		{Username: "yuksel", Balance: 350},
		{Username: "lacin", Balance: 0},
	}
	mockRepo.EXPECT().GetAll().Return(mockWallets)

	s := service.NewWallet(nil, mockRepo)
	wallets := s.GetAll()

	assert.Equal(t, mockWallets, wallets)
}

func TestGet_WalletDoesNotExist_ReturnsNilAndErrWalletNotExists(t *testing.T) {
	mockRepo := createMockRepository(t)
	mockRepo.EXPECT().Exists("doga").Return(false)

	s := service.NewWallet(nil, mockRepo)
	wallet, err := s.Get("doga")

	assert.Empty(t, wallet)
	assert.EqualError(t, err, service.ErrWalletNotExists.Error())
}

func TestGet_WalletExists_ReturnsWalletAndNil(t *testing.T) {
	mockRepo := createMockRepository(t)
	mockRepo.EXPECT().Exists("lacin").Return(true).Times(1)

	mockWallet := model.Wallet{Username: "lacin", Balance: 70}
	mockRepo.EXPECT().Get("lacin").Return(mockWallet).Times(1)

	s := service.NewWallet(nil, mockRepo)
	wallet, err := s.Get("lacin")

	assert.Equal(t, mockWallet, wallet)
	assert.Nil(t, err)
}

func TestCreate_WalletAlreadyExists_ReturnsExistingWallet(t *testing.T) {
	mockRepo := createMockRepository(t)
	mockRepo.EXPECT().Exists("lacin").Return(true).Times(1)

	mockWallet := model.Wallet{Username: "lacin", Balance: 0}
	mockRepo.EXPECT().Get("lacin").Return(mockWallet).Times(1)

	s := service.NewWallet(&config.Conf{}, mockRepo)
	wallet := s.Create("lacin")

	assert.Equal(t, mockWallet, wallet)
}

func TestCreate_ReturnsNewWallet(t *testing.T) {
	mockRepo := createMockRepository(t)
	mockRepo.EXPECT().Exists("yulet").Return(false).Times(1)

	mockWallet := &model.Wallet{Username: "yulet", Balance: 0}
	mockRepo.EXPECT().Save(mockWallet).Return().Times(1)

	s := service.NewWallet(&config.Conf{}, mockRepo)
	wallet := s.Create("yulet")

	assert.Equal(t, *mockWallet, wallet)
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		desc           string
		username       string
		newBalance     int
		exists         bool
		currWallet     model.Wallet
		expectedWallet model.Wallet
		expectedError  error
	}{
		{
			desc:           "given lacin and 400, expect lacin's updated wallet and nil",
			username:       "lacin",
			newBalance:     400,
			exists:         true,
			currWallet:     model.Wallet{Username: "lacin", Balance: 0},
			expectedWallet: model.Wallet{Username: "lacin", Balance: 400},
			expectedError:  nil,
		},
		{
			desc:           "given doga and 100, expect nil and wallet not exists error",
			username:       "doga",
			newBalance:     100,
			exists:         false,
			currWallet:     model.Wallet{},
			expectedWallet: model.Wallet{},
			expectedError:  service.ErrWalletNotExists,
		},
		{
			desc:           "given lacin and -300, expect nil and balance below limit error",
			username:       "lacin",
			newBalance:     -300,
			exists:         true,
			currWallet:     model.Wallet{Username: "lacin", Balance: 0},
			expectedWallet: model.Wallet{},
			expectedError:  service.ErrBalanceBelowLimit,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockRepo := createMockRepository(t)
			mockRepo.EXPECT().Exists(tC.username).Return(tC.exists)
			if tC.exists {
				mockRepo.EXPECT().Get(tC.username).Return(tC.currWallet)
			}
			if tC.expectedWallet != (model.Wallet{}) {
				mockRepo.EXPECT().Update(tC.username, tC.newBalance).Return(tC.expectedWallet)
			}

			s := service.NewWallet(&config.Conf{MinimumBalanceAmount: -100}, mockRepo)
			updatedWallet, err := s.Update(tC.username, tC.newBalance)

			assert.Equal(t, tC.expectedWallet, updatedWallet)
			assert.Equal(t, tC.expectedError, err)
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		desc        string
		username    string
		exists      bool
		expectedErr error
	}{
		{
			desc:        "wallet belonging to given username exists, return nil error",
			username:    "lacin",
			exists:      true,
			expectedErr: nil,
		},
		{
			desc:        "wallet belonging to given username does not exist, return not exists error",
			username:    "doga",
			exists:      false,
			expectedErr: service.ErrWalletNotExists,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			mockRepo := createMockRepository(t)
			mockRepo.EXPECT().Exists(tC.username).Return(tC.exists)

			if tC.exists {
				mockRepo.EXPECT().Delete(tC.username)
			}

			s := service.NewWallet(&config.Conf{}, mockRepo)
			err := s.Delete(tC.username)

			assert.Equal(t, tC.expectedErr, err)
		})
	}
}
