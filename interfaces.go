package builder

// NewResourceer is a wire type whose resource is created in [Build] via [newResourceSpec].
type NewResourceer interface {
	NewResource() (any, error)
}

// BuildConfiger registers a config spec in the registry during seed.
type BuildConfiger interface {
	BuildConfig() (Builder, error)
}
