package controller

import "errors"

var (
	ErrEstablishmentNotFound         = errors.New("establishment not found")
	ErrQuantityMustBeGreaterThanZero = errors.New("quantity must be greater than zero")
	ErrQuantityMustBeAPositiveNumber = errors.New("quantity must be a positive number")
	ErrEstablishmentHasNoTables      = errors.New("establishment has no tables")
	ErrTableNotFound                 = errors.New("table not found")
	ErrProductNotFound               = errors.New("product not found")
	ErrTableNotAvailable             = errors.New("table not available")
	ErrTableIsInUse                  = errors.New("table is in use")
	ErrTableIsInAnotherEstablishment = errors.New("table is in another establishment")
	ErrNoRowsAffected                = errors.New("no rows affected")
)

type unauthorizedErr string

func (u unauthorizedErr) Error() string {
	return string(u)
}

func (u unauthorizedErr) IsUnauthorized() {}
