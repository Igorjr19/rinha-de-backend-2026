package fraud

import "sync/atomic"

type Scorer struct {
	ready atomic.Bool
}

func NewScorer() *Scorer {
	s := &Scorer{}
	s.ready.Store(true)
	return s
}

func (s *Scorer) Ready() bool {
	return s.ready.Load()
}

func (s *Scorer) Score(req *Request) Response {
	return Response{
		Approved:   true,
		FraudScore: 0,
	}
}
