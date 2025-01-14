// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.
package oa3gen

import (
	"strings"
	"time"

	"github.com/aarondl/chrono"
	"github.com/aarondl/oa3/support"
	"github.com/aarondl/opt/null"
	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Checks to see that all Go primitives work
type Primitives struct {
	Bool         bool                        `json:"bool"`
	BoolNull     null.Val[bool]              `json:"bool_null"`
	DateNull     null.Val[chrono.Date]       `json:"date_null"`
	DateVal      chrono.Date                 `json:"date_val"`
	DatetimeNull null.Val[chrono.DateTime]   `json:"datetime_null"`
	DatetimeVal  chrono.DateTime             `json:"datetime_val"`
	Decimal      decimal.Decimal             `json:"decimal"`
	DecimalNull  null.Val[decimal.Decimal]   `json:"decimal_null"`
	DurationNull omitnull.Val[time.Duration] `json:"duration_null,omitempty"`
	DurationVal  omit.Val[time.Duration]     `json:"duration_val,omitempty"`
	Float        float64                     `json:"float"`
	Float32      float32                     `json:"float32"`
	Float32Null  null.Val[float32]           `json:"float32_null"`
	Float64      float64                     `json:"float64"`
	Float64Null  null.Val[float64]           `json:"float64_null"`
	FloatNull    null.Val[float64]           `json:"float_null"`
	// Normal int
	Int       int                         `json:"int"`
	Int32     int32                       `json:"int32"`
	Int32Null null.Val[int32]             `json:"int32_null"`
	Int64     int64                       `json:"int64"`
	Int64Null null.Val[int64]             `json:"int64_null"`
	IntNull   null.Val[int]               `json:"int_null"`
	Str       PrimitivesStr               `json:"str"`
	StrNull   null.Val[PrimitivesStrNull] `json:"str_null"`
	TimeNull  null.Val[chrono.Time]       `json:"time_null"`
	TimeVal   chrono.Time                 `json:"time_val"`
	Uuid      uuid.UUID                   `json:"uuid"`
	UuidNull  null.Val[uuid.UUID]         `json:"uuid_null"`
}

// validateSchema validates the object and returns
// errors that can be returned to the user.
func (o Primitives) validateSchema() support.Errors {
	var ctx []string
	var ers []error
	var errs support.Errors
	_, _, _ = ctx, ers, errs

	ers = nil
	if err := support.ValidateMultipleOfFloat(o.Float, 5.5); err != nil {
		ers = append(ers, err)
	}
	if len(ers) != 0 {
		ctx = append(ctx, "float")
		errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
		ctx = ctx[:len(ctx)-1]
	}
	ers = nil
	if err := support.ValidateMaxNumber(o.Float32, 5.5, false); err != nil {
		ers = append(ers, err)
	}
	if len(ers) != 0 {
		ctx = append(ctx, "float32")
		errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
		ctx = ctx[:len(ctx)-1]
	}

	if val, ok := o.Float32Null.Get(); ok {

		ers = nil
		if err := support.ValidateMaxNumber(val, 5, false); err != nil {
			ers = append(ers, err)
		}
		if len(ers) != 0 {
			ctx = append(ctx, "float32_null")
			errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
			ctx = ctx[:len(ctx)-1]
		}
	}
	ers = nil
	if err := support.ValidateMinNumber(o.Float64, 5.5, false); err != nil {
		ers = append(ers, err)
	}
	if len(ers) != 0 {
		ctx = append(ctx, "float64")
		errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
		ctx = ctx[:len(ctx)-1]
	}

	if val, ok := o.Float64Null.Get(); ok {

		ers = nil
		if err := support.ValidateMinNumber(val, 5, false); err != nil {
			ers = append(ers, err)
		}
		if len(ers) != 0 {
			ctx = append(ctx, "float64_null")
			errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
			ctx = ctx[:len(ctx)-1]
		}
	}

	if val, ok := o.FloatNull.Get(); ok {

		ers = nil
		if err := support.ValidateMultipleOfFloat(val, 5.5); err != nil {
			ers = append(ers, err)
		}
		if len(ers) != 0 {
			ctx = append(ctx, "float_null")
			errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
			ctx = ctx[:len(ctx)-1]
		}
	}
	ers = nil
	if err := support.ValidateMultipleOfInt(o.Int, 5); err != nil {
		ers = append(ers, err)
	}
	if len(ers) != 0 {
		ctx = append(ctx, "int")
		errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
		ctx = ctx[:len(ctx)-1]
	}
	ers = nil
	if err := support.ValidateMaxNumber(o.Int32, 5, false); err != nil {
		ers = append(ers, err)
	}
	if len(ers) != 0 {
		ctx = append(ctx, "int32")
		errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
		ctx = ctx[:len(ctx)-1]
	}

	if val, ok := o.Int32Null.Get(); ok {

		ers = nil
		if err := support.ValidateMaxNumber(val, 5, false); err != nil {
			ers = append(ers, err)
		}
		if len(ers) != 0 {
			ctx = append(ctx, "int32_null")
			errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
			ctx = ctx[:len(ctx)-1]
		}
	}
	ers = nil
	if err := support.ValidateMinNumber(o.Int64, 5, false); err != nil {
		ers = append(ers, err)
	}
	if len(ers) != 0 {
		ctx = append(ctx, "int64")
		errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
		ctx = ctx[:len(ctx)-1]
	}

	if val, ok := o.Int64Null.Get(); ok {

		ers = nil
		if err := support.ValidateMinNumber(val, 5, false); err != nil {
			ers = append(ers, err)
		}
		if len(ers) != 0 {
			ctx = append(ctx, "int64_null")
			errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
			ctx = ctx[:len(ctx)-1]
		}
	}

	if val, ok := o.IntNull.Get(); ok {

		ers = nil
		if err := support.ValidateMultipleOfInt(val, 5); err != nil {
			ers = append(ers, err)
		}
		if len(ers) != 0 {
			ctx = append(ctx, "int_null")
			errs = support.AddErrs(errs, strings.Join(ctx, "."), ers...)
			ctx = ctx[:len(ctx)-1]
		}
	}
	if newErrs := Validate(o.Str); newErrs != nil {
		ctx = append(ctx, "str")
		errs = support.AddErrsFlatten(errs, strings.Join(ctx, "."), newErrs)
		ctx = ctx[:len(ctx)-1]
	}
	if val, ok := o.StrNull.Get(); ok {
		if newErrs := Validate(val); newErrs != nil {
			ctx = append(ctx, "str_null")
			errs = support.AddErrsFlatten(errs, strings.Join(ctx, "."), newErrs)
			ctx = ctx[:len(ctx)-1]
		}
	}

	return errs
}
