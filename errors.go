package crex

import "errors"

var (
	ErrNotImplemented    = errors.New("not implement")
	ErrWebSocketDisabled = errors.New("websocket disabled")
	ErrApiKeysRequired   = errors.New("api keys required")

	ErrInvalidAmount = errors.New("amount is not valid")
)
