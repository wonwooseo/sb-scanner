package sentiment

import "context"

type Sentiment struct {
	Score float64 // -1.0 (negative) to 1.0 (positive)
	Model string  // name of the model used for evaluation
}

type Evaluator interface {
	Evaluate(context.Context, string) (Sentiment, error)
}
