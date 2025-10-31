package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/expr"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/grammar"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/reflection"
	"github.com/verily-src/fhirpath-go/fhirpath/system"
	"github.com/verily-src/fhirpath-go/internal/resource"
	"github.com/verily-src/fhirpath-go/internal/slices"
)

var (
	errNotSupported       = errors.New("expression not currently supported")
	errTooManyQualifiers  = errors.New("too many type qualifiers")
	errVisitingChildren   = errors.New("error while visiting child expressions")
	errUnresolvedFunction = errors.New("function identifier can't be resolved")
)

type FHIRPathVisitor struct {
	*grammar.BasefhirpathVisitor
	visitedRoot bool
	Functions   funcs.FunctionTable
	Transform   VisitorTransform
	Permissive  bool
}

type VisitResult struct {
	Result expr.Expression
	Error  error
}

type typeResult struct {
	result reflection.TypeSpecifier
	err    error
}

// clone produces a shallow-clone of the visitor, to be used when visiting sub-expressions.
func (v *FHIRPathVisitor) clone() *FHIRPathVisitor {
	return &FHIRPathVisitor{
		Functions:   v.Functions,
		Transform:   v.Transform,
		Permissive:  v.Permissive,
		visitedRoot: false,
	}
}

func (v *FHIRPathVisitor) transformedVisitResult(resultExpr expr.Expression) *VisitResult {
	if v.Transform == nil {
		v.Transform = IdentityTransform
	}
	return &VisitResult{v.Transform(resultExpr), nil}
}

