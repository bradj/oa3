// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package oa3gen

import (
	"github.com/aarondl/oa3/support"
)

type TestInlineInline struct {
	Id string `json:"id"`
}

// VVValidateSchemaTestInlineInline validates the object and returns
// errors that can be returned to the user.
func (o TestInlineInline) VVValidateSchema() support.Errors {
	var ctx []string
	var ers []error
	var errs support.Errors
	_, _, _ = ctx, ers, errs

	errs = support.AddErrs(errs, "", ers...)

	return errs
}
