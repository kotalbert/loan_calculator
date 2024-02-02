package main

import (
	"flag"
	"fmt"
	"log"
	"math"
)

// CalcOption will be my enum to track what I'm calculating
type CalcOption int

// todo: include diff type calculation
const (
	Payment CalcOption = iota
	Principal
	Periods
)

// I will use this to test if flag is set or not
const defaultFlagValue = -1

func main() {
	payment, principal, periods, interest, _, err := parseArguments()

	if err != nil {
		log.Fatal(err)
	}

	calcOption := whatCalcWe(payment, principal)

	switch calcOption {
	case Payment:
		paymentResult := getPayment(*principal, *periods, *interest)
		fmt.Printf("Your monthly payment = %d!\n!", paymentResult)
	case Principal:
		principalResult := getPrincipal(*payment, *periods, *interest)
		fmt.Printf("Your loan principal = %d!\n", principalResult)
	case Periods:
		periodsResult := getPeriods(*principal, *payment, *interest)
		outputPeriodsResult(periodsResult)

	}

}

func outputPeriodsResult(periods int) {
	years := periods / 12
	months := periods % 12
	if months == 0 {
		fmt.Printf("It will take %d years to repay this loan!\n", years)
	} else if years == 0 {
		fmt.Printf("It will take %d months to repay this loan!\n", months)
	} else {
		fmt.Printf("It will take %d years and %d months to repay this loan!\n", years, months)
	}
}

func getPeriods(principal float64, payment float64, interest float64) int {
	i := getMonthlyInterestRate(interest)

	n := math.Log(payment/(payment-i*principal)) / math.Log(1+i)

	return int(math.Ceil(n))
}

func getPrincipal(payment float64, periods float64, interest float64) int {
	i := getMonthlyInterestRate(interest)

	p := payment / ((i * math.Pow(1+i, periods)) / (math.Pow(1+i, periods) - 1))

	return int(math.Ceil(p))
}

func getPayment(principal float64, periods float64, interest float64) int {
	i := getMonthlyInterestRate(interest)

	payment := principal * (i * math.Pow(1+i, periods)) / (math.Pow(1+i, periods) - 1)
	return int(math.Ceil(payment))

}

func getMonthlyInterestRate(interest float64) float64 {
	return interest / (12.0 * 100)
}

// whatCalcWe is a function to find which flag is not unset from default -1.
// It will return my enum, to use in a switch statement.
func whatCalcWe(payment *float64, principal *float64) CalcOption {
	if *payment < 0 {
		return Payment
	}
	if *principal < 0 {
		return Principal
	}
	return Periods
}

func parseArguments() (*float64, *float64, *float64, *float64, *string, error) {

	payment := flag.Float64("payment", defaultFlagValue, "payment amount")
	principal := flag.Float64("principal", defaultFlagValue, "loan principal")
	periods := flag.Float64("periods", defaultFlagValue, "number of months needed to repay the loan")
	interest := flag.Float64("interest", defaultFlagValue, "loan interest")
	calcType := flag.String("type", "", "type of calculation, must be either 'annuity' or 'diff'")

	flag.Parse()

	if err := validateTypeFlag(*calcType); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if err := validatePaymentFlag(*calcType, *payment); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if err := validateAllFlagsSetWhenTypeIsDiff(*calcType, *principal, *periods); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if err := validateAllFlagsSetExceptOneWhenTypeIsAnnuity(*calcType, *principal, *periods, *payment, *interest); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if err := validateInterestFlag(*interest); err != nil {
		return nil, nil, nil, nil, nil, err

	}
	if err := validatePositiveFlagValues(*payment, *principal, *periods, *interest); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return payment, principal, periods, interest, calcType, nil
}

// validateTypeFlag validation for --type flag;
//
//	it should be either 'annuity' or 'diff'
func validateTypeFlag(calcType string) error {
	if calcType == "" {
		return fmt.Errorf("type flag is not set")
	}
	if calcType != "annuity" && calcType != "diff" {
		return fmt.Errorf("invalid type flag value: %s", calcType)
	}
	return nil
}

// validatePaymentFlag validation for --payment flag;
//
//	it should not be set when type is diff
func validatePaymentFlag(calcType string, payment float64) error {
	if calcType == "diff" && payment != -1 {
		return fmt.Errorf("payment flag should not be set when type is diff")
	}
	return nil
}

// validateAllFlagsSetWhenTypeIsDiff validation for --type flag == diff;
//
//	all flags should be set, except --payment
func validateAllFlagsSetWhenTypeIsDiff(calcType string, principal float64, periods float64) error {
	if calcType == "diff" && (principal == defaultFlagValue || periods == defaultFlagValue) {
		return fmt.Errorf("all flags, except --payment, should be set, when --type is diff")
	}
	return nil
}

// validateAllFlagsSetExceptOneWhenTypeIsAnnuity validation for --type flag == annuity;
//
//	exactly 4 of 5 flags should be set
func validateAllFlagsSetExceptOneWhenTypeIsAnnuity(calcType string, principal float64, periods float64, payment float64, interest float64) error {
	if calcType == "annuity" {
		setFlags := 0
		if principal != defaultFlagValue {
			setFlags++
		}
		if periods != defaultFlagValue {
			setFlags++
		}
		if payment != defaultFlagValue {
			setFlags++
		}
		if interest != defaultFlagValue {
			setFlags++
		}
		if setFlags != 4 {
			return fmt.Errorf("exactly 4 of 5 flags should be set when --type is annuity")
		}
	}
	return nil

}

func validateInterestFlag(interest float64) error {
	if interest == defaultFlagValue {
		return fmt.Errorf("interest flag is not set")
	}
	return nil
}

// validatePositiveFlagValues validation for all the flags;
//
//	will check if none of the flags has value less than defaultFlagValue.
//	potentially it will not catch edge case when value -1 is passed
func validatePositiveFlagValues(payment float64, principal float64, periods float64, interest float64) error {
	if payment < defaultFlagValue || principal < defaultFlagValue || periods < defaultFlagValue || interest < defaultFlagValue {
		return fmt.Errorf("all the values should be positive")
	}
	return nil
}
