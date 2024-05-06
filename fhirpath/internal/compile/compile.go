package compile

import (
	"github.com/antlr4-go/antlr/v4"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/funcs"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/grammar"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/opts"
	"github.com/verily-src/fhirpath-go/fhirpath/internal/parser"
)

// PopulateConfig creates a CompileConfig and prepopulates it with
// a function table and any provided options.
func PopulateConfig(options ...opts.CompileOption) (*opts.CompileConfig, error) {
	config := &opts.CompileConfig{
		Table: funcs.Clone(),
	}
	config, err := opts.ApplyOptions(config, options...)
	if err != nil {
		return nil, err
	}
	return config, err
}

// Tree creates an ANTLR parsing context from the provided FHIRPath string.
func Tree(expr string) (grammar.IProgContext, error) {
	inputStream := antlr.NewInputStream(expr)
	errorListener := &parser.FHIRPathErrorListener{}

	// Lex the input stream
	lexer := grammar.NewfhirpathLexer(inputStream)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Parse the tokens
	p := grammar.NewfhirpathParser(tokens)
	p.RemoveErrorListeners()
	p.AddErrorListener(errorListener)
	tree := p.Prog()

	if err := errorListener.Error(); err != nil {
		return nil, err
	}

	return tree, nil
}
