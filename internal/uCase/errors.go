package uCase

import "errors"

var (
	ErrTraceProviderIsNotSet   = errors.New("the global trace provider isn't set")
	ErrRequestAlreadyProcessed = errors.New("request already processed")
)
