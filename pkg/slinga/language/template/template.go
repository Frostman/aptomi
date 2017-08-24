package template

import (
	"bytes"
	"fmt"
	"strings"
	t "text/template"
	"github.com/Aptomi/aptomi/pkg/slinga/errors"
)

type Template struct {
	templateStr      string
	templateCompiled *t.Template
}

func NewTemplate(templateStr string) (*Template, error) {
	templateCompiled, err := t.New("").Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("Unable to compile template '%s': %s", templateStr, err.Error())
	}
	return &Template{
		templateStr:      templateStr,
		templateCompiled: templateCompiled,
	}, nil
}

// Evaluate an expression, given a set of labels
func (template *Template) Evaluate(params *TemplateParameters) (string, error) {
	// Evaluate
	var doc bytes.Buffer
	err := template.templateCompiled.Execute(&doc, params.params)
	if err != nil {
		return "", errors.NewErrorWithDetails(
			fmt.Sprintf("Unable to evaluate template '%s': %s", template.templateStr, err.Error()),
			errors.Details{
				"template": template.templateStr,
				"params": params,
			},
		)
	}

	// Convert result to bool
	result := doc.String()
	if strings.Contains(result, "<no value>") {
		return "", errors.NewErrorWithDetails(
			fmt.Sprintf("Unable to evaluate template '%s': <no value>", template.templateStr),
			errors.Details{
				"template": template.templateStr,
				"params": params,
			},
		)
	}

	return result, nil
}
