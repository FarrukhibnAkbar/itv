package constants

import "errors"

type Sentinel string

func (s Sentinel) Error() string {
	return string(s)
}

const (
	// ErrXozmakConstantNotExists ...
	ErrXozmakConstantNotExists = Sentinel("no company constant exists")
	//ErrXozmakAlreadyExists ...
	ErrXozmakAlreadyExists = Sentinel("the xozmak already exists")
	//ErrXozmakAlreadyExists ...
	ErrAdminAlreadyExists = Sentinel("the admin already exists")
	// ErrProductConstantNotExists ...
	ErrProductConstantNotExists = Sentinel("no product constant exists")
	//ErrProductAlreadyExists ...
	ErrProductAlreadyExists = Sentinel("the product already exists")
	//ErrTransaction ...
	ErrTransaction = Sentinel("transaction start error")
	// ErrOrderAlreadyExists ...
	ErrOrderAlreadyExists = Sentinel("the order already exists")
	// ErrOrderNotFound ...
	ErrOrderNotFound = Sentinel("the order not found")
)

const (
	// PostgreSQL xatolik kodlari
	PGForeignKeyViolationCode   = "23503" // Tashqi kalit xatosi
	PGCheckViolationCode        = "23514" // Tekshiruv cheklovi xatosi
	PGUniqueKeyViolationCode    = "23505" // Unikal kalit xatosi
	PGNotNullViolationCode      = "23502" // Null bo'lishi mumkin emas xatosi
	PGDuplicateKeyViolationCode = "23505" // Takroriy kalit xatosi
	PGDataExceptionCode         = "22000" // Ma'lumotlar xatosi
)

var (
	ErrForeignKeyViolation      = errors.New("foreign key violation")
	ErrCheckConstraintViolation = errors.New("check constraint violation")
	ErrUniqueKeyViolation       = errors.New("unique key violation")
	ErrNotNullViolation         = errors.New("not null violation")
	ErrDuplicateKeyViolation    = errors.New("duplicate key violation")
	ErrDataException            = errors.New("data exception")
	ErrNotFound                 = errors.New("record not found")
	ErrRowsAffectedIsZero       = errors.New("no rows affected")
	ErrInvalidCredentials       = errors.New("invalid password")
	ErrUserNotFound             = errors.New("user not found")
)
