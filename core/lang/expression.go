package lang

type Expression struct {
	Open  string
	Close string
}

func NewExpression(open string, close string) *Expression {
	return &Expression{
		Open:  open,
		Close: close,
	}
}
