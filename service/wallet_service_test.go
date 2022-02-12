package service_test

import (
	"testing"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/mocks"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	mockRepo := createMockRepository(t)
	mockWallets := map[string]*model.Wallet{
		"yuksel": {
			Username: "yuksel",
			Balance:  350,
		},
		"lacin": {
			Username: "lacin",
			Balance:  0,
		},
	}
	mockRepo.EXPECT().GetAll().Return(mockWallets)

	s := service.NewWallet(mockRepo)
	wallets := s.GetAll()

	assert.Equal(t, mockWallets, wallets)
}

func TestGet_WalletDoesNotExist_ReturnsNilAndErrWalletNotExists(t *testing.T) {
	mockRepo := createMockRepository(t)
	mockRepo.EXPECT().Exists("doga").Return(false)

	s := service.NewWallet(mockRepo)
	wallet, err := s.Get("doga")

	assert.Nil(t, wallet)
	assert.EqualError(t, err, service.ErrWalletNotExists.Error())
}

func TestGet_WalletExists_ReturnsWalletAndNil(t *testing.T) {
	mockRepo := createMockRepository(t)
	mockRepo.EXPECT().Exists("lacin").Return(true).Times(1)

	mockWallet := &model.Wallet{
		Username: "lacin",
		Balance:  70,
	}
	mockRepo.EXPECT().Get("lacin").Return(mockWallet).Times(1)

	s := service.NewWallet(mockRepo)
	wallet, err := s.Get("lacin")

	assert.Equal(t, mockWallet, wallet)
	assert.Nil(t, err)
}

func TestCreate_WalletAlreadyExists_ReturnsExistingWallet(t *testing.T) {
	mockRepo := createMockRepository(t)
	mockRepo.EXPECT().Exists("lacin").Return(true).Times(1)

	mockWallet := &model.Wallet{
		Username: "lacin",
		Balance:  0,
	}
	mockRepo.EXPECT().Get("lacin").Return(mockWallet).Times(1)

	s := service.NewWallet(mockRepo)
	wallet := s.Create("lacin")

	assert.Equal(t, mockWallet, wallet)
}

func TestCreate_ReturnsNewWallet(t *testing.T) {
	mockRepo := createMockRepository(t)
	mockRepo.EXPECT().Exists("yulet").Return(false).Times(1)

	mockWallet := &model.Wallet{
		Username: "yulet",
		Balance:  -10,
	}
	mockRepo.EXPECT().Create("yulet", 0).Return(mockWallet).Times(1)

	s := service.NewWallet(mockRepo)
	wallet := s.Create("yulet")

	assert.Equal(t, mockWallet, wallet)
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		desc           string
		username       string
		newBalance     int
		exists         bool
		currWallet     *model.Wallet
		expectedWallet *model.Wallet
		expectedError  error
	}{
		{
			desc:       "given lacin and 400, expect lacin's updated wallet and nil",
			username:   "lacin",
			newBalance: 400,
			exists:     true,
			currWallet: &model.Wallet{
				Username: "lacin",
				Balance:  0,
			},
			expectedWallet: &model.Wallet{
				Username: "lacin",
				Balance:  400,
			},
			expectedError: nil,
		},
		{
			desc:           "given doga and 100, expect nil and wallet not exists error",
			username:       "doga",
			newBalance:     100,
			exists:         false,
			currWallet:     nil,
			expectedWallet: nil,
			expectedError:  service.ErrWalletNotExists,
		},
		{
			desc:       "given lacin and -300, expect nil and balance below limit error",
			username:   "lacin",
			newBalance: -300,
			exists:     true,
			currWallet: &model.Wallet{
				Username: "lacin",
				Balance:  0,
			},
			expectedWallet: nil,
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
			if tC.expectedWallet != nil {
				mockRepo.EXPECT().Update(tC.username, tC.newBalance).Return(tC.expectedWallet)
			}

			s := service.NewWallet(mockRepo)
			updatedWallet, err := s.Update(tC.username, tC.newBalance)

			assert.Equal(t, tC.expectedWallet, updatedWallet)
			assert.Equal(t, tC.expectedError, err)
		})
	}
}

func createMockRepository(t *testing.T) *mocks.MockWalletRepository {
	return mocks.NewMockWalletRepository(gomock.NewController(t))
}
