package domain

type CreateLinkResponse struct {
	LinkURL  string `json:"linkUrl"`
	LinkCode string `json:"linkCode"`
}

type NotifyRequest struct {
	UserID string `json:"userId"`
	Text   string `json:"text"`
}
