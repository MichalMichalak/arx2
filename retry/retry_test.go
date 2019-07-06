package retry_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/unravelin/core/arx2/retry"
	"testing"
)

func TestRetry(t *testing.T) {
	c := 0
	f := func() error {
		fmt.Println("inside function, count", c)
		if c == 20 {
			return nil
		}
		c++
		return errors.New("test")
	}
	r := retry.DefaultRetryer{MaxTries: 3, WaitInterval: 1000, Multiplier: 2}
	err := r.Retry(f)
	require.Error(t, err)
	require.Equal(t, "failed to try 3 times", err.Error())
}

func TestRetry2(t *testing.T) {
	c := 0
	var act retry.Action = func() error {
		fmt.Println("inside function, count", c)
		if c == 3 {
			return nil
		}
		c++
		return errors.New("test")
	}
	var succChk retry.SuccessCheck = func(e error) bool {
		return e == nil
	}
	r := retry.DefaultRetryer{MaxTries: 3, WaitInterval: 500, Multiplier: 2}
	err := r.RetryWithSuccessCheck(act, succChk)
	assert.NotNil(t, err)
	assert.Equal(t, "failed to try 3 times", err.Error())
}
