package util_test

import (
	"testing"

	"github.com/claytten/golang-simplebank/internal/util"
	"github.com/stretchr/testify/require"
)

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	IDR = "IDR"
)

func TestIsSupportCurrency(t *testing.T) {
	testCases := []struct {
		name string
		val  string
		res  bool
	}{
		//Valid USD
		{
			name: USD,
			val:  USD,
			res:  true,
		},
		//Valid EUR
		{
			name: EUR,
			val:  EUR,
			res:  true,
		},
		// Valid CAD
		{
			name: CAD,
			val:  CAD,
			res:  true,
		},
		//Invalid
		{
			name: IDR,
			val:  IDR,
			res:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			isValid := util.IsSupportCurrency(tc.val)
			require.Equal(t, tc.res, isValid)
		})
	}
}
