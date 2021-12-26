package mathmaker

type Mathmaker interface {
	MakeMathEquation(Operations) string
	GetEquationResult() (result string)
}

type mathMake struct {
	a      float64
	b      float64
	oper   Operation
	result string
}

type Operation interface {
	Perform(a, b float64) float64
	GetName() string
}

type Sum struct {
	name string
}
type Minus struct {
	name string
}
type Mult struct {
	name string
}
type Div struct {
	name string
}

type Operations []Operation
