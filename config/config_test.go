package config_test

import (
	"testing"

	"github.com/152-Modanisa-FullStack-Bootcamp/week-6-assignment-gokcelb/config"
	"github.com/stretchr/testify/assert"
)

func TestGetconf(t *testing.T) {
	c := config.Getconf()

	assert.Equal(t, 0, c.InitialBalanceAmount)
	assert.Equal(t, -100, c.MinimumBalanceAmount)
}
