// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package oa3gen

import (
	"strings"

	"github.com/aarondl/oa3/support"
)

// Recursively defined maps
type MapRecursive map[string]map[string]map[string]string

// ValidateSchemaMapRecursive validates the object and returns
// errors that can be returned to the user.
func (o MapRecursive) ValidateSchemaMapRecursive() support.Errors {
	var ctx []string
	var ers []error
	var errs support.Errors
	_, _ = ers, errs

	if err := support.ValidateMaxProperties(o, 3); err != nil {
		ers = append(ers, err)
	}
	if err := support.ValidateMinProperties(o, 2); err != nil {
		ers = append(ers, err)
	}
	for k, o := range o {
		var ers []error
		ctx = append(ctx, k)

		if err := support.ValidateMaxProperties(o, 4); err != nil {
			ers = append(ers, err)
		}
		if err := support.ValidateMinProperties(o, 3); err != nil {
			ers = append(ers, err)
		}
		for k, o := range o {
			var ers []error
			ctx = append(ctx, k)

			if err := support.ValidateMaxProperties(o, 6); err != nil {
				ers = append(ers, err)
			}
			if err := support.ValidateMinProperties(o, 5); err != nil {
				ers = append(ers, err)
			}
			for k, o := range o {
				var ers []error
				ctx = append(ctx, k)

				errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
				ctx = ctx[:len(ctx)-1]
			}
			errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
			ctx = ctx[:len(ctx)-1]
		}
		errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
		ctx = ctx[:len(ctx)-1]
	}

	errs = support.AddErrs(errs, "", ers...)

	return errs
}
