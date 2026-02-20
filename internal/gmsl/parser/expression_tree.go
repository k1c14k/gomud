package parser

type ExpressionTree struct {
	root    Expression
	current *Expression
}

func NewExpressionTree() *ExpressionTree {
	return &ExpressionTree{}
}

func (et *ExpressionTree) GetExpression() Expression {
	return et.root
}

func (et *ExpressionTree) AddExpression(expression Expression) {
	if et.root == nil {
		et.root = expression
	} else {
		switch e := expression.(type) {
		case *BinaryExpression:
			e.Left, et.root = et.root, e
		default:
			et.root.(*BinaryExpression).Right = expression
		}
	}
}

func (et *ExpressionTree) CanAddLeaf() bool {
	if et.root == nil {
		return true
	}

	elem := et.root
	if _, ok := elem.(*BinaryExpression); ok {
		return elem.(*BinaryExpression).Right == nil
	}

	return false
}

func (et *ExpressionTree) CanAddBranch() bool {
	return !et.CanAddLeaf()
}
