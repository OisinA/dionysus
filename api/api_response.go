package api

type APIResponse struct {
	StatusCode int         `json:"status_code"`
	Content    interface{} `json:"content"`
}
