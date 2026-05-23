package domain

type Link struct {
	Id          int64
	OriginalUrl string
	ShortName   string
	ShortUrl    string
}

type LinkVO struct {
	OriginalUrl string
	ShortName   string
}

type LinkShortenedVO struct {
	OriginalUrl string
	ShortName   string
	ShortUrl    string
}
