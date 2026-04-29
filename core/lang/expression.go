package lang

type Expression struct {
	Open  string
	Close string
}

const (
	DefaultOpenExpr  = "{{"
	DefaultCloseExpr = "}}"
)

func NewExpression() *Expression {
	return &Expression{
		Open:  DefaultOpenExpr,
		Close: DefaultCloseExpr,
	}
}

func NewExpressionConfigurable(open string, close string) *Expression {
	return &Expression{
		Open:  open,
		Close: close,
	}
}
