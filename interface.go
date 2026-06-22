package builder

// Registry is the minimal resource store [Seed] and [Build] need.
// [github.com/omcrgnt/res.Registry] satisfies it via [github.com/omcrgnt/res.ForBuilder].
type Registry interface {
	Add(v any) error
	AddWithTags(v any, tags ...any) error
	WalkEntries(fn func(value any, tags []any) bool)
	Remove(v any) error
}

// Builder is a config spec registered in the store before materialization.
type Builder interface {
	Build() (any, error)
}

// NewResourceer is a wire type whose resource is created in [Build] via [newResourceSpec].
type NewResourceer interface {
	NewResource() (any, error)
}

// BuildConfiger registers a config spec in the registry during seed.
type BuildConfiger interface {
	BuildConfig() (Builder, error)
}
