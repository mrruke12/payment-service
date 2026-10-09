package storage

type Scanner interface {
	Scan(dest ...any) error
}
