package model

type Stock struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Exchange string `json:"exchange"`
	MICCode  string `json:"mic_code"`
	Country  string `json:"country"`
	Type     string `json:"type"`
	FIGICode string `json:"figi_code"`
	CFICode  string `json:"cfi_code"`
	ISIN     string `json:"isin"`
	CUSIP    string `json:"cusip"`
}