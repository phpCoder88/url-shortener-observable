package interfaces

type URLVisitRepository interface {
	AddURLVisit(int64, string) error
}
