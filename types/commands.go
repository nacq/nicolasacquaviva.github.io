package types

type Command struct {
	Name            string `json:"name"`
	ExecutionStatus bool   `json:"executionStatus"`
	Date            string `json:"date"`
	ClientIP        string `json:"clientIP"`
	UserAgent       string `json:"userAgent"`
}
