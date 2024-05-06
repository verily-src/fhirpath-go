package parser

import (
	"errors"
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

type FHIRPathErrorListener struct {
	*antlr.DefaultErrorListener
	errors []error
}

func (l *FHIRPathErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	err := fmt.Errorf("syntax error on line %d:%d - %s", line, column, msg)
	l.errors = append(l.errors, err)
}

func (l *FHIRPathErrorListener) Error() error {
	return errors.Join(l.errors...)
}
