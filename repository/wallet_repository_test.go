package repository_test

import (
	"testing"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/model"
	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/repository"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	wallets := map[string]*model.Wallet{
		"lacin": {
			Username: "lacin",
			Balance:  100,
		},
		"fatma": {
			Username: "fatma",
			Balance:  -50,
		},
	}
	r := repository.NewWallet(wallets)

	testCases := []struct {
		desc          string
		givenUsername string
		expectedValue bool
	}{
		{
			desc:          "given name that does not exist, expect false",
			givenUsername: "doga",
			expectedValue: false,
		},
		{
			desc:          "given name that exists, expect true",
			givenUsername: "lacin",
			expectedValue: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			exists := r.Exists(tC.givenUsername)
			assert.Equal(t, tC.expectedValue, exists)
		})
	}
}

func TestGetAll(t *testing.T) {
	wallets := map[string]*model.Wallet{
		"lacin": {
			Username: "lacin",
			Balance:  100,
		},
		"fatma": {
			Username: "fatma",
			Balance:  -50,
		},
	}
	r := repository.NewWallet(wallets)

	result := r.GetAll()

	assert.Equal(t, wallets, result)
}

func TestGet(t *testing.T) {
	wallets := map[string]*model.Wallet{
		"lacin": {
			Username: "lacin",
			Balance:  100,
		},
		"fatma": {
			Username: "fatma",
			Balance:  -50,
		},
	}
	r := repository.NewWallet(wallets)

	testCases := []struct {
		desc           string
		givenUsername  string
		expectedWallet *model.Wallet
	}{
		{
			desc:          "given lacin, expect lacin's wallet",
			givenUsername: "lacin",
			expectedWallet: &model.Wallet{
				Username: "lacin",
				Balance:  100,
			},
		},
		{
			desc:          "given fatma, expect fatma's wallet",
			givenUsername: "fatma",
			expectedWallet: &model.Wallet{
				Username: "fatma",
				Balance:  -50,
			},
		},
		{
			desc:           "given non existant name, expect zero values",
			givenUsername:  "doga",
			expectedWallet: nil,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := r.Get(tC.givenUsername)

			assert.Equal(t, tC.expectedWallet, result)
		})
	}
}

func TestCreate(t *testing.T) {
	wallets := map[string]*model.Wallet{
		"lacin": {
			Username: "lacin",
			Balance:  100,
		},
		"fatma": {
			Username: "fatma",
			Balance:  -50,
		},
	}
	r := repository.NewWallet(wallets)

	newWallet := r.Create("halil", 0)

	assert.Contains(t, wallets, newWallet.Username)
	assert.Equal(t, "halil", newWallet.Username)
	assert.Equal(t, 0, newWallet.Balance)
}

func TestUpdate(t *testing.T) {
	wallets := map[string]*model.Wallet{
		"lacin": {
			Username: "lacin",
			Balance:  100,
		},
		"fatma": {
			Username: "fatma",
			Balance:  -50,
		},
	}
	r := repository.NewWallet(wallets)

	testCases := []struct {
		desc           string
		givenUsername  string
		givenBalance   int
		expectedWallet *model.Wallet
	}{
		{
			desc:          "given lacin and -20, expect wallet's balance to be 80",
			givenUsername: "lacin",
			givenBalance:  -20,
			expectedWallet: &model.Wallet{
				Username: "lacin",
				Balance:  80,
			},
		},
		{
			desc:          "given fatma and 50, expect wallet's balance to be 0",
			givenUsername: "fatma",
			givenBalance:  50,
			expectedWallet: &model.Wallet{
				Username: "fatma",
				Balance:  0,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := r.Update(tC.givenUsername, tC.givenBalance)

			assert.Equal(t, tC.expectedWallet, result)
		})
	}
}
