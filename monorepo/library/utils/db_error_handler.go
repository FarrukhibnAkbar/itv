package utils

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"itv/monorepo/library/constants"
)

// HandleDBError – bu universal baza operatsiyalari uchun xatolikni boshqarish funksiyasi
func HandleDBError(operation string, err error, tx *gorm.DB) error {
	if err != nil {
		var pgErr *pgconn.PgError
		// Agar PostgreSQL xatoligi bo'lsa
		if errors.As(err, &pgErr) {
			// Xatolik kodiga qarab tegishli xatolikni qaytarish
			switch pgErr.Code {
			case constants.PGForeignKeyViolationCode:
				return fmt.Errorf("error in %s operatsiyasida: %w", operation, constants.ErrForeignKeyViolation)
			case constants.PGCheckViolationCode:
				return fmt.Errorf("error in %s operatsiyasida: %w", operation, constants.ErrCheckConstraintViolation)
			case constants.PGNotNullViolationCode:
				return fmt.Errorf("error in %s operatsiyasida: %w", operation, constants.ErrNotNullViolation)
			case constants.PGDuplicateKeyViolationCode:
				return fmt.Errorf("error in %s operatsiyasida: %w", operation, constants.ErrDuplicateKeyViolation)
			case constants.PGDataExceptionCode:
				return fmt.Errorf("error in %s operatsiyasida: %w", operation, constants.ErrDataException)
			default:
				// Agar boshqa PostgreSQL xatoligi bo'lsa, uni umumiy tarzda qaytaramiz
				return fmt.Errorf("error in  %s operatsiyasida: %w", operation, pgErr)
			}
		}

		// Agar o'zgartirilgan satrlar soni 0 bo'lsa (UPDATE yoki DELETE operatsiyasida)
		if tx != nil && tx.RowsAffected == 0 {
			return fmt.Errorf("error in  %s operatsiyasida: %w", operation, err) // Not Found holati
		}

		// Agar boshqa xatolik bo'lsa, uni umumiy tarzda qaytaramiz
		return fmt.Errorf("error in  %s operatsiyasida: %w", operation, err)
	}
	return nil
}
