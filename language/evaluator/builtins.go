package evaluator

import (
	"fmt"
	"glass/language/object"
)

var builtins = map[string]*object.Builtin{
	"print": {
		Function: func(arguments ...object.Object) object.Object {
			for _, argument := range arguments {
				fmt.Print(argument.Inspect())
			}

			return NULL
		},
	},
}
