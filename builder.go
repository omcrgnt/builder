package builder

import "fmt"

// Build materializes every config entry in reg that implements [Builder]:
// Build(), register the resource (inheriting entry tags), remove the config.
// Non-Builder entries are left unchanged.
func Build(reg Registry) error {
	if reg == nil {
		return fmt.Errorf("builder: nil registry")
	}

	type job struct {
		config any
		tags   []any
	}

	var jobs []job
	reg.WalkEntries(func(value any, tags []any) bool {
		b, ok := value.(Builder)
		if !ok {
			return true
		}
		jobs = append(jobs, job{config: b, tags: tags})
		return true
	})

	for _, j := range jobs {
		built, err := j.config.(Builder).Build()
		if err != nil {
			return fmt.Errorf("builder: %T: %w", j.config, err)
		}

		if len(j.tags) > 0 {
			if err := reg.AddWithTags(built, j.tags...); err != nil {
				return fmt.Errorf("builder: %T: %w", j.config, err)
			}
		} else if err := reg.Add(built); err != nil {
			return fmt.Errorf("builder: %T: %w", j.config, err)
		}

		if err := reg.Remove(j.config); err != nil {
			return fmt.Errorf("builder: %T: remove config: %w", j.config, err)
		}
	}

	return nil
}
