package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGivenValidParams_WhenCallsNewOrder_ThenShouldReceiveCreateOrderWithAllParams(t *testing.T) {
	order, err := NewOrder("123", 100.0, 10.0)
	assert.NoError(t, err)
	assert.Equal(t, "123", order.ID)
	assert.Equal(t, 100.0, order.Price)
	assert.Equal(t, 10.0, order.Tax)
	assert.Nil(t, order.IsValid())
}

func TestGivenEmptyID_WhenCallsNewOrder_ThenShouldReceiveAnError(t *testing.T) {
	order := Order{}
	assert.Error(t, order.IsValid(), "invalid id")
}

func TestGivenInvalidPrice_WhenCallsNewOrder_ThenShouldReceiveAnError(t *testing.T) {
	order := Order{ID: "23"}
	assert.Error(t, order.IsValid(), "invalid price")
}
func TestGivenInvalidTax_WhenCallsNewOrder_ThenShouldReceiveAnError(t *testing.T) {
	order := Order{ID: "23", Price: 10}
	assert.Error(t, order.IsValid(), "invalid tax")
}

func TestGivenAPriceAndTax_WhenCallsCalculatePrice_ThenShouldUpdateFinalPrice(t *testing.T) {
	order, err := NewOrder("123", 100.0, 10.0)
	assert.Nil(t, err)

	assert.Nil(t, order.CalculateFinalPrice())
	assert.Equal(t, 1000.0, order.FinalPrice)
}
