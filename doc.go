/*
Package builder registers and materializes config specs in a [Registry].

AppResources pipeline:

	builder.Seed(reg, &appResources)  // BuildConfiger → spec; NewResourceer → deferred build
	ecfg.Apply(reg, &appResources, …) // env into specs
	builder.Build(reg)                // Spec.Build() / NewResource() → resources

Pass [github.com/omcrgnt/res.ForBuilder] when using [github.com/omcrgnt/res.Registry].
Builder does not perform DI — wiring happens later via [sdi.Resolve].
*/
package builder
