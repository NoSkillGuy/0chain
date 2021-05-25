package magmasc

import (
	"errors"
)

var (
	// ErrStakePoolExist represents an error that can occur while a new stake pool is creating
	// and saving in state.StateContextI when stake pool for provided node already exists in Magma Smart Contract.
	ErrStakePoolExist = errors.New("stake pool already exist")
)
