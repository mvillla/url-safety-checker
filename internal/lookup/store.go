package lookup

type Store interface {
	Contains(url string) bool
}

type MemoryStore struct {
	urls map[string]struct{}
}

func NewMemoryStore(urls []string) *MemoryStore {
	store := &MemoryStore{
		urls: make(map[string]struct{}, len(urls)),
	}

	for _, url := range urls {
		store.urls[url] = struct{}{}
	}

	return store
}

func (s *MemoryStore) Contains(url string) bool {
	_, ok := s.urls[url]
	return ok
}
