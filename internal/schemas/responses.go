package schemas

type PingResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type LinkResponse struct {
	ID          int64  `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
	ShortURL    string `json:"short_url"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Payload any    `json:"payload,omitempty"`
}

type CreateUpdateLinkRequest struct {
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
}

type LinkVisitResponse struct {
	ID        int64  `json:"id"`
	LinkID    int64  `json:"link_id"`
	CreatedAt string `json:"created_at"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Status    int    `json:"status"`
}
