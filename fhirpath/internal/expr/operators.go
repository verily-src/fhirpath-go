package expr

// Operator constants.
const (
	Equals        = "="
	NotEquals     = "!="
	Equivalence   = "~"
	Inequivalence = "!~"
	Is            = "is"
	As            = "as"
	And           = "and"
	Or            = "or"
	Xor           = "xor"
	Implies       = "implies"
	Lt            = "<"
	Gt            = ">"
	Lte           = "<="
	Gte           = ">="
	Add           = "+"
	Sub           = "-"
	Concat        = "&"
	Mul           = "*"
	Div           = "/"
	FloorDiv      = "div"
	Mod           = "mod"
	In            = "in"
	Contains      = "contains"
)

// Operator represents a valid expression operator.
type Operator string
