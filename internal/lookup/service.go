package lookup

type Verdict string

const (
	VerdictSafe      Verdict = "safe"
	VerdictMalicious Verdict = "malicious"
)

type Result struct {
	NormalizedURL string  `json:"normalized_url"`
	Verdict       Verdict `json:"verdict"`
	Matched       bool    `json:"matched"`
	Reason        string  `json:"reason,omitempty"`
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Lookup(normalizedURL string) Result {
	if s.store.Contains(normalizedURL) {
		return Result{
			NormalizedURL: normalizedURL,
			Verdict:       VerdictMalicious,
			Matched:       true,
			Reason:        "known malware URL",
		}
	}

	return Result{
		NormalizedURL: normalizedURL,
		Verdict:       VerdictSafe,
		Matched:       false,
	}
}
