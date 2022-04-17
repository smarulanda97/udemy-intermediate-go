package utils

type JsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `jsond:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}
