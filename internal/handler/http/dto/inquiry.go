package dto

// AskRequest represents the request payload for asking a question
type AskRequest struct {
	Msg string `json:"msg"`
}

// AskResponse represents the response payload for a question
type AskResponse struct {
	Response string `json:"response"`
}
