/*
Package builder materializes config entries from res into runtime resources.

Configs are registered in res by library use init (AddWithTags) and by
ecfg.Register from AppConfig. Build walks the registry, calls Build() on
every entry that implements Builder, registers the result (inheriting entry
tags), and removes the config entry.

Builder does not perform DI — wiring happens later via sdi.

Typical pipeline:

	cfg, _ := ecfg.Parse(...)
	ecfg.Register(cfg, res.Default)
	builder.Build(res.Default)
	res.Transform(...)
	sdi.Resolve(res.Default)
*/
package builder
