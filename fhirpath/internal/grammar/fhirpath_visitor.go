// Code generated from fhirpath.g4 by ANTLR 4.13.0. DO NOT EDIT.

package grammar // fhirpath
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by fhirpathParser.
type fhirpathVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by fhirpathParser#prog.
	VisitProg(ctx *ProgContext) interface{}

	// Visit a parse tree produced by fhirpathParser#indexerExpression.
	VisitIndexerExpression(ctx *IndexerExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#polarityExpression.
	VisitPolarityExpression(ctx *PolarityExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#additiveExpression.
	VisitAdditiveExpression(ctx *AdditiveExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#multiplicativeExpression.
	VisitMultiplicativeExpression(ctx *MultiplicativeExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#unionExpression.
	VisitUnionExpression(ctx *UnionExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#orExpression.
	VisitOrExpression(ctx *OrExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#andExpression.
	VisitAndExpression(ctx *AndExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#membershipExpression.
	VisitMembershipExpression(ctx *MembershipExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#inequalityExpression.
	VisitInequalityExpression(ctx *InequalityExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#invocationExpression.
	VisitInvocationExpression(ctx *InvocationExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#equalityExpression.
	VisitEqualityExpression(ctx *EqualityExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#impliesExpression.
	VisitImpliesExpression(ctx *ImpliesExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#termExpression.
	VisitTermExpression(ctx *TermExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#typeExpression.
	VisitTypeExpression(ctx *TypeExpressionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#invocationTerm.
	VisitInvocationTerm(ctx *InvocationTermContext) interface{}

	// Visit a parse tree produced by fhirpathParser#literalTerm.
	VisitLiteralTerm(ctx *LiteralTermContext) interface{}

	// Visit a parse tree produced by fhirpathParser#externalConstantTerm.
	VisitExternalConstantTerm(ctx *ExternalConstantTermContext) interface{}

	// Visit a parse tree produced by fhirpathParser#parenthesizedTerm.
	VisitParenthesizedTerm(ctx *ParenthesizedTermContext) interface{}

	// Visit a parse tree produced by fhirpathParser#nullLiteral.
	VisitNullLiteral(ctx *NullLiteralContext) interface{}

	// Visit a parse tree produced by fhirpathParser#booleanLiteral.
	VisitBooleanLiteral(ctx *BooleanLiteralContext) interface{}

	// Visit a parse tree produced by fhirpathParser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by fhirpathParser#numberLiteral.
	VisitNumberLiteral(ctx *NumberLiteralContext) interface{}

	// Visit a parse tree produced by fhirpathParser#dateLiteral.
	VisitDateLiteral(ctx *DateLiteralContext) interface{}

	// Visit a parse tree produced by fhirpathParser#dateTimeLiteral.
	VisitDateTimeLiteral(ctx *DateTimeLiteralContext) interface{}

	// Visit a parse tree produced by fhirpathParser#timeLiteral.
	VisitTimeLiteral(ctx *TimeLiteralContext) interface{}

	// Visit a parse tree produced by fhirpathParser#quantityLiteral.
	VisitQuantityLiteral(ctx *QuantityLiteralContext) interface{}

	// Visit a parse tree produced by fhirpathParser#externalConstant.
	VisitExternalConstant(ctx *ExternalConstantContext) interface{}

	// Visit a parse tree produced by fhirpathParser#memberInvocation.
	VisitMemberInvocation(ctx *MemberInvocationContext) interface{}

	// Visit a parse tree produced by fhirpathParser#functionInvocation.
	VisitFunctionInvocation(ctx *FunctionInvocationContext) interface{}

	// Visit a parse tree produced by fhirpathParser#thisInvocation.
	VisitThisInvocation(ctx *ThisInvocationContext) interface{}

	// Visit a parse tree produced by fhirpathParser#indexInvocation.
	VisitIndexInvocation(ctx *IndexInvocationContext) interface{}

	// Visit a parse tree produced by fhirpathParser#totalInvocation.
	VisitTotalInvocation(ctx *TotalInvocationContext) interface{}

	// Visit a parse tree produced by fhirpathParser#function.
	VisitFunction(ctx *FunctionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#paramList.
	VisitParamList(ctx *ParamListContext) interface{}

	// Visit a parse tree produced by fhirpathParser#quantity.
	VisitQuantity(ctx *QuantityContext) interface{}

	// Visit a parse tree produced by fhirpathParser#unit.
	VisitUnit(ctx *UnitContext) interface{}

	// Visit a parse tree produced by fhirpathParser#dateTimePrecision.
	VisitDateTimePrecision(ctx *DateTimePrecisionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#pluralDateTimePrecision.
	VisitPluralDateTimePrecision(ctx *PluralDateTimePrecisionContext) interface{}

	// Visit a parse tree produced by fhirpathParser#typeSpecifier.
	VisitTypeSpecifier(ctx *TypeSpecifierContext) interface{}

	// Visit a parse tree produced by fhirpathParser#qualifiedIdentifier.
	VisitQualifiedIdentifier(ctx *QualifiedIdentifierContext) interface{}

	// Visit a parse tree produced by fhirpathParser#identifier.
	VisitIdentifier(ctx *IdentifierContext) interface{}
}
