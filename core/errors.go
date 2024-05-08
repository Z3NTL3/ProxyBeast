package core

import "errors"

var (
	ErrPropsInvalid = errors.New("[props] of invalid type")
	ErrPropsOPInvalid = errors.New("[props] operation is invalid")
)