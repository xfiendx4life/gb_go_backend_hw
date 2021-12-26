package mathmaker

import (
	"fmt"
	"math/rand"
	"time"
)

func (o *Sum) Perform(a, b float64) float64 {
	return a + b
}

func (o *Sum) GetName() string {
	return o.name
}
func (o *Minus) Perform(a, b float64) float64 {
	return a - b
}
func (o *Minus) GetName() string {
	return o.name
}

func (o *Mult) Perform(a, b float64) float64 {
	return a * b
}

func (o *Mult) GetName() string {
	return o.name

}
func (o *Div) Perform(a, b float64) float64 {
	return a / b
}
func (o *Div) GetName() string {
	return o.name
}

func NewMathmaker() Mathmaker {
	return &mathMake{}
}

func (m *mathMake) MakeMathEquation(ops Operations) string {
	rand.Seed(time.Now().UnixNano())
	m.a = rand.Float64() + float64(rand.Intn(100))
	rand.Seed(time.Now().UnixNano())
	m.b = rand.Float64() + float64(rand.Intn(100))
	m.oper = ops[rand.Intn(len(ops))]
	m.result = fmt.Sprintf("%.3f", m.oper.Perform(m.a, m.b))
	return fmt.Sprintf("%.3f %s %.3f", m.a, m.oper.GetName(), m.b)
}

func (m *mathMake) GetEquationResult() (result string) {
	return m.result
}

func CreateDefaultOperations() Operations {
	s := &Sum{"+"}
	m := &Minus{"-"}
	mul := &Mult{"*"}
	div := &Div{"/"}
	return Operations{s, m, mul, div}
}
