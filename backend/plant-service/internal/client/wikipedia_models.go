package client

type PlantDescriptionResponse struct {
	Description string       `json:"description"`
	Extract     string       `json:"extract"`
	ExtractHTML string       `json:"extract_html"`
	ContentURLs *ContentURLs `json:"content_urls"`
}

type ContentURLs struct {
	Desktop *DesktopURL `json:"desktop"`
}

type DesktopURL struct {
	Page string `json:"page"`
}
