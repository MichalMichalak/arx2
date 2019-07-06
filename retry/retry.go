package retry

import (
	"fmt"
	"time"
)

type Retryer interface {
	Retry(action Action) error

	// WithRetry2 executes `action` multiple times until it succeeds or hit the failures limit and returns error.
	// A `success` function is called to determine if action ended up successfuly (not all errors must always be
	// treated as failure.)
	RetryWithSuccessCheck(action Action, check SuccessCheck) error
}

type Action func() error
type SuccessCheck func(error) bool

type DefaultRetryer struct {
	MaxTries     int
	WaitInterval int
	Multiplier   int
}

func (r *DefaultRetryer) Retry(action Action) error {
	count := 0
	interval := r.WaitInterval
	var errs []error
	for {
		fmt.Println("TRY NUMBER", count)
		err := action()
		count++
		if err == nil {
			return nil
		}
		fmt.Println("TRY", count, "FAILED")
		errs = append(errs, err)
		if count == r.MaxTries {
			return RetryError{Errors: errs, Tries: count}
		}
		fmt.Println("WAITING", interval, "MS")
		time.Sleep(time.Duration(interval) * time.Millisecond)
		if r.Multiplier != 0 {
			interval = interval * r.Multiplier
		}
	}
}

func (r *DefaultRetryer) RetryWithSuccessCheck(action Action, successCheck SuccessCheck) error {
	count := 0
	interval := r.WaitInterval
	var errs []error
	for {
		fmt.Println("TRY NUMBER", count)
		err := action()
		count++
		if successCheck(err) {
			return nil
		}
		fmt.Println("TRY", count, "FAILED")
		errs = append(errs, err)
		if count == r.MaxTries {
			return RetryError{Errors: errs, Tries: count}
		}
		fmt.Println("WAITING", interval, "MS")
		time.Sleep(time.Duration(interval) * time.Millisecond)
		if r.Multiplier != 0 {
			interval = interval * r.Multiplier
		}
	}
}

type RetryError struct {
	Errors []error
	Tries  int
}

func (re RetryError) Error() string {
	return fmt.Sprintf("failed to try %d times", re.Tries)
}

func (RetryError) RuntimeError() {
	panic("implement me")
}
