package logger

// Config holds logger configuration.
type Config struct {
	Level           string // debug|info|warn|error
	Format          string // json|text
	Output          string // stdout or filepath
	RequestIDHeader string // header name for request id
}

// DefaultConfig returns sensible defaults.
func DefaultConfig() Config {
	return Config{
		Level:           "info",
		Format:          "json",
		Output:          "stdout",
		RequestIDHeader: "X-Request-Id",
	}
}
