package shortener

type CreateAddressRequest struct {
	URL string `json:"url"`
}

type CreateAddressResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