func (v *FHIRPathVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

func (v *FHIRPathVisitor) VisitProg(ctx *grammar.ProgContext) interface{} {
	return v.Visit(ctx.Expression()).(*VisitResult)
}

// VisitIndexerExpression visits both the left side expression and right side expression, and
// constructs an index expression. If the right side expression does not evaluate to an Integer,
// returns an error
func (v *FHIRPathVisitor) VisitIndexerExpression(ctx *grammar.IndexerExpressionContext) interface{} {
	// visit left side expression
	leftResult := v.Visit(ctx.Expression(0)).(*VisitResult)
	if leftResult.Error != nil {
		return &VisitResult{nil, leftResult.Error}
	}

	// visit contained expression with new Visitor to reset root node, and construct index
	rightResult := v.clone().Visit(ctx.Expression(1)).(*VisitResult)
	if rightResult.Error != nil {
		return &VisitResult{nil, rightResult.Error}
	}
	indexExpr := &expr.IndexExpression{Index: rightResult.Result}

	sequence := &expr.ExpressionSequence{Expressions: []expr.Expression{leftResult.Result, indexExpr}}
	return v.transformedVisitResult(sequence)
}

func (v *FHIRPathVisitor) VisitPolarityExpression(ctx *grammar.PolarityExpressionContext) interface{} {
	operator := expr.Operator(ctx.GetChild(0).(antlr.TerminalNode).GetText())
	result := v.Visit(ctx.Expression()).(*VisitResult)
	if result.Error != nil {
		return &VisitResult{nil, result.Error}
	}

	// Return initial expression if the operator is '+'
	if operator != expr.Sub {
		return v.transformedVisitResult(result.Result)
	}

	return v.transformedVisitResult(&expr.NegationExpression{Expr: result.Result})
}

func (v *FHIRPathVisitor) VisitAdditiveExpression(ctx *grammar.AdditiveExpressionContext) interface{} {
	leftResult := v.Visit(ctx.Expression(0)).(*VisitResult)
	if leftResult.Error != nil {
		return &VisitResult{nil, leftResult.Error}
	}
	rightResult := v.clone().Visit(ctx.Expression(1)).(*VisitResult)
	if rightResult.Error != nil {
		return &VisitResult{nil, rightResult.Error}
	}

	operator := expr.Operator(ctx.GetChild(1).(antlr.TerminalNode).GetText())

	var expression expr.Expression
	switch operator {
	case expr.Concat:
		expression = &expr.ConcatExpression{Left: leftResult.Result, Right: rightResult.Result}
	case expr.Add:
		expression = &expr.ArithmeticExpression{Left: leftResult.Result, Right: rightResult.Result, Op: expr.EvaluateAdd}
	case expr.Sub:
		expression = &expr.ArithmeticExpression{Left: leftResult.Result, Right: rightResult.Result, Op: expr.EvaluateSub}
	}
	return v.transformedVisitResult(expression)
}

func (v *FHIRPathVisitor) VisitMultiplicativeExpression(ctx *grammar.MultiplicativeExpressionContext) interface{} {
	leftResult := v.Visit(ctx.Expression(0)).(*VisitResult)
	if leftResult.Error != nil {
		return &VisitResult{nil, leftResult.Error}
	}
	rightResult := v.clone().Visit(ctx.Expression(1)).(*VisitResult)
	if rightResult.Error != nil {
		return &VisitResult{nil, rightResult.Error}
	}

	operator := expr.Operator(ctx.GetChild(1).(antlr.TerminalNode).GetText())

	// Select correct operator function.
	var op func(system.Any, system.Any) (system.Any, error)
	switch operator {
	case expr.Mul:
		op = expr.EvaluateMul
	case expr.Div:
		op = expr.EvaluateDiv
	case expr.FloorDiv:
		op = expr.EvaluateFloorDiv
	case expr.Mod:
		op = expr.EvaluateMod
	}

	return v.transformedVisitResult(
		&expr.ArithmeticExpression{Left: leftResult.Result, Right: rightResult.Result, Op: op},
	)
}

func (v *FHIRPathVisitor) VisitUnionExpression(ctx *grammar.UnionExpressionContext) interface{} {
	leftResult := v.Visit(ctx.Expression(0)).(*VisitResult)
	if leftResult.Error != nil {
		return &VisitResult{nil, leftResult.Error}
	}
	rightResult := v.clone().Visit(ctx.Expression(1)).(*VisitResult)
	if rightResult.Error != nil {
		return &VisitResult{nil, rightResult.Error}
	}

	expression := &expr.UnionExpression{Left: leftResult.Result, Right: rightResult.Result}
	return v.transformedVisitResult(expression)
}

func (v *FHIRPathVisitor) VisitOrExpression(ctx *grammar.OrExpressionContext) interface{} {
	leftResult := v.Visit(ctx.Expression(0)).(*VisitResult)
	if leftResult.Error != nil {
		return &VisitResult{nil, leftResult.Error}
	}
	rightResult := v.clone().Visit(ctx.Expression(1)).(*VisitResult)
	if rightResult.Error != nil {
		return &VisitResult{nil, rightResult.Error}
	}

	operator := expr.Operator(ctx.GetChild(1).(antlr.TerminalNode).GetText())

	expression := &expr.BooleanExpression{Left: leftResult.Result, Right: rightResult.Result, Op: operator}
	return v.transformedVisitResult(expression)
}

func (v *FHIRPathVisitor) VisitAndExpression(ctx *grammar.AndExpressionContext) interface{} {
	leftResult := v.Visit(ctx.Expression(0)).(*VisitResult)
	if leftResult.Error != nil {
		return &VisitResult{nil, leftResult.Error}
	}
	rightResult := v.clone().Visit(ctx.Expression(1)).(*VisitResult)
	if rightResult.Error != nil {
		return &VisitResult{nil, rightResult.Error}
	}

	expression := &expr.BooleanExpression{Left: leftResult.Result, Right: rightResult.Result, Op: expr.And}
	return v.transformedVisitResult(expression)
}

func (v *FHIRPathVisitor) VisitMembershipExpression(ctx *grammar.MembershipExpressionContext) interface{} {
	leftResult := v.Visit(ctx.Expression(0)).(*VisitResult)
	if leftResult.Error != nil {
		return &VisitResult{nil, leftResult.Error}
	}
	rightResult := v.clone().Visit(ctx.Expression(1)).(*VisitResult)
	if rightResult.Error != nil {
		return &VisitResult{nil, rightResult.Error}
	}

	operator := expr.Operator(ctx.GetChild(1).(antlr.TerminalNode).GetText())

	expression := &expr.MembershipExpression{Left: leftResult.Result, Right: rightResult.Result, Operator: operator}
	return v.transformedVisitResult(expression)
}

func (v *FHIRPathVisitor) VisitInequalityExpression(ctx *grammar.InequalityExpressionContext) interface{} {
	leftResult := v.Visit(ctx.Expression(0)).(*VisitResult)
	if leftResult.Error != nil {
		return &VisitResult{nil, leftResult.Error}
	}
	rightResult := v.clone().Visit(ctx.Expression(1)).(*VisitResult)
	if rightResult.Error != nil {
		return &VisitResult{nil, rightResult.Error}
	}

	operator := expr.Operator(ctx.GetChild(1).(antlr.TerminalNode).GetText())

	expression := &expr.ComparisonExpression{Left: leftResult.Result, Right: rightResult.Result, Op: operator}
	return v.transformedVisitResult(expression)
}

// VisitInvocationExpression visits both sides, and constructs an expression sequence.
func (v *FHIRPathVisitor) VisitInvocationExpression(ctx *grammar.InvocationExpressionContext) interface{} {
	// Visit left side with new visitor, raising error if necessary
	leftResult := v.Visit(ctx.Expression()).(*VisitResult)
	if leftResult.Error != nil {
		return &VisitResult{nil, leftResult.Error}
	}

	// Visit right side, raising error if necessary
	rightResult := v.Visit(ctx.Invocation()).(*VisitResult)
	if rightResult.Error != nil {
		return &VisitResult{nil, rightResult.Error}
	}

	// Construct and return ExpressionSequence
	expressions := []expr.Expression{leftResult.Result, rightResult.Result}
	sequence := &expr.ExpressionSequence{Expressions: expressions}
	return v.transformedVisitResult(sequence)
}

// VisitEqualityExpression both equality subexpressions and constructs an Equality Expression
// from the results of each subexpression
func (v *FHIRPathVisitor) VisitEqualityExpression(ctx *grammar.EqualityExpressionContext) interface{} {
	leftResult := v.Visit(ctx.Expression(0)).(*VisitResult)
	if leftResult.Error != nil {
		return &VisitResult{nil, leftResult.Error}
	}

	rightResult := v.clone().Visit(ctx.Expression(1)).(*VisitResult)
	if rightResult.Error != nil {
		return &VisitResult{nil, rightResult.Error}
	}
	operator := ctx.GetChild(1).(antlr.TerminalNode).GetText()
	var expression expr.Expression
	switch operator {
	case expr.Equals:
		expression = &expr.EqualityExpression{Left: leftResult.Result, Right: rightResult.Result}
	case expr.NotEquals:
		expression = &expr.EqualityExpression{Left: leftResult.Result, Right: rightResult.Result, Not: true}
	case expr.Equivalence:
		// TODO: Implement equivalence expressions
	case expr.Inequivalence:
		// TODO: Implement non-equivalence expressions
	}
	return v.transformedVisitResult(expression)
}

func (v *FHIRPathVisitor) VisitImpliesExpression(ctx *grammar.ImpliesExpressionContext) interface{} {
	leftResult := v.Visit(ctx.Expression(0)).(*VisitResult)
	if leftResult.Error != nil {
		return &VisitResult{nil, leftResult.Error}
	}
	rightResult := v.clone().Visit(ctx.Expression(1)).(*VisitResult)
	if rightResult.Error != nil {
		return &VisitResult{nil, rightResult.Error}
	}

	expression := &expr.BooleanExpression{Left: leftResult.Result, Right: rightResult.Result, Op: expr.Implies}
	return v.transformedVisitResult(expression)
}

func (v *FHIRPathVisitor) VisitTermExpression(ctx *grammar.TermExpressionContext) interface{} {
	return v.Visit(ctx.Term())
}

func (v *FHIRPathVisitor) VisitTypeExpression(ctx *grammar.TypeExpressionContext) interface{} {
	expression := v.Visit(ctx.Expression()).(*VisitResult)
	if expression.Error != nil {
		return &VisitResult{nil, expression.Error}
	}
	typeSpecifier := v.Visit(ctx.TypeSpecifier()).(*typeResult)
	if typeSpecifier.err != nil {
		return &VisitResult{nil, typeSpecifier.err}
	}
	operator := ctx.GetChild(1).(antlr.TerminalNode).GetText()
	var typeExpression expr.Expression
	if operator == expr.Is {
		typeExpression = &expr.IsExpression{Expr: expression.Result, Type: typeSpecifier.result}
	}
	if operator == expr.As {
		typeExpression = &expr.AsExpression{Expr: expression.Result, Type: typeSpecifier.result}
	}
	return v.transformedVisitResult(typeExpression)
}

func (v *FHIRPathVisitor) VisitInvocationTerm(ctx *grammar.InvocationTermContext) interface{} {
	return v.Visit(ctx.Invocation())
}

func (v *FHIRPathVisitor) VisitLiteralTerm(ctx *grammar.LiteralTermContext) interface{} {
	return v.Visit(ctx.Literal())
}

func (v *FHIRPathVisitor) VisitExternalConstantTerm(ctx *grammar.ExternalConstantTermContext) interface{} {
	ident := ctx.ExternalConstant().GetText()
	ident = strings.TrimPrefix(ident, "%")
	return v.transformedVisitResult(&expr.ExternalConstantExpression{Identifier: ident})
}

func (v *FHIRPathVisitor) VisitParenthesizedTerm(ctx *grammar.ParenthesizedTermContext) interface{} {
	return v.Visit(ctx.Expression())
}

// VisitNullLiteral returns a NullLiteralExpression, without any error.
func (v *FHIRPathVisitor) VisitNullLiteral(ctx *grammar.NullLiteralContext) interface{} {
	result := &expr.LiteralExpression{}
	return v.transformedVisitResult(result)
}

// VisitBooleanLiteral returns a BooleanLiteralExpression, returning an error if
// there is an error during creation of the Boolean.
func (v *FHIRPathVisitor) VisitBooleanLiteral(ctx *grammar.BooleanLiteralContext) interface{} {
	result, err := system.ParseBoolean(ctx.GetText())
	if err != nil {
		return &VisitResult{nil, err}
	}
	expr := &expr.LiteralExpression{Literal: result}
	return v.transformedVisitResult(expr)
}

// VisitStringLiteral returns a StringLiteralExpression, returning an error if there is an
// error during creation of the String.
func (v *FHIRPathVisitor) VisitStringLiteral(ctx *grammar.StringLiteralContext) interface{} {
	result, err := system.ParseString(ctx.STRING().GetText())
	if err != nil {
		return &VisitResult{nil, err}
	}
	expr := &expr.LiteralExpression{Literal: result}
	return v.transformedVisitResult(expr)
}

// VisitNumberLiteral returns either an integer or decimal, depending on whether or not
// the number contains a decimal. Returns an error if there is an error during creation
// of the number.
func (v *FHIRPathVisitor) VisitNumberLiteral(ctx *grammar.NumberLiteralContext) interface{} {
	number := ctx.NUMBER().GetText()

	if strings.Contains(number, ".") {
		result, err := system.ParseDecimal(number)
		if err != nil {
			return &VisitResult{nil, err}
		}
		expr := &expr.LiteralExpression{Literal: result}
		return v.transformedVisitResult(expr)
	}

	result, err := system.ParseInteger(number)
	if err != nil {
		return &VisitResult{nil, err}
	}
	expr := &expr.LiteralExpression{Literal: result}
	return v.transformedVisitResult(expr)
}

// VisitDateLiteral returns a DateLiteralExpression, returning an error if there is an
// error during creation of the Date type.
func (v *FHIRPathVisitor) VisitDateLiteral(ctx *grammar.DateLiteralContext) interface{} {
	date, err := system.ParseDate(ctx.DATE().GetText())
	if err != nil {
		return &VisitResult{nil, err}
	}
	expr := &expr.LiteralExpression{Literal: date}
	return v.transformedVisitResult(expr)
}

// VisitDateTimeLiteral returns a DateTimeLiteralExpression, returning an error if there
// is an error during creation of the DateTime type.
func (v *FHIRPathVisitor) VisitDateTimeLiteral(ctx *grammar.DateTimeLiteralContext) interface{} {
	dateTime, err := system.ParseDateTime(ctx.DATETIME().GetText())
	if err != nil {
		return &VisitResult{nil, err}
	}
	expr := &expr.LiteralExpression{Literal: dateTime}
	return v.transformedVisitResult(expr)
}

// VisitTimeLiteral returns a TimeLiteralExpression, returning an error if there is an error
// during creation of the Time type.
func (v *FHIRPathVisitor) VisitTimeLiteral(ctx *grammar.TimeLiteralContext) interface{} {
	time, err := system.ParseTime(ctx.TIME().GetText())
	if err != nil {
		return &VisitResult{nil, err}
	}
	expr := &expr.LiteralExpression{Literal: time}
	return v.transformedVisitResult(expr)
}

// VisitQuantityLiteral returns a QuantityLiteralExpression, returning an error if there
// is an error during creation of the Quantity type.
func (v *FHIRPathVisitor) VisitQuantityLiteral(ctx *grammar.QuantityLiteralContext) interface{} {
	// remove string quotes from unit
	unit := ctx.Quantity().Unit().GetText()
	unit = strings.TrimPrefix(unit, "'")
	unit = strings.TrimSuffix(unit, "'")

	quantity, err := system.ParseQuantity(ctx.Quantity().NUMBER().GetText(), unit)
	if err != nil {
		return &VisitResult{nil, err}
	}
	expr := &expr.LiteralExpression{Literal: quantity}
	return v.transformedVisitResult(expr)
}

func (v *FHIRPathVisitor) VisitExternalConstant(ctx *grammar.ExternalConstantContext) interface{} {
	return &VisitResult{nil, errNotSupported}
}

// VisitMemberInvocation checks to see if the identifier corresponds to a resource type and is the
// root of the expression. If so, it will return a TypeExpression. Otherwise, it returns a FieldExpression.
func (v *FHIRPathVisitor) VisitMemberInvocation(ctx *grammar.MemberInvocationContext) interface{} {
	identifier := ctx.GetText()
	var expression expr.Expression

	if resource.IsType(identifier) && !v.visitedRoot {
		expression = &expr.TypeExpression{Type: identifier}
		v.visitedRoot = true
	} else {
		expression = &expr.FieldExpression{FieldName: identifier, Permissive: v.Permissive}
	}

	return v.transformedVisitResult(expression)
}

func (v *FHIRPathVisitor) VisitFunctionInvocation(ctx *grammar.FunctionInvocationContext) interface{} {
	return v.Visit(ctx.Function())
}

func (v *FHIRPathVisitor) VisitThisInvocation(ctx *grammar.ThisInvocationContext) interface{} {
	return &VisitResult{&expr.IdentityExpression{}, nil}
}

func (v *FHIRPathVisitor) VisitIndexInvocation(ctx *grammar.IndexInvocationContext) interface{} {
	return &VisitResult{nil, errNotSupported}
}

func (v *FHIRPathVisitor) VisitTotalInvocation(ctx *grammar.TotalInvocationContext) interface{} {
	return &VisitResult{nil, errNotSupported}
}

func (v *FHIRPathVisitor) VisitFunction(ctx *grammar.FunctionContext) interface{} {
	var fn funcs.Function
	if ctx.Identifier() == nil {
		fn = v.Functions["ofType"]
	} else {
		var ok bool
		fn, ok = v.Functions[ctx.Identifier().GetText()]
		if !ok {
			return &VisitResult{nil, fmt.Errorf("%w: %s", errUnresolvedFunction, ctx.Identifier().GetText())}
		}
	}

	// Handling for type functions
	if fn.IsTypeFunction {
		paramList := ctx.ParamList()
		if paramList == nil || len(paramList.AllExpression()) != 1 {
			return &VisitResult{nil, fmt.Errorf("type function expects exactly one argument")}
		}

		// Visit the first argument as a type specifier
		typeExpr := paramList.Expression(0)
		// This gets the string, e.g. "string" or "FHIR.Patient"
		typeName := typeExpr.GetText()

		typeSpecifier, err := reflection.NewTypeSpecifier(typeName)
		if err != nil {
			return &VisitResult{nil, err}
		}

		return v.transformedVisitResult(
			&expr.FunctionExpression{
				Fn:   fn.Func,
				Args: []expr.Expression{&expr.TypeExpression{Type: typeSpecifier.String()}},
			},
		)
	}

	results := []*VisitResult{}
	if args := ctx.ParamList(); args != nil {
		results = v.Visit(args).([]*VisitResult)
	}

	errs := slices.Map(results, func(r *VisitResult) error { return r.Error })
	if err := errors.Join(errs...); err != nil {
		return &VisitResult{nil, fmt.Errorf("%w: %w", errVisitingChildren, err)}
	}

	expressions := slices.Map(results, func(r *VisitResult) expr.Expression { return r.Result })
	if len(expressions) < fn.MinArity || len(expressions) > fn.MaxArity {
		return &VisitResult{nil, fmt.Errorf("%w: input arity outside of function arity bounds", impl.ErrWrongArity)}
	}
	return v.transformedVisitResult(&expr.FunctionExpression{Fn: fn.Func, Args: expressions})
}

func (v *FHIRPathVisitor) VisitParamList(ctx *grammar.ParamListContext) interface{} {
	return slices.Map(ctx.AllExpression(), func(e grammar.IExpressionContext) *VisitResult { return v.Visit(e).(*VisitResult) })
}

func (v *FHIRPathVisitor) VisitQuantity(ctx *grammar.QuantityContext) interface{} {
	return &VisitResult{nil, errNotSupported}
}

func (v *FHIRPathVisitor) VisitUnit(ctx *grammar.UnitContext) interface{} {
	return &VisitResult{nil, errNotSupported}
}

func (v *FHIRPathVisitor) VisitDateTimePrecision(ctx *grammar.DateTimePrecisionContext) interface{} {
	return &VisitResult{nil, errNotSupported}
}

func (v *FHIRPathVisitor) VisitPluralDateTimePrecision(ctx *grammar.PluralDateTimePrecisionContext) interface{} {
	return &VisitResult{nil, errNotSupported}
}

func (v *FHIRPathVisitor) VisitTypeSpecifier(ctx *grammar.TypeSpecifierContext) interface{} {
	identifiers := v.Visit(ctx.QualifiedIdentifier()).([]string)
	if len(identifiers) == 1 {
		specifier, err := reflection.NewTypeSpecifier(identifiers[0])
		return &typeResult{specifier, err}
	}
	if len(identifiers) == 2 {
		specifier, err := reflection.NewQualifiedTypeSpecifier(identifiers[0], identifiers[1])
		return &typeResult{specifier, err}
	}
	return &typeResult{err: fmt.Errorf("%w: %s", errTooManyQualifiers, strings.Join(identifiers, ","))}
}

func (v *FHIRPathVisitor) VisitQualifiedIdentifier(ctx *grammar.QualifiedIdentifierContext) interface{} {
	return slices.Map(ctx.AllIdentifier(), func(i grammar.IIdentifierContext) string { return i.GetText() })
}

func (v *FHIRPathVisitor) VisitIdentifier(ctx *grammar.IdentifierContext) interface{} {
	return &VisitResult{nil, errNotSupported}
}
