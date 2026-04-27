package fraud

type Scorer struct{}

func NewScorer() *Scorer {
	return &Scorer{}
}

func (s *Scorer) Score(req *Request) Response {
	return Response{
		Approved:   true,
		FraudScore: 0,
	}
}
