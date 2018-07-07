package config

// Interface for fetching configuration settings
type Config interface {
	// Get an optional setting by name, with a default value to return if it does
	// exist.
	Optional(key string, defaultValue string) string

	// Get a required setting by name.
	Required(key string) (string, error)
}
