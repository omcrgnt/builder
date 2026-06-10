/*
Пакет builder — сборка ресурсов из cfg перед регистрацией в res.

Обходит поля структуры первого уровня; для каждого поля, реализующего Builder,
вызывает Build() и передаёт результат в Registrar (например, res).

Builder не выполняет DI — связывание зависимостей выполняется позже через sdi.

Типичный pipeline:

	cfg, _ := ecfg.Parse(...)
	builder.Build(cfg, res)
	res.Transform(...)
	sdi.Resolve(res)
*/
package builder
