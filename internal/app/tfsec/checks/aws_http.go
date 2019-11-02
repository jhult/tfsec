package checks

import (
	"fmt"

	"github.com/liamg/tfsec/internal/app/tfsec/scanner"

	"github.com/liamg/tfsec/internal/app/tfsec/parser"
	"github.com/zclconf/go-cty/cty"
)

// AWSPlainHTTP See https://github.com/liamg/tfsec#included-checks for check info
const AWSPlainHTTP scanner.CheckCode = "AWS004"

func init() {
	scanner.RegisterCheck(scanner.Check{
		Code:           AWSPlainHTTP,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_lb_listener", "aws_alb_listener"},
		CheckFunc: func(check *scanner.Check, block *parser.Block) []scanner.Result {
			if protocolAttr := block.GetAttribute("protocol"); protocolAttr == nil || (protocolAttr.Type() == cty.String && protocolAttr.Value().AsString() == "HTTP") {
				// check if this is a redirect to HTTPS - if it is, then no problem
				if actionBlock := block.GetBlock("default_action"); actionBlock != nil {
					actionTypeAttr := actionBlock.GetAttribute("type")
					if actionTypeAttr != nil && actionTypeAttr.Type() == cty.String && actionTypeAttr.Value().AsString() == "redirect" {
						if redirectBlock := actionBlock.GetBlock("redirect"); redirectBlock != nil {
							redirectProtocolAttr := redirectBlock.GetAttribute("protocol")
							if redirectProtocolAttr != nil && redirectProtocolAttr.Type() == cty.String && redirectProtocolAttr.Value().AsString() == "HTTPS" {
								return nil
							}
						}
					}
				}
				reportRange := block.Range()
				if protocolAttr != nil {
					reportRange = protocolAttr.Range()
				}
				return []scanner.Result{
					check.NewResultWithValueAnnotation(
						fmt.Sprintf("Resource '%s' uses plain HTTP instead of HTTPS.", block.Name()),
						reportRange,
						protocolAttr,
					),
				}
			}
			return nil
		},
	})
}
