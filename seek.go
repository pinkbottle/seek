package seek

type Resource struct {
	Content string
	URL     string
}

type Result struct {
	Content string
	URL     string
	Score   float64
}
