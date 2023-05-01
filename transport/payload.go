package transport

import "fmt"

type ConvertedPriceRequest struct {
	TargetCurrency string `json:"target_currency"`
	Amount         int    `json:"amount"`
}

type ConvertedPriceResponse struct {
	Total float64 `json:"total"`
	Err   string  `json:"err,omitempty"`
}

type ErrorResponse struct {
	Err string `json:"err, omitEmpty"`
}

func (e *ErrorResponse) Error() string {
	return e.Err
}

func (e *ErrorResponse) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"err":"%s"}`, e.Err)), nil
}
