package generator

import "time"

type GeneratorConfig struct {
	UseBuffer bool
}

type Generator interface {
	Generate(t time.Time) (string, error)
}
