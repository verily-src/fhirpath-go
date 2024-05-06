// Code generated from fhirpath.g4 by ANTLR 4.13.0. DO NOT EDIT.

package grammar // fhirpath
import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type fhirpathParser struct {
	*antlr.BaseParser
}

var FhirpathParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func fhirpathParserInit() {
	staticData := &FhirpathParserStaticData
	staticData.LiteralNames = []string{
		"", "'.'", "'['", "']'", "'+'", "'-'", "'*'", "'/'", "'div'", "'mod'",
		"'&'", "'is'", "'as'", "'|'", "'<='", "'<'", "'>'", "'>='", "'='", "'~'",
		"'!='", "'!~'", "'in'", "'contains'", "'and'", "'or'", "'xor'", "'implies'",
		"'('", "')'", "'{'", "'}'", "'true'", "'false'", "'%'", "'$this'", "'$index'",
		"'$total'", "','", "'year'", "'month'", "'week'", "'day'", "'hour'",
		"'minute'", "'second'", "'millisecond'", "'years'", "'months'", "'weeks'",
		"'days'", "'hours'", "'minutes'", "'seconds'", "'milliseconds'",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "DATE", "DATETIME", "TIME", "IDENTIFIER", "DELIMITEDIDENTIFIER",
		"STRING", "NUMBER", "WS", "COMMENT", "LINE_COMMENT",
	}
	staticData.RuleNames = []string{
		"prog", "expression", "term", "literal", "externalConstant", "invocation",
		"function", "paramList", "quantity", "unit", "dateTimePrecision", "pluralDateTimePrecision",
		"typeSpecifier", "qualifiedIdentifier", "identifier",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 64, 155, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 1, 0, 1, 0,
		1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 38, 8, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 5, 1, 78, 8, 1,
		10, 1, 12, 1, 81, 9, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3, 2,
		90, 8, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 3, 3, 101,
		8, 3, 1, 4, 1, 4, 1, 4, 3, 4, 106, 8, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5,
		3, 5, 113, 8, 5, 1, 6, 1, 6, 1, 6, 3, 6, 118, 8, 6, 1, 6, 1, 6, 1, 7, 1,
		7, 1, 7, 5, 7, 125, 8, 7, 10, 7, 12, 7, 128, 9, 7, 1, 8, 1, 8, 3, 8, 132,
		8, 8, 1, 9, 1, 9, 1, 9, 3, 9, 137, 8, 9, 1, 10, 1, 10, 1, 11, 1, 11, 1,
		12, 1, 12, 1, 13, 1, 13, 1, 13, 5, 13, 148, 8, 13, 10, 13, 12, 13, 151,
		9, 13, 1, 14, 1, 14, 1, 14, 0, 1, 2, 15, 0, 2, 4, 6, 8, 10, 12, 14, 16,
		18, 20, 22, 24, 26, 28, 0, 12, 1, 0, 4, 5, 1, 0, 6, 9, 2, 0, 4, 5, 10,
		10, 1, 0, 14, 17, 1, 0, 18, 21, 1, 0, 22, 23, 1, 0, 25, 26, 1, 0, 11, 12,
		1, 0, 32, 33, 1, 0, 39, 46, 1, 0, 47, 54, 3, 0, 11, 12, 22, 23, 58, 59,
		173, 0, 30, 1, 0, 0, 0, 2, 37, 1, 0, 0, 0, 4, 89, 1, 0, 0, 0, 6, 100, 1,
		0, 0, 0, 8, 102, 1, 0, 0, 0, 10, 112, 1, 0, 0, 0, 12, 114, 1, 0, 0, 0,
		14, 121, 1, 0, 0, 0, 16, 129, 1, 0, 0, 0, 18, 136, 1, 0, 0, 0, 20, 138,
		1, 0, 0, 0, 22, 140, 1, 0, 0, 0, 24, 142, 1, 0, 0, 0, 26, 144, 1, 0, 0,
		0, 28, 152, 1, 0, 0, 0, 30, 31, 3, 2, 1, 0, 31, 32, 5, 0, 0, 1, 32, 1,
		1, 0, 0, 0, 33, 34, 6, 1, -1, 0, 34, 38, 3, 4, 2, 0, 35, 36, 7, 0, 0, 0,
		36, 38, 3, 2, 1, 11, 37, 33, 1, 0, 0, 0, 37, 35, 1, 0, 0, 0, 38, 79, 1,
		0, 0, 0, 39, 40, 10, 10, 0, 0, 40, 41, 7, 1, 0, 0, 41, 78, 3, 2, 1, 11,
		42, 43, 10, 9, 0, 0, 43, 44, 7, 2, 0, 0, 44, 78, 3, 2, 1, 10, 45, 46, 10,
		7, 0, 0, 46, 47, 5, 13, 0, 0, 47, 78, 3, 2, 1, 8, 48, 49, 10, 6, 0, 0,
		49, 50, 7, 3, 0, 0, 50, 78, 3, 2, 1, 7, 51, 52, 10, 5, 0, 0, 52, 53, 7,
		4, 0, 0, 53, 78, 3, 2, 1, 6, 54, 55, 10, 4, 0, 0, 55, 56, 7, 5, 0, 0, 56,
		78, 3, 2, 1, 5, 57, 58, 10, 3, 0, 0, 58, 59, 5, 24, 0, 0, 59, 78, 3, 2,
		1, 4, 60, 61, 10, 2, 0, 0, 61, 62, 7, 6, 0, 0, 62, 78, 3, 2, 1, 3, 63,
		64, 10, 1, 0, 0, 64, 65, 5, 27, 0, 0, 65, 78, 3, 2, 1, 2, 66, 67, 10, 13,
		0, 0, 67, 68, 5, 1, 0, 0, 68, 78, 3, 10, 5, 0, 69, 70, 10, 12, 0, 0, 70,
		71, 5, 2, 0, 0, 71, 72, 3, 2, 1, 0, 72, 73, 5, 3, 0, 0, 73, 78, 1, 0, 0,
		0, 74, 75, 10, 8, 0, 0, 75, 76, 7, 7, 0, 0, 76, 78, 3, 24, 12, 0, 77, 39,
		1, 0, 0, 0, 77, 42, 1, 0, 0, 0, 77, 45, 1, 0, 0, 0, 77, 48, 1, 0, 0, 0,
		77, 51, 1, 0, 0, 0, 77, 54, 1, 0, 0, 0, 77, 57, 1, 0, 0, 0, 77, 60, 1,
		0, 0, 0, 77, 63, 1, 0, 0, 0, 77, 66, 1, 0, 0, 0, 77, 69, 1, 0, 0, 0, 77,
		74, 1, 0, 0, 0, 78, 81, 1, 0, 0, 0, 79, 77, 1, 0, 0, 0, 79, 80, 1, 0, 0,
		0, 80, 3, 1, 0, 0, 0, 81, 79, 1, 0, 0, 0, 82, 90, 3, 10, 5, 0, 83, 90,
		3, 6, 3, 0, 84, 90, 3, 8, 4, 0, 85, 86, 5, 28, 0, 0, 86, 87, 3, 2, 1, 0,
		87, 88, 5, 29, 0, 0, 88, 90, 1, 0, 0, 0, 89, 82, 1, 0, 0, 0, 89, 83, 1,
		0, 0, 0, 89, 84, 1, 0, 0, 0, 89, 85, 1, 0, 0, 0, 90, 5, 1, 0, 0, 0, 91,
		92, 5, 30, 0, 0, 92, 101, 5, 31, 0, 0, 93, 101, 7, 8, 0, 0, 94, 101, 5,
		60, 0, 0, 95, 101, 5, 61, 0, 0, 96, 101, 5, 55, 0, 0, 97, 101, 5, 56, 0,
		0, 98, 101, 5, 57, 0, 0, 99, 101, 3, 16, 8, 0, 100, 91, 1, 0, 0, 0, 100,
		93, 1, 0, 0, 0, 100, 94, 1, 0, 0, 0, 100, 95, 1, 0, 0, 0, 100, 96, 1, 0,
		0, 0, 100, 97, 1, 0, 0, 0, 100, 98, 1, 0, 0, 0, 100, 99, 1, 0, 0, 0, 101,
		7, 1, 0, 0, 0, 102, 105, 5, 34, 0, 0, 103, 106, 3, 28, 14, 0, 104, 106,
		5, 60, 0, 0, 105, 103, 1, 0, 0, 0, 105, 104, 1, 0, 0, 0, 106, 9, 1, 0,
		0, 0, 107, 113, 3, 28, 14, 0, 108, 113, 3, 12, 6, 0, 109, 113, 5, 35, 0,
		0, 110, 113, 5, 36, 0, 0, 111, 113, 5, 37, 0, 0, 112, 107, 1, 0, 0, 0,
		112, 108, 1, 0, 0, 0, 112, 109, 1, 0, 0, 0, 112, 110, 1, 0, 0, 0, 112,
		111, 1, 0, 0, 0, 113, 11, 1, 0, 0, 0, 114, 115, 3, 28, 14, 0, 115, 117,
		5, 28, 0, 0, 116, 118, 3, 14, 7, 0, 117, 116, 1, 0, 0, 0, 117, 118, 1,
		0, 0, 0, 118, 119, 1, 0, 0, 0, 119, 120, 5, 29, 0, 0, 120, 13, 1, 0, 0,
		0, 121, 126, 3, 2, 1, 0, 122, 123, 5, 38, 0, 0, 123, 125, 3, 2, 1, 0, 124,
		122, 1, 0, 0, 0, 125, 128, 1, 0, 0, 0, 126, 124, 1, 0, 0, 0, 126, 127,
		1, 0, 0, 0, 127, 15, 1, 0, 0, 0, 128, 126, 1, 0, 0, 0, 129, 131, 5, 61,
		0, 0, 130, 132, 3, 18, 9, 0, 131, 130, 1, 0, 0, 0, 131, 132, 1, 0, 0, 0,
		132, 17, 1, 0, 0, 0, 133, 137, 3, 20, 10, 0, 134, 137, 3, 22, 11, 0, 135,
		137, 5, 60, 0, 0, 136, 133, 1, 0, 0, 0, 136, 134, 1, 0, 0, 0, 136, 135,
		1, 0, 0, 0, 137, 19, 1, 0, 0, 0, 138, 139, 7, 9, 0, 0, 139, 21, 1, 0, 0,
		0, 140, 141, 7, 10, 0, 0, 141, 23, 1, 0, 0, 0, 142, 143, 3, 26, 13, 0,
		143, 25, 1, 0, 0, 0, 144, 149, 3, 28, 14, 0, 145, 146, 5, 1, 0, 0, 146,
		148, 3, 28, 14, 0, 147, 145, 1, 0, 0, 0, 148, 151, 1, 0, 0, 0, 149, 147,
		1, 0, 0, 0, 149, 150, 1, 0, 0, 0, 150, 27, 1, 0, 0, 0, 151, 149, 1, 0,
		0, 0, 152, 153, 7, 11, 0, 0, 153, 29, 1, 0, 0, 0, 12, 37, 77, 79, 89, 100,
		105, 112, 117, 126, 131, 136, 149,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// fhirpathParserInit initializes any static state used to implement fhirpathParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewfhirpathParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func FhirpathParserInit() {
	staticData := &FhirpathParserStaticData
	staticData.once.Do(fhirpathParserInit)
}

// NewfhirpathParser produces a new parser instance for the optional input antlr.TokenStream.
func NewfhirpathParser(input antlr.TokenStream) *fhirpathParser {
	FhirpathParserInit()
	this := new(fhirpathParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &FhirpathParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "fhirpath.g4"

	return this
}

// fhirpathParser tokens.
const (
	fhirpathParserEOF                 = antlr.TokenEOF
	fhirpathParserT__0                = 1
	fhirpathParserT__1                = 2
	fhirpathParserT__2                = 3
	fhirpathParserT__3                = 4
	fhirpathParserT__4                = 5
	fhirpathParserT__5                = 6
	fhirpathParserT__6                = 7
	fhirpathParserT__7                = 8
	fhirpathParserT__8                = 9
	fhirpathParserT__9                = 10
	fhirpathParserT__10               = 11
	fhirpathParserT__11               = 12
	fhirpathParserT__12               = 13
	fhirpathParserT__13               = 14
	fhirpathParserT__14               = 15
	fhirpathParserT__15               = 16
	fhirpathParserT__16               = 17
	fhirpathParserT__17               = 18
	fhirpathParserT__18               = 19
	fhirpathParserT__19               = 20
	fhirpathParserT__20               = 21
	fhirpathParserT__21               = 22
	fhirpathParserT__22               = 23
	fhirpathParserT__23               = 24
	fhirpathParserT__24               = 25
	fhirpathParserT__25               = 26
	fhirpathParserT__26               = 27
	fhirpathParserT__27               = 28
	fhirpathParserT__28               = 29
	fhirpathParserT__29               = 30
	fhirpathParserT__30               = 31
	fhirpathParserT__31               = 32
	fhirpathParserT__32               = 33
	fhirpathParserT__33               = 34
	fhirpathParserT__34               = 35
	fhirpathParserT__35               = 36
	fhirpathParserT__36               = 37
	fhirpathParserT__37               = 38
	fhirpathParserT__38               = 39
	fhirpathParserT__39               = 40
	fhirpathParserT__40               = 41
	fhirpathParserT__41               = 42
	fhirpathParserT__42               = 43
	fhirpathParserT__43               = 44
	fhirpathParserT__44               = 45
	fhirpathParserT__45               = 46
	fhirpathParserT__46               = 47
	fhirpathParserT__47               = 48
	fhirpathParserT__48               = 49
	fhirpathParserT__49               = 50
	fhirpathParserT__50               = 51
	fhirpathParserT__51               = 52
	fhirpathParserT__52               = 53
	fhirpathParserT__53               = 54
	fhirpathParserDATE                = 55
	fhirpathParserDATETIME            = 56
	fhirpathParserTIME                = 57
	fhirpathParserIDENTIFIER          = 58
	fhirpathParserDELIMITEDIDENTIFIER = 59
	fhirpathParserSTRING              = 60
	fhirpathParserNUMBER              = 61
	fhirpathParserWS                  = 62
	fhirpathParserCOMMENT             = 63
	fhirpathParserLINE_COMMENT        = 64
)

// fhirpathParser rules.
const (
	fhirpathParserRULE_prog                    = 0
	fhirpathParserRULE_expression              = 1
	fhirpathParserRULE_term                    = 2
	fhirpathParserRULE_literal                 = 3
	fhirpathParserRULE_externalConstant        = 4
	fhirpathParserRULE_invocation              = 5
	fhirpathParserRULE_function                = 6
	fhirpathParserRULE_paramList               = 7
	fhirpathParserRULE_quantity                = 8
	fhirpathParserRULE_unit                    = 9
	fhirpathParserRULE_dateTimePrecision       = 10
	fhirpathParserRULE_pluralDateTimePrecision = 11
	fhirpathParserRULE_typeSpecifier           = 12
	fhirpathParserRULE_qualifiedIdentifier     = 13
	fhirpathParserRULE_identifier              = 14
)

// IProgContext is an interface to support dynamic dispatch.
type IProgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	EOF() antlr.TerminalNode

	// IsProgContext differentiates from other interfaces.
	IsProgContext()
}

type ProgContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgContext() *ProgContext {
	var p = new(ProgContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_prog
	return p
}

func InitEmptyProgContext(p *ProgContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_prog
}

func (*ProgContext) IsProgContext() {}

func NewProgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgContext {
	var p = new(ProgContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_prog

	return p
}

func (s *ProgContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ProgContext) EOF() antlr.TerminalNode {
	return s.GetToken(fhirpathParserEOF, 0)
}

func (s *ProgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitProg(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) Prog() (localctx IProgContext) {
	localctx = NewProgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, fhirpathParserRULE_prog)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(30)
		p.expression(0)
	}
	{
		p.SetState(31)
		p.Match(fhirpathParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) CopyAll(ctx *ExpressionContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type IndexerExpressionContext struct {
	ExpressionContext
}

func NewIndexerExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IndexerExpressionContext {
	var p = new(IndexerExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *IndexerExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IndexerExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *IndexerExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *IndexerExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitIndexerExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type PolarityExpressionContext struct {
	ExpressionContext
}

func NewPolarityExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PolarityExpressionContext {
	var p = new(PolarityExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *PolarityExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PolarityExpressionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *PolarityExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitPolarityExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type AdditiveExpressionContext struct {
	ExpressionContext
}

func NewAdditiveExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AdditiveExpressionContext {
	var p = new(AdditiveExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *AdditiveExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AdditiveExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *AdditiveExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *AdditiveExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitAdditiveExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type MultiplicativeExpressionContext struct {
	ExpressionContext
}

func NewMultiplicativeExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MultiplicativeExpressionContext {
	var p = new(MultiplicativeExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *MultiplicativeExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MultiplicativeExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *MultiplicativeExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *MultiplicativeExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitMultiplicativeExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type UnionExpressionContext struct {
	ExpressionContext
}

func NewUnionExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UnionExpressionContext {
	var p = new(UnionExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *UnionExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnionExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *UnionExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *UnionExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitUnionExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type OrExpressionContext struct {
	ExpressionContext
}

func NewOrExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *OrExpressionContext {
	var p = new(OrExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *OrExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *OrExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *OrExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitOrExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type AndExpressionContext struct {
	ExpressionContext
}

func NewAndExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *AndExpressionContext {
	var p = new(AndExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *AndExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AndExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *AndExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *AndExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitAndExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type MembershipExpressionContext struct {
	ExpressionContext
}

func NewMembershipExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MembershipExpressionContext {
	var p = new(MembershipExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *MembershipExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MembershipExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *MembershipExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *MembershipExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitMembershipExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type InequalityExpressionContext struct {
	ExpressionContext
}

func NewInequalityExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InequalityExpressionContext {
	var p = new(InequalityExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *InequalityExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InequalityExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *InequalityExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *InequalityExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitInequalityExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type InvocationExpressionContext struct {
	ExpressionContext
}

func NewInvocationExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InvocationExpressionContext {
	var p = new(InvocationExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *InvocationExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InvocationExpressionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *InvocationExpressionContext) Invocation() IInvocationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInvocationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInvocationContext)
}

func (s *InvocationExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitInvocationExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type EqualityExpressionContext struct {
	ExpressionContext
}

func NewEqualityExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EqualityExpressionContext {
	var p = new(EqualityExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *EqualityExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *EqualityExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *EqualityExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *EqualityExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitEqualityExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type ImpliesExpressionContext struct {
	ExpressionContext
}

func NewImpliesExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ImpliesExpressionContext {
	var p = new(ImpliesExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *ImpliesExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ImpliesExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ImpliesExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ImpliesExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitImpliesExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type TermExpressionContext struct {
	ExpressionContext
}

func NewTermExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TermExpressionContext {
	var p = new(TermExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *TermExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TermExpressionContext) Term() ITermContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITermContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITermContext)
}

func (s *TermExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitTermExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type TypeExpressionContext struct {
	ExpressionContext
}

func NewTypeExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TypeExpressionContext {
	var p = new(TypeExpressionContext)

	InitEmptyExpressionContext(&p.ExpressionContext)
	p.parser = parser
	p.CopyAll(ctx.(*ExpressionContext))

	return p
}

func (s *TypeExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeExpressionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *TypeExpressionContext) TypeSpecifier() ITypeSpecifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeSpecifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeSpecifierContext)
}

func (s *TypeExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitTypeExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *fhirpathParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 2
	p.EnterRecursionRule(localctx, 2, fhirpathParserRULE_expression, _p)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(37)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case fhirpathParserT__10, fhirpathParserT__11, fhirpathParserT__21, fhirpathParserT__22, fhirpathParserT__27, fhirpathParserT__29, fhirpathParserT__31, fhirpathParserT__32, fhirpathParserT__33, fhirpathParserT__34, fhirpathParserT__35, fhirpathParserT__36, fhirpathParserDATE, fhirpathParserDATETIME, fhirpathParserTIME, fhirpathParserIDENTIFIER, fhirpathParserDELIMITEDIDENTIFIER, fhirpathParserSTRING, fhirpathParserNUMBER:
		localctx = NewTermExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(34)
			p.Term()
		}

	case fhirpathParserT__3, fhirpathParserT__4:
		localctx = NewPolarityExpressionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(35)
			_la = p.GetTokenStream().LA(1)

			if !(_la == fhirpathParserT__3 || _la == fhirpathParserT__4) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(36)
			p.expression(11)
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(79)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(77)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
			case 1:
				localctx = NewMultiplicativeExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, fhirpathParserRULE_expression)
				p.SetState(39)

				if !(p.Precpred(p.GetParserRuleContext(), 10)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 10)", ""))
					goto errorExit
				}
				{
					p.SetState(40)
					_la = p.GetTokenStream().LA(1)

					if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&960) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(41)
					p.expression(11)
				}

			case 2:
				localctx = NewAdditiveExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, fhirpathParserRULE_expression)
				p.SetState(42)

				if !(p.Precpred(p.GetParserRuleContext(), 9)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 9)", ""))
					goto errorExit
				}
				{
					p.SetState(43)
					_la = p.GetTokenStream().LA(1)

					if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1072) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(44)
					p.expression(10)
				}

			case 3:
				localctx = NewUnionExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, fhirpathParserRULE_expression)
				p.SetState(45)

				if !(p.Precpred(p.GetParserRuleContext(), 7)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 7)", ""))
					goto errorExit
				}
				{
					p.SetState(46)
					p.Match(fhirpathParserT__12)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(47)
					p.expression(8)
				}

			case 4:
				localctx = NewInequalityExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, fhirpathParserRULE_expression)
				p.SetState(48)

				if !(p.Precpred(p.GetParserRuleContext(), 6)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 6)", ""))
					goto errorExit
				}
				{
					p.SetState(49)
					_la = p.GetTokenStream().LA(1)

					if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&245760) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(50)
					p.expression(7)
				}

			case 5:
				localctx = NewEqualityExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, fhirpathParserRULE_expression)
				p.SetState(51)

				if !(p.Precpred(p.GetParserRuleContext(), 5)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 5)", ""))
					goto errorExit
				}
				{
					p.SetState(52)
					_la = p.GetTokenStream().LA(1)

					if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&3932160) != 0) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(53)
					p.expression(6)
				}

			case 6:
				localctx = NewMembershipExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, fhirpathParserRULE_expression)
				p.SetState(54)

				if !(p.Precpred(p.GetParserRuleContext(), 4)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 4)", ""))
					goto errorExit
				}
				{
					p.SetState(55)
					_la = p.GetTokenStream().LA(1)

					if !(_la == fhirpathParserT__21 || _la == fhirpathParserT__22) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(56)
					p.expression(5)
				}

			case 7:
				localctx = NewAndExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, fhirpathParserRULE_expression)
				p.SetState(57)

				if !(p.Precpred(p.GetParserRuleContext(), 3)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 3)", ""))
					goto errorExit
				}
				{
					p.SetState(58)
					p.Match(fhirpathParserT__23)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(59)
					p.expression(4)
				}

			case 8:
				localctx = NewOrExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, fhirpathParserRULE_expression)
				p.SetState(60)

				if !(p.Precpred(p.GetParserRuleContext(), 2)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
					goto errorExit
				}
				{
					p.SetState(61)
					_la = p.GetTokenStream().LA(1)

					if !(_la == fhirpathParserT__24 || _la == fhirpathParserT__25) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(62)
					p.expression(3)
				}

			case 9:
				localctx = NewImpliesExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, fhirpathParserRULE_expression)
				p.SetState(63)

				if !(p.Precpred(p.GetParserRuleContext(), 1)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
					goto errorExit
				}
				{
					p.SetState(64)
					p.Match(fhirpathParserT__26)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(65)
					p.expression(2)
				}

			case 10:
				localctx = NewInvocationExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, fhirpathParserRULE_expression)
				p.SetState(66)

				if !(p.Precpred(p.GetParserRuleContext(), 13)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 13)", ""))
					goto errorExit
				}
				{
					p.SetState(67)
					p.Match(fhirpathParserT__0)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(68)
					p.Invocation()
				}

			case 11:
				localctx = NewIndexerExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, fhirpathParserRULE_expression)
				p.SetState(69)

				if !(p.Precpred(p.GetParserRuleContext(), 12)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 12)", ""))
					goto errorExit
				}
				{
					p.SetState(70)
					p.Match(fhirpathParserT__1)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(71)
					p.expression(0)
				}
				{
					p.SetState(72)
					p.Match(fhirpathParserT__2)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}

			case 12:
				localctx = NewTypeExpressionContext(p, NewExpressionContext(p, _parentctx, _parentState))
				p.PushNewRecursionContext(localctx, _startState, fhirpathParserRULE_expression)
				p.SetState(74)

				if !(p.Precpred(p.GetParserRuleContext(), 8)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 8)", ""))
					goto errorExit
				}
				{
					p.SetState(75)
					_la = p.GetTokenStream().LA(1)

					if !(_la == fhirpathParserT__10 || _la == fhirpathParserT__11) {
						p.GetErrorHandler().RecoverInline(p)
					} else {
						p.GetErrorHandler().ReportMatch(p)
						p.Consume()
					}
				}
				{
					p.SetState(76)
					p.TypeSpecifier()
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(81)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITermContext is an interface to support dynamic dispatch.
type ITermContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsTermContext differentiates from other interfaces.
	IsTermContext()
}

type TermContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTermContext() *TermContext {
	var p = new(TermContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_term
	return p
}

func InitEmptyTermContext(p *TermContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_term
}

func (*TermContext) IsTermContext() {}

func NewTermContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TermContext {
	var p = new(TermContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_term

	return p
}

func (s *TermContext) GetParser() antlr.Parser { return s.parser }

func (s *TermContext) CopyAll(ctx *TermContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *TermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TermContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ExternalConstantTermContext struct {
	TermContext
}

func NewExternalConstantTermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ExternalConstantTermContext {
	var p = new(ExternalConstantTermContext)

	InitEmptyTermContext(&p.TermContext)
	p.parser = parser
	p.CopyAll(ctx.(*TermContext))

	return p
}

func (s *ExternalConstantTermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExternalConstantTermContext) ExternalConstant() IExternalConstantContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExternalConstantContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExternalConstantContext)
}

func (s *ExternalConstantTermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitExternalConstantTerm(s)

	default:
		return t.VisitChildren(s)
	}
}

type LiteralTermContext struct {
	TermContext
}

func NewLiteralTermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LiteralTermContext {
	var p = new(LiteralTermContext)

	InitEmptyTermContext(&p.TermContext)
	p.parser = parser
	p.CopyAll(ctx.(*TermContext))

	return p
}

