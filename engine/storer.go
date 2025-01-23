package engine

type KVStorer interface {
	Get(string) (any, error)
	Store(string, any) error
	StoreOW(string, any) error
	Update(string, any) error
	Delete(string) error
	Flush() error
	Load() error
}
