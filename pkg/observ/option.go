package observ

type Option struct {
	ServiceName  string
	CollectorURL string // Without https:// (example: localhost:4317)
	Headers      map[string]string
	Insecure     bool
}
