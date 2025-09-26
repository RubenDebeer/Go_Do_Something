package port

// Repository is the port for persistence, implemented by the file adapter.
type Repository interface {
	Add(value string) error
	List() ([]string, error)
	DeleteLast() error
}
