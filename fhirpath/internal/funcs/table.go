package funcs

import "github.com/verily-src/fhirpath-go/fhirpath/internal/funcs/impl"

// BaseTable holds the default mapping of all
// FHIRPath functions. Unimplemented functions return an
// unimplemented error.
// This table is a part of the N1 normative spec.
// See https://hl7.org/fhirpath/N1/
var baseTable = FunctionTable{
	"empty": Function{
		impl.Empty,
		0,
		0,
		false,
	},
	"exists": Function{
		impl.Exists,
		0,
		1,
		false,
	},
	"extension": Function{
		impl.Extension,
		1,
		1,
		false,
	},
	"all": Function{
		impl.All,
		1,
		1,
		false,
	},
	"allTrue": Function{
		impl.AllTrue,
		0,
		0,
		false,
	},
	"anyTrue": Function{
		impl.AnyTrue,
		0,
		0,
		false,
	},
	"allFalse": Function{
		impl.AllFalse,
		0,
		0,
		false,
	},
	"anyFalse": Function{
		impl.AnyFalse,
		0,
		0,
		false,
	},
	"subsetOf":   notImplemented,
	"supersetOf": notImplemented,
	"count": Function{
		impl.Count,
		0,
		0,
		false,
	},
	"distinct": Function{
		impl.Distinct,
		0,
		0,
		false,
	},
	"isDistinct": Function{
		impl.IsDistinct,
		0,
		0,
		false,
	},
	"where": Function{
		impl.Where,
		1,
		1,
		false,
	},
	"select": Function{
		impl.Select,
		1,
		1,
		false,
	},
	"repeat": notImplemented,
	"ofType": Function{
		impl.OfType,
		1,
		1,
		true,
	},
	"single": notImplemented,
	"first": Function{
		impl.First,
		0,
		0,
		false,
	},
	"last": Function{
		impl.Last,
		0,
		0,
		false,
	},
	"tail": Function{
		impl.Tail,
		0,
		0,
		false,
	},
	"skip": Function{
		impl.Skip,
		1,
		1,
		false,
	},
	"take": Function{
		impl.Take,
		1,
		1,
		false,
	},
	"intersect": Function{
		impl.Intersect,
		1,
		1,
		false,
	},
	"exclude": Function{
		impl.Exclude,
		1,
		1,
		false,
	},
	"union": notImplemented,
	"combine": Function{
		impl.Combine,
		0,
		1,
		false,
	},
	"iif": Function{
		impl.Iif,
		2,
		3,
		false,
	},
	"toBoolean": Function{
		impl.ToBoolean,
		0,
		0,
		false,
	},
	"convertsToBoolean": Function{
		impl.ConvertsToBoolean,
		0,
		0,
		false,
	},
	"toInteger": Function{
		impl.ToInteger,
		0,
		0,
		false,
	},
	"convertsToInteger": Function{
		impl.ConvertsToInteger,
		0,
		0,
		false,
	},
	"toDate": Function{
		impl.ToDate,
		0,
		0,
		false,
	},
	"convertsToDate": Function{
		impl.ConvertsToDate,
		0,
		0,
		false,
	},
	"toDateTime": Function{
		impl.ToDateTime,
		0,
		0,
		false,
	},
	"convertToDateTime": Function{
		impl.ConvertsToDateTime,
		0,
		0,
		false,
	},
	"toDecimal": Function{
		impl.ToDecimal,
		0,
		0,
		false,
	},
	"convertsToDecimal": Function{
		impl.ConvertsToDecimal,
		0,
		0,
		false,
	},
	"toQuantity": Function{
		impl.ToInteger,
		0,
		1,
		false,
	},
	"convertsToQuantity": Function{
		impl.ConvertsToQuantity,
		0,
		1,
		false,
	},
	"toString": Function{
		impl.ToString,
		0,
		0,
		false,
	},
	"convertsToString": Function{
		impl.ConvertsToString,
		0,
		0,
		false,
	},
	"toTime": Function{
		impl.ToTime,
		0,
		0,
		false,
	},
	"convertsToTime": Function{
		impl.ConvertsToTime,
		0,
		0,
		false,
	},
	"indexOf": Function{
		impl.IndexOf,
		1,
		1,
		false,
	},
	"substring": Function{
		impl.Substring,
		1,
		2,
		false,
	},
	"startsWith": Function{
		impl.StartsWith,
		1,
		1,
		false,
	},
	"endsWith": Function{
		impl.EndsWith,
		1,
		1,
		false,
	},
	"contains": Function{
		impl.Contains,
		1,
		1,
		false,
	},
	"upper": Function{
		impl.Upper,
		0,
		0,
		false,
	},
	"lower": Function{
		impl.Lower,
		0,
		0,
		false,
	},
	"replace": Function{
		impl.Replace,
		2,
		2,
		false,
	},
	"matches": Function{
		impl.Matches,
		1,
		1,
		false,
	},
	"replaceMatches": Function{
		impl.ReplaceMatches,
		2,
		2,
		false,
	},
	"length": Function{
		impl.Length,
		0,
		0,
		false,
	},
	"toChars": Function{
		impl.ToChars,
		0,
		0,
		false,
	},
	"abs": Function{
		impl.Abs,
		0,
		0,
		false,
	},
	"ceiling": Function{
		impl.Ceiling,
		0,
		0,
		false,
	},
	"exp": Function{
		impl.Exp,
		0,
		0,
		false,
	},
	"floor": Function{
		impl.Floor,
		0,
		0,
		false,
	},
	"ln": Function{
		impl.Ln,
		0,
		0,
		false,
	},
	"log": Function{
		impl.Log,
		0,
		0,
		false,
	},
	"power": Function{
		impl.Power,
		0,
		0,
		false,
	},
	"round": Function{
		impl.Round,
		0,
		0,
		false,
	},
	"sqrt": Function{
		impl.Sqrt,
		0,
		0,
		false,
	},
	"truncate": Function{
		impl.Truncate,
		0,
		0,
		false,
	},
	"children": Function{
		impl.Children,
		0,
		0,
		false,
	},
	"descendants": Function{
		impl.Descendants,
		0,
		0,
		false,
	},
	"trace": notImplemented,
	"now": Function{
		impl.Now,
		0,
		0,
		false,
	},
	"timeOfDay": Function{
		impl.TimeOfDay,
		0,
		0,
		false,
	},
	"today": Function{
		impl.Today,
		0,
		0,
		false,
	},
	"not": Function{
		impl.Not,
		0,
		0,
		false,
	},
	"resolve": Function{
		impl.Resolve,
		0,
		0,
		false,
	},
}

// ExperimentalTable holds the mapping of all
// experimental FHIRPath functions. These functions
// are not a part of the N1 normative spec.
// See https://build.fhir.org/ig/HL7/FHIRPath/
var experimentalTable = FunctionTable{
	"join": Function{
		impl.Join,
		0,
		1,
		false,
	},
}

// Clone returns a deep copy of the base
// function table.
func Clone() FunctionTable {
	// TODO: Optimize
	table := make(FunctionTable)
	for k, v := range baseTable {
		table[k] = v
	}
	return table
}

// AddExperimentalFuncs adds experimental functions
// to the given function table and returns it.
// If a function already exists in the table, it is not overridden.
func AddExperimentalFuncs(table FunctionTable) FunctionTable {
	for k, v := range experimentalTable {
		if _, exists := table[k]; exists {
			continue
		}
		table[k] = v
	}
	return table
}
