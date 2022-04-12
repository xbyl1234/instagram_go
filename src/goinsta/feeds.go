package goinsta

type FeedInterface interface {
	NextPost() (*Media, error)
}

type Feed struct {
	FeedInterface
	Inst       *Instagram
	media      []*Media
	indexMedia int
}
