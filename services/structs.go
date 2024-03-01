package services

type tokenValidationRequest struct {
	token string `json:"token"`
}

type tokenValidationResponse struct {
	err error `json:"err,invalidRoken"`
}

type tokenGenerationRequest struct {
	userID uint64 `json:"userID"`
}
type tokenGenerationResponse struct {
	token string `json:"token"`
	err   error  `json:"err,couldnGenerateToken"`
}
