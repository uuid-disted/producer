package generator

import "time"

type UUIDGenerator interface {
	Generate(t time.Time) (string, error)
}
