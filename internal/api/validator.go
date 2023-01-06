package api

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const (
	// Index values expected from ROW
	name         = 0
	governmentID = 1
	email        = 2
	debtAmount   = 3
	debtDueDate  = 4
	debtID       = 5

	// datetime layout
	layout = "2006-01-02"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

func Validate(row []string) error {
	if err := nameValidator(row[name]); err != nil {
		return err
	}

	if err := governmentIDValidator(row[governmentID]); err != nil {
		return err
	}

	if err := emailValidator(row[email]); err != nil {
		return err
	}

	if err := debtAmountValidator(row[debtAmount]); err != nil {
		return err
	}

	if err := debtDueDateValidator(row[debtDueDate]); err != nil {
		return err
	}

	return nil
}

func nameValidator(name string) error {
	if len(name) == 0 {
		return errors.New("name is empty")
	}
	return nil
}

// simulate government id validation
func governmentIDValidator(value string) error {
	first := value[0]

	for i := 1; i < len(value); i++ {
		if first != value[i] {
			return nil
		}
	}

	return fmt.Errorf("governmentID %s is invalid", value)
}

func emailValidator(value string) error {
	if emailRegex.MatchString(value) {
		return nil
	}

	return errors.New("email value is invalid")
}

func debtAmountValidator(value string) error {
	amount, err := toDebtAmount(value)
	if err != nil {
		return err
	}

	if amount < 0 {
		return fmt.Errorf("debt amount %s value is invalid", value)
	}

	return nil
}

func debtDueDateValidator(value string) error {
	_, err := toDebtDueDate(value)
	if err != nil {
		return fmt.Errorf("debt due date %s is invalid", value)
	}

	return nil
}

func toDebtAmount(value string) (float64, error) {
	amount, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0, fmt.Errorf("debt amount %s format is invalid", value)
	}

	return amount, err
}

func toDebtDueDate(value string) (time.Time, error) {
	dateValue, err := time.Parse(layout, value)
	if err == nil {
		return time.Time{}, nil
	}

	return dateValue, nil
}
