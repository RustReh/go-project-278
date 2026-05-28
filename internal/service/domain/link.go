package domain

type Link struct {
	ID          int64
	OriginalURL string
	ShortName   string
	ShortURL    string
}

type LinkVO struct {
	OriginalURL string
	ShortName   string
}

type LinkShortenedVO struct {
	OriginalURL string
	ShortName   string
	ShortURL    string
}
