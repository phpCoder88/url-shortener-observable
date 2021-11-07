package interfaces

import "context"

type URLVisitRepository interface {
	AddURLVisit(context.Context, int64, string) error
}
