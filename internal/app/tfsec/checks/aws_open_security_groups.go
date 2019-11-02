package checks

import (
	"fmt"
	"strings"

	"github.com/liamg/tfsec/internal/app/tfsec/scanner"

	"github.com/liamg/tfsec/internal/app/tfsec/parser"
)

// AWSOpenIngressSecurityGroupInlineRule See https://github.com/liamg/tfsec#included-checks for check info
const AWSOpenIngressSecurityGroupInlineRule scanner.CheckCode = "AWS008"

// AWSOpenEgressSecurityGroupInlineRule See https://github.com/liamg/tfsec#included-checks for check info
const AWSOpenEgressSecurityGroupInlineRule scanner.CheckCode = "AWS009"

func init() {
	scanner.RegisterCheck(scanner.Check{
		Code:           AWSOpenIngressSecurityGroupInlineRule,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_security_group"},
		CheckFunc: func(check *scanner.Check, block *parser.Block) []scanner.Result {

			var results []scanner.Result

			if directionBlock := block.GetBlock("ingress"); directionBlock != nil {
				if cidrBlocksAttr := directionBlock.GetAttribute("cidr_blocks"); cidrBlocksAttr != nil {

					if cidrBlocksAttr.Value().LengthInt() == 0 {
						return nil
					}

					for _, cidr := range cidrBlocksAttr.Value().AsValueSlice() {
						if strings.HasSuffix(cidr.AsString(), "/0") {
							results = append(results,
								check.NewResult(
									fmt.Sprintf("Resource '%s' defines a fully open ingress security group.", block.Name()),
									cidrBlocksAttr.Range(),
								),
							)
						}
					}
				}
			}

			return results
		},
	})

	scanner.RegisterCheck(scanner.Check{
		Code:           AWSOpenEgressSecurityGroupInlineRule,
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"aws_security_group"},
		CheckFunc: func(check *scanner.Check, block *parser.Block) []scanner.Result {

			var results []scanner.Result

			if directionBlock := block.GetBlock("egress"); directionBlock != nil {
				if cidrBlocksAttr := directionBlock.GetAttribute("cidr_blocks"); cidrBlocksAttr != nil {

					if cidrBlocksAttr.Value().LengthInt() == 0 {
						return nil
					}

					for _, cidr := range cidrBlocksAttr.Value().AsValueSlice() {
						if strings.HasSuffix(cidr.AsString(), "/0") {
							results = append(results,
								check.NewResultWithValueAnnotation(
									fmt.Sprintf("Resource '%s' defines a fully open egress security group.", block.Name()),
									cidrBlocksAttr.Range(),
									cidrBlocksAttr,
								),
							)
						}
					}
				}
			}

			return results
		},
	})
}
