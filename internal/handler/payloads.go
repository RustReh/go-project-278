package handler

// createLinkPayload — POST /api/links
type createLinkPayload struct {
	OriginalURL string `json:"original_url" binding:"required,url,max=2048"`
	ShortName   string `json:"short_name" binding:"omitempty,min=3,max=32"`
}

// updateLinkPayload — PUT /api/links/:id
type updateLinkPayload struct {
	OriginalURL string `json:"original_url" binding:"required,url,max=2048"`
	ShortName   string `json:"short_name" binding:"required,min=3,max=32"`
}
