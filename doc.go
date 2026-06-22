/*
Package builder registers and materializes config specs in [res.Registry].

AppResources pipeline:

	builder.Seed(reg, &appResources)  // BuildConfiger → spec; NewResourceer → deferred build
	ecfg.Apply(reg, &appResources, …) // env into specs
	builder.Build(reg)                // Spec.Build() / NewResource() → resources

Library use init may still [res.AddWithTags] replaceable defaults.
Builder does not perform DI — wiring happens later via [sdi.Resolve].
*/
package builder