func (s *LiteralTermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralTermContext) Literal() ILiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *LiteralTermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitLiteralTerm(s)

	default:
		return t.VisitChildren(s)
	}
}

type ParenthesizedTermContext struct {
	TermContext
}

func NewParenthesizedTermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ParenthesizedTermContext {
	var p = new(ParenthesizedTermContext)

	InitEmptyTermContext(&p.TermContext)
	p.parser = parser
	p.CopyAll(ctx.(*TermContext))

	return p
}

func (s *ParenthesizedTermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParenthesizedTermContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ParenthesizedTermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitParenthesizedTerm(s)

	default:
		return t.VisitChildren(s)
	}
}

type InvocationTermContext struct {
	TermContext
}

func NewInvocationTermContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *InvocationTermContext {
	var p = new(InvocationTermContext)

	InitEmptyTermContext(&p.TermContext)
	p.parser = parser
	p.CopyAll(ctx.(*TermContext))

	return p
}

func (s *InvocationTermContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InvocationTermContext) Invocation() IInvocationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInvocationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInvocationContext)
}

func (s *InvocationTermContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitInvocationTerm(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) Term() (localctx ITermContext) {
	localctx = NewTermContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, fhirpathParserRULE_term)
	p.SetState(89)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case fhirpathParserT__10, fhirpathParserT__11, fhirpathParserT__21, fhirpathParserT__22, fhirpathParserT__34, fhirpathParserT__35, fhirpathParserT__36, fhirpathParserIDENTIFIER, fhirpathParserDELIMITEDIDENTIFIER:
		localctx = NewInvocationTermContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(82)
			p.Invocation()
		}

	case fhirpathParserT__29, fhirpathParserT__31, fhirpathParserT__32, fhirpathParserDATE, fhirpathParserDATETIME, fhirpathParserTIME, fhirpathParserSTRING, fhirpathParserNUMBER:
		localctx = NewLiteralTermContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(83)
			p.Literal()
		}

	case fhirpathParserT__33:
		localctx = NewExternalConstantTermContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(84)
			p.ExternalConstant()
		}

	case fhirpathParserT__27:
		localctx = NewParenthesizedTermContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(85)
			p.Match(fhirpathParserT__27)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(86)
			p.expression(0)
		}
		{
			p.SetState(87)
			p.Match(fhirpathParserT__28)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILiteralContext is an interface to support dynamic dispatch.
type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}

type LiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_literal
	return p
}

func InitEmptyLiteralContext(p *LiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_literal
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) CopyAll(ctx *LiteralContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type TimeLiteralContext struct {
	LiteralContext
}

func NewTimeLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TimeLiteralContext {
	var p = new(TimeLiteralContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

	return p
}

func (s *TimeLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TimeLiteralContext) TIME() antlr.TerminalNode {
	return s.GetToken(fhirpathParserTIME, 0)
}

func (s *TimeLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitTimeLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type NullLiteralContext struct {
	LiteralContext
}

func NewNullLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NullLiteralContext {
	var p = new(NullLiteralContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

	return p
}

func (s *NullLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NullLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitNullLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type DateTimeLiteralContext struct {
	LiteralContext
}

func NewDateTimeLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DateTimeLiteralContext {
	var p = new(DateTimeLiteralContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

	return p
}

func (s *DateTimeLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DateTimeLiteralContext) DATETIME() antlr.TerminalNode {
	return s.GetToken(fhirpathParserDATETIME, 0)
}

func (s *DateTimeLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitDateTimeLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type StringLiteralContext struct {
	LiteralContext
}

func NewStringLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringLiteralContext {
	var p = new(StringLiteralContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

	return p
}

func (s *StringLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StringLiteralContext) STRING() antlr.TerminalNode {
	return s.GetToken(fhirpathParserSTRING, 0)
}

func (s *StringLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitStringLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type DateLiteralContext struct {
	LiteralContext
}

func NewDateLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DateLiteralContext {
	var p = new(DateLiteralContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

	return p
}

func (s *DateLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DateLiteralContext) DATE() antlr.TerminalNode {
	return s.GetToken(fhirpathParserDATE, 0)
}

func (s *DateLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitDateLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type BooleanLiteralContext struct {
	LiteralContext
}

func NewBooleanLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BooleanLiteralContext {
	var p = new(BooleanLiteralContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

	return p
}

func (s *BooleanLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BooleanLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitBooleanLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type NumberLiteralContext struct {
	LiteralContext
}

func NewNumberLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NumberLiteralContext {
	var p = new(NumberLiteralContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

	return p
}

func (s *NumberLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NumberLiteralContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(fhirpathParserNUMBER, 0)
}

func (s *NumberLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitNumberLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

type QuantityLiteralContext struct {
	LiteralContext
}

func NewQuantityLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *QuantityLiteralContext {
	var p = new(QuantityLiteralContext)

	InitEmptyLiteralContext(&p.LiteralContext)
	p.parser = parser
	p.CopyAll(ctx.(*LiteralContext))

	return p
}

func (s *QuantityLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QuantityLiteralContext) Quantity() IQuantityContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQuantityContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQuantityContext)
}

func (s *QuantityLiteralContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitQuantityLiteral(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) Literal() (localctx ILiteralContext) {
	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, fhirpathParserRULE_literal)
	var _la int

	p.SetState(100)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		localctx = NewNullLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(91)
			p.Match(fhirpathParserT__29)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(92)
			p.Match(fhirpathParserT__30)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		localctx = NewBooleanLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(93)
			_la = p.GetTokenStream().LA(1)

			if !(_la == fhirpathParserT__31 || _la == fhirpathParserT__32) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	case 3:
		localctx = NewStringLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(94)
			p.Match(fhirpathParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		localctx = NewNumberLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(95)
			p.Match(fhirpathParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		localctx = NewDateLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(96)
			p.Match(fhirpathParserDATE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		localctx = NewDateTimeLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(97)
			p.Match(fhirpathParserDATETIME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		localctx = NewTimeLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(98)
			p.Match(fhirpathParserTIME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		localctx = NewQuantityLiteralContext(p, localctx)
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(99)
			p.Quantity()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExternalConstantContext is an interface to support dynamic dispatch.
type IExternalConstantContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	STRING() antlr.TerminalNode

	// IsExternalConstantContext differentiates from other interfaces.
	IsExternalConstantContext()
}

type ExternalConstantContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExternalConstantContext() *ExternalConstantContext {
	var p = new(ExternalConstantContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_externalConstant
	return p
}

func InitEmptyExternalConstantContext(p *ExternalConstantContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_externalConstant
}

func (*ExternalConstantContext) IsExternalConstantContext() {}

func NewExternalConstantContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExternalConstantContext {
	var p = new(ExternalConstantContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_externalConstant

	return p
}

func (s *ExternalConstantContext) GetParser() antlr.Parser { return s.parser }

func (s *ExternalConstantContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *ExternalConstantContext) STRING() antlr.TerminalNode {
	return s.GetToken(fhirpathParserSTRING, 0)
}

func (s *ExternalConstantContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExternalConstantContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExternalConstantContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitExternalConstant(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) ExternalConstant() (localctx IExternalConstantContext) {
	localctx = NewExternalConstantContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, fhirpathParserRULE_externalConstant)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(102)
		p.Match(fhirpathParserT__33)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(105)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case fhirpathParserT__10, fhirpathParserT__11, fhirpathParserT__21, fhirpathParserT__22, fhirpathParserIDENTIFIER, fhirpathParserDELIMITEDIDENTIFIER:
		{
			p.SetState(103)
			p.Identifier()
		}

	case fhirpathParserSTRING:
		{
			p.SetState(104)
			p.Match(fhirpathParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IInvocationContext is an interface to support dynamic dispatch.
type IInvocationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsInvocationContext differentiates from other interfaces.
	IsInvocationContext()
}

type InvocationContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInvocationContext() *InvocationContext {
	var p = new(InvocationContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_invocation
	return p
}

func InitEmptyInvocationContext(p *InvocationContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_invocation
}

func (*InvocationContext) IsInvocationContext() {}

func NewInvocationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InvocationContext {
	var p = new(InvocationContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_invocation

	return p
}

func (s *InvocationContext) GetParser() antlr.Parser { return s.parser }

func (s *InvocationContext) CopyAll(ctx *InvocationContext) {
	s.CopyFrom(&ctx.BaseParserRuleContext)
}

func (s *InvocationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InvocationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type TotalInvocationContext struct {
	InvocationContext
}

func NewTotalInvocationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TotalInvocationContext {
	var p = new(TotalInvocationContext)

	InitEmptyInvocationContext(&p.InvocationContext)
	p.parser = parser
	p.CopyAll(ctx.(*InvocationContext))

	return p
}

func (s *TotalInvocationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TotalInvocationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitTotalInvocation(s)

	default:
		return t.VisitChildren(s)
	}
}

type ThisInvocationContext struct {
	InvocationContext
}

func NewThisInvocationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ThisInvocationContext {
	var p = new(ThisInvocationContext)

	InitEmptyInvocationContext(&p.InvocationContext)
	p.parser = parser
	p.CopyAll(ctx.(*InvocationContext))

	return p
}

func (s *ThisInvocationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ThisInvocationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitThisInvocation(s)

	default:
		return t.VisitChildren(s)
	}
}

type IndexInvocationContext struct {
	InvocationContext
}

func NewIndexInvocationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IndexInvocationContext {
	var p = new(IndexInvocationContext)

	InitEmptyInvocationContext(&p.InvocationContext)
	p.parser = parser
	p.CopyAll(ctx.(*InvocationContext))

	return p
}

func (s *IndexInvocationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IndexInvocationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitIndexInvocation(s)

	default:
		return t.VisitChildren(s)
	}
}

type FunctionInvocationContext struct {
	InvocationContext
}

func NewFunctionInvocationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FunctionInvocationContext {
	var p = new(FunctionInvocationContext)

	InitEmptyInvocationContext(&p.InvocationContext)
	p.parser = parser
	p.CopyAll(ctx.(*InvocationContext))

	return p
}

func (s *FunctionInvocationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionInvocationContext) Function() IFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionContext)
}

func (s *FunctionInvocationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitFunctionInvocation(s)

	default:
		return t.VisitChildren(s)
	}
}

type MemberInvocationContext struct {
	InvocationContext
}

func NewMemberInvocationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MemberInvocationContext {
	var p = new(MemberInvocationContext)

	InitEmptyInvocationContext(&p.InvocationContext)
	p.parser = parser
	p.CopyAll(ctx.(*InvocationContext))

	return p
}

func (s *MemberInvocationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MemberInvocationContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *MemberInvocationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitMemberInvocation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) Invocation() (localctx IInvocationContext) {
	localctx = NewInvocationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, fhirpathParserRULE_invocation)
	p.SetState(112)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 6, p.GetParserRuleContext()) {
	case 1:
		localctx = NewMemberInvocationContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(107)
			p.Identifier()
		}

	case 2:
		localctx = NewFunctionInvocationContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(108)
			p.Function()
		}

	case 3:
		localctx = NewThisInvocationContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(109)
			p.Match(fhirpathParserT__34)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		localctx = NewIndexInvocationContext(p, localctx)
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(110)
			p.Match(fhirpathParserT__35)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		localctx = NewTotalInvocationContext(p, localctx)
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(111)
			p.Match(fhirpathParserT__36)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctionContext is an interface to support dynamic dispatch.
type IFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Identifier() IIdentifierContext
	ParamList() IParamListContext

	// IsFunctionContext differentiates from other interfaces.
	IsFunctionContext()
}

type FunctionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionContext() *FunctionContext {
	var p = new(FunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_function
	return p
}

func InitEmptyFunctionContext(p *FunctionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_function
}

func (*FunctionContext) IsFunctionContext() {}

func NewFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionContext {
	var p = new(FunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_function

	return p
}

func (s *FunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionContext) Identifier() IIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *FunctionContext) ParamList() IParamListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IParamListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IParamListContext)
}

func (s *FunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) Function() (localctx IFunctionContext) {
	localctx = NewFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, fhirpathParserRULE_function)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(114)
		p.Identifier()
	}
	{
		p.SetState(115)
		p.Match(fhirpathParserT__27)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(117)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&4575657493346129968) != 0 {
		{
			p.SetState(116)
			p.ParamList()
		}

	}
	{
		p.SetState(119)
		p.Match(fhirpathParserT__28)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IParamListContext is an interface to support dynamic dispatch.
type IParamListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext

	// IsParamListContext differentiates from other interfaces.
	IsParamListContext()
}

type ParamListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyParamListContext() *ParamListContext {
	var p = new(ParamListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_paramList
	return p
}

func InitEmptyParamListContext(p *ParamListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_paramList
}

func (*ParamListContext) IsParamListContext() {}

func NewParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamListContext {
	var p = new(ParamListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_paramList

	return p
}

func (s *ParamListContext) GetParser() antlr.Parser { return s.parser }

func (s *ParamListContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ParamListContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ParamListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ParamListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ParamListContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitParamList(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) ParamList() (localctx IParamListContext) {
	localctx = NewParamListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, fhirpathParserRULE_paramList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(121)
		p.expression(0)
	}
	p.SetState(126)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == fhirpathParserT__37 {
		{
			p.SetState(122)
			p.Match(fhirpathParserT__37)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(123)
			p.expression(0)
		}

		p.SetState(128)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IQuantityContext is an interface to support dynamic dispatch.
type IQuantityContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMBER() antlr.TerminalNode
	Unit() IUnitContext

	// IsQuantityContext differentiates from other interfaces.
	IsQuantityContext()
}

type QuantityContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQuantityContext() *QuantityContext {
	var p = new(QuantityContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_quantity
	return p
}

func InitEmptyQuantityContext(p *QuantityContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_quantity
}

func (*QuantityContext) IsQuantityContext() {}

func NewQuantityContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QuantityContext {
	var p = new(QuantityContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_quantity

	return p
}

func (s *QuantityContext) GetParser() antlr.Parser { return s.parser }

func (s *QuantityContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(fhirpathParserNUMBER, 0)
}

func (s *QuantityContext) Unit() IUnitContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnitContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnitContext)
}

func (s *QuantityContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QuantityContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QuantityContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitQuantity(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) Quantity() (localctx IQuantityContext) {
	localctx = NewQuantityContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, fhirpathParserRULE_quantity)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(129)
		p.Match(fhirpathParserNUMBER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(131)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 9, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(130)
			p.Unit()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnitContext is an interface to support dynamic dispatch.
type IUnitContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DateTimePrecision() IDateTimePrecisionContext
	PluralDateTimePrecision() IPluralDateTimePrecisionContext
	STRING() antlr.TerminalNode

	// IsUnitContext differentiates from other interfaces.
	IsUnitContext()
}

type UnitContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnitContext() *UnitContext {
	var p = new(UnitContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_unit
	return p
}

func InitEmptyUnitContext(p *UnitContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_unit
}

func (*UnitContext) IsUnitContext() {}

func NewUnitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnitContext {
	var p = new(UnitContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_unit

	return p
}

func (s *UnitContext) GetParser() antlr.Parser { return s.parser }

func (s *UnitContext) DateTimePrecision() IDateTimePrecisionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDateTimePrecisionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDateTimePrecisionContext)
}

func (s *UnitContext) PluralDateTimePrecision() IPluralDateTimePrecisionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPluralDateTimePrecisionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPluralDateTimePrecisionContext)
}

func (s *UnitContext) STRING() antlr.TerminalNode {
	return s.GetToken(fhirpathParserSTRING, 0)
}

func (s *UnitContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnitContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnitContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitUnit(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) Unit() (localctx IUnitContext) {
	localctx = NewUnitContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, fhirpathParserRULE_unit)
	p.SetState(136)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case fhirpathParserT__38, fhirpathParserT__39, fhirpathParserT__40, fhirpathParserT__41, fhirpathParserT__42, fhirpathParserT__43, fhirpathParserT__44, fhirpathParserT__45:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(133)
			p.DateTimePrecision()
		}

	case fhirpathParserT__46, fhirpathParserT__47, fhirpathParserT__48, fhirpathParserT__49, fhirpathParserT__50, fhirpathParserT__51, fhirpathParserT__52, fhirpathParserT__53:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(134)
			p.PluralDateTimePrecision()
		}

	case fhirpathParserSTRING:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(135)
			p.Match(fhirpathParserSTRING)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDateTimePrecisionContext is an interface to support dynamic dispatch.
type IDateTimePrecisionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsDateTimePrecisionContext differentiates from other interfaces.
	IsDateTimePrecisionContext()
}

type DateTimePrecisionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDateTimePrecisionContext() *DateTimePrecisionContext {
	var p = new(DateTimePrecisionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_dateTimePrecision
	return p
}

func InitEmptyDateTimePrecisionContext(p *DateTimePrecisionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_dateTimePrecision
}

func (*DateTimePrecisionContext) IsDateTimePrecisionContext() {}

func NewDateTimePrecisionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DateTimePrecisionContext {
	var p = new(DateTimePrecisionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_dateTimePrecision

	return p
}

func (s *DateTimePrecisionContext) GetParser() antlr.Parser { return s.parser }
func (s *DateTimePrecisionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DateTimePrecisionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DateTimePrecisionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitDateTimePrecision(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) DateTimePrecision() (localctx IDateTimePrecisionContext) {
	localctx = NewDateTimePrecisionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, fhirpathParserRULE_dateTimePrecision)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(138)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&140187732541440) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPluralDateTimePrecisionContext is an interface to support dynamic dispatch.
type IPluralDateTimePrecisionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsPluralDateTimePrecisionContext differentiates from other interfaces.
	IsPluralDateTimePrecisionContext()
}

type PluralDateTimePrecisionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPluralDateTimePrecisionContext() *PluralDateTimePrecisionContext {
	var p = new(PluralDateTimePrecisionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_pluralDateTimePrecision
	return p
}

func InitEmptyPluralDateTimePrecisionContext(p *PluralDateTimePrecisionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_pluralDateTimePrecision
}

func (*PluralDateTimePrecisionContext) IsPluralDateTimePrecisionContext() {}

func NewPluralDateTimePrecisionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PluralDateTimePrecisionContext {
	var p = new(PluralDateTimePrecisionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_pluralDateTimePrecision

	return p
}

func (s *PluralDateTimePrecisionContext) GetParser() antlr.Parser { return s.parser }
func (s *PluralDateTimePrecisionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PluralDateTimePrecisionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PluralDateTimePrecisionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitPluralDateTimePrecision(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) PluralDateTimePrecision() (localctx IPluralDateTimePrecisionContext) {
	localctx = NewPluralDateTimePrecisionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, fhirpathParserRULE_pluralDateTimePrecision)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(140)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&35888059530608640) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeSpecifierContext is an interface to support dynamic dispatch.
type ITypeSpecifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QualifiedIdentifier() IQualifiedIdentifierContext

	// IsTypeSpecifierContext differentiates from other interfaces.
	IsTypeSpecifierContext()
}

type TypeSpecifierContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeSpecifierContext() *TypeSpecifierContext {
	var p = new(TypeSpecifierContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_typeSpecifier
	return p
}

func InitEmptyTypeSpecifierContext(p *TypeSpecifierContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_typeSpecifier
}

func (*TypeSpecifierContext) IsTypeSpecifierContext() {}

func NewTypeSpecifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeSpecifierContext {
	var p = new(TypeSpecifierContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_typeSpecifier

	return p
}

func (s *TypeSpecifierContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeSpecifierContext) QualifiedIdentifier() IQualifiedIdentifierContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedIdentifierContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedIdentifierContext)
}

func (s *TypeSpecifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeSpecifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeSpecifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitTypeSpecifier(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) TypeSpecifier() (localctx ITypeSpecifierContext) {
	localctx = NewTypeSpecifierContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, fhirpathParserRULE_typeSpecifier)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(142)
		p.QualifiedIdentifier()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IQualifiedIdentifierContext is an interface to support dynamic dispatch.
type IQualifiedIdentifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIdentifier() []IIdentifierContext
	Identifier(i int) IIdentifierContext

	// IsQualifiedIdentifierContext differentiates from other interfaces.
	IsQualifiedIdentifierContext()
}

type QualifiedIdentifierContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQualifiedIdentifierContext() *QualifiedIdentifierContext {
	var p = new(QualifiedIdentifierContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_qualifiedIdentifier
	return p
}

func InitEmptyQualifiedIdentifierContext(p *QualifiedIdentifierContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_qualifiedIdentifier
}

func (*QualifiedIdentifierContext) IsQualifiedIdentifierContext() {}

func NewQualifiedIdentifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QualifiedIdentifierContext {
	var p = new(QualifiedIdentifierContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_qualifiedIdentifier

	return p
}

func (s *QualifiedIdentifierContext) GetParser() antlr.Parser { return s.parser }

func (s *QualifiedIdentifierContext) AllIdentifier() []IIdentifierContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IIdentifierContext); ok {
			len++
		}
	}

	tst := make([]IIdentifierContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IIdentifierContext); ok {
			tst[i] = t.(IIdentifierContext)
			i++
		}
	}

	return tst
}

func (s *QualifiedIdentifierContext) Identifier(i int) IIdentifierContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifierContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifierContext)
}

func (s *QualifiedIdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QualifiedIdentifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QualifiedIdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitQualifiedIdentifier(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) QualifiedIdentifier() (localctx IQualifiedIdentifierContext) {
	localctx = NewQualifiedIdentifierContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, fhirpathParserRULE_qualifiedIdentifier)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(144)
		p.Identifier()
	}
	p.SetState(149)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 11, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(145)
				p.Match(fhirpathParserT__0)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(146)
				p.Identifier()
			}

		}
		p.SetState(151)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 11, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIdentifierContext is an interface to support dynamic dispatch.
type IIdentifierContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	DELIMITEDIDENTIFIER() antlr.TerminalNode

	// IsIdentifierContext differentiates from other interfaces.
	IsIdentifierContext()
}

type IdentifierContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentifierContext() *IdentifierContext {
	var p = new(IdentifierContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_identifier
	return p
}

func InitEmptyIdentifierContext(p *IdentifierContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = fhirpathParserRULE_identifier
}

func (*IdentifierContext) IsIdentifierContext() {}

func NewIdentifierContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IdentifierContext {
	var p = new(IdentifierContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = fhirpathParserRULE_identifier

	return p
}

func (s *IdentifierContext) GetParser() antlr.Parser { return s.parser }

func (s *IdentifierContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(fhirpathParserIDENTIFIER, 0)
}

func (s *IdentifierContext) DELIMITEDIDENTIFIER() antlr.TerminalNode {
	return s.GetToken(fhirpathParserDELIMITEDIDENTIFIER, 0)
}

func (s *IdentifierContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifierContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IdentifierContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case fhirpathVisitor:
		return t.VisitIdentifier(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *fhirpathParser) Identifier() (localctx IIdentifierContext) {
	localctx = NewIdentifierContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, fhirpathParserRULE_identifier)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(152)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&864691128467724288) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

func (p *fhirpathParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 1:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *fhirpathParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 10)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 9)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 7)

	case 3:
		return p.Precpred(p.GetParserRuleContext(), 6)

	case 4:
		return p.Precpred(p.GetParserRuleContext(), 5)

	case 5:
		return p.Precpred(p.GetParserRuleContext(), 4)

	case 6:
		return p.Precpred(p.GetParserRuleContext(), 3)

	case 7:
		return p.Precpred(p.GetParserRuleContext(), 2)

	case 8:
		return p.Precpred(p.GetParserRuleContext(), 1)

	case 9:
		return p.Precpred(p.GetParserRuleContext(), 13)

	case 10:
		return p.Precpred(p.GetParserRuleContext(), 12)

	case 11:
		return p.Precpred(p.GetParserRuleContext(), 8)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
