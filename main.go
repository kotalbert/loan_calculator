package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
)

// AnnuityTypeCalculatedParameter will be my enum to track what I'm calculating, when `--type` flag is set to `annuity`
type AnnuityTypeCalculatedParameter int

const (
	Payment AnnuityTypeCalculatedParameter = iota
	Principal
	Periods
)

type PaymentType int

const (
	Annuity = iota
	Differentiate
)

// I will use this to test if flag is set or not
const defaultFlagValue = -1

func main() {
	payment, principal, periods, interest, paymentType, err := parseArguments()

	if err != nil {
		log.Fatal(err)
	}

	switch *paymentType {
	case Differentiate:
		calculateDifferentiateTypeLoanValues(principal, periods, interest)
	default:
		calculateAnnuityTypeLoanValues(principal, periods, interest, payment)
	}

}

// calculateDifferentiateTypeLoanValues calculates the differentiate payment
//
//	over the time specified by the periods flag
func calculateDifferentiateTypeLoanValues(principal *float64, periods *float64, interest *float64) {
	mir := getMonthlyInterestRate(*interest)
	m := 1
	for i := 0; i < int(*periods); i++ {
		d := getDifferentiatePayment(*principal, *periods, mir, m)
		fmt.Printf("Month %d: payment is %d\n", m, d)
		m++
	}
}

// getDifferentiatePayment calculates the differentiate payment in the given month (m)
func getDifferentiatePayment(principal float64, periods float64, interest float64, m int) interface{} {

}

// calculateAnnuityTypeLoanValues handles the logic for calculating the annuity payment, principal or periods
//
// this logic is applied when `--type` flag is set to `annuity`; it also decides
// if we are calculating payment, principal or periods, based on what flags are set
func calculateAnnuityTypeLoanValues(principal *float64, periods *float64, interest *float64, payment *float64) {

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
	mir := getMonthlyInterestRate(interest)

	periods := math.Log(payment/(payment-mir*principal)) / math.Log(1+mir)

	return int(math.Ceil(periods))
}

func getPrincipal(payment float64, periods float64, interest float64) int {
	mir := getMonthlyInterestRate(interest)

	pmt := payment / ((mir * math.Pow(1+mir, periods)) / (math.Pow(1+mir, periods) - 1))

	return int(math.Ceil(pmt))
}

func getPayment(principal float64, periods float64, interest float64) int {
	mir := getMonthlyInterestRate(interest)

	payment := principal * (mir * math.Pow(1+mir, periods)) / (math.Pow(1+mir, periods) - 1)
	return int(math.Ceil(payment))

}

func getMonthlyInterestRate(interest float64) float64 {
	return interest / (12.0 * 100)
}

// whatCalcWe is a function to find which flag is not unset from default -1.
//
//	It will return my enum, to use in a switch statement.
func whatCalcWe(payment *float64, principal *float64) AnnuityTypeCalculatedParameter {
	if *payment < 0 {
		return Payment
	}
	if *principal < 0 {
		return Principal
	}
	return Periods
}

// parseArguments parses the command line arguments, validates them and returns pointers to the values
//
//	will return error if any of the validation fails
func parseArguments() (*float64, *float64, *float64, *float64, *PaymentType, error) {

	payment := flag.Float64("payment", defaultFlagValue, "payment amount")
	principal := flag.Float64("principal", defaultFlagValue, "loan principal")
	periods := flag.Float64("periods", defaultFlagValue, "number of months needed to repay the loan")
	interest := flag.Float64("interest", defaultFlagValue, "loan interest")
	typeFlagValue := flag.String("type", "", "type of calculation, must be either 'annuity' or 'diff'")

	flag.Parse()

	if err := validateTypeFlag(*typeFlagValue); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if err := validatePaymentFlag(*typeFlagValue, *payment); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if err := validateAllFlagsSetWhenTypeIsDiff(*typeFlagValue, *principal, *periods); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if err := validateAllFlagsSetExceptOneWhenTypeIsAnnuity(*typeFlagValue, *principal, *periods, *payment, *interest); err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if err := validateInterestFlag(*interest); err != nil {
		return nil, nil, nil, nil, nil, err

	}
	if err := validatePositiveFlagValues(*payment, *principal, *periods, *interest); err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// convert string to enum
	var paymentType *PaymentType
	if *typeFlagValue == "diff" {
		*paymentType = Differentiate
	} else {
		*paymentType = Annuity
	}

	return payment, principal, periods, interest, paymentType, nil
}

// validateTypeFlag validation for --type flag;
//
//	it should be either 'annuity' or 'diff'
func validateTypeFlag(calcType string) error {
	if calcType == "" {
		return errors.New("incorrect parameters")
	}
	if calcType != "annuity" && calcType != "diff" {
		return errors.New("incorrect parameters")
	}
	return nil
}

// validatePaymentFlag validation for --payment flag;
//
//	it should not be set when type is diff
func validatePaymentFlag(calcType string, payment float64) error {
	if calcType == "diff" && payment != -1 {
		return errors.New("incorrect parameters")
	}
	return nil
}

// validateAllFlagsSetWhenTypeIsDiff validation for --type flag == diff;
//
//	all flags should be set, except --payment
func validateAllFlagsSetWhenTypeIsDiff(calcType string, principal float64, periods float64) error {
	if calcType == "diff" && (principal == defaultFlagValue || periods == defaultFlagValue) {
		return errors.New("incorrect parameters")
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
		if setFlags != 3 {
			return errors.New("incorrect parameters")
		}
	}
	return nil

}

func validateInterestFlag(interest float64) error {
	if interest == defaultFlagValue {
		return errors.New("incorrect parameters")
	}
	return nil
}

// validatePositiveFlagValues validation for all the flags;
//
//	will check if none of the flags has value less than defaultFlagValue.
//	potentially it will not catch edge case when value -1 is passed
func validatePositiveFlagValues(payment float64, principal float64, periods float64, interest float64) error {
	if payment < defaultFlagValue || principal < defaultFlagValue || periods < defaultFlagValue || interest < defaultFlagValue {
		return errors.New("incorrect parameters")
	}
	return nil
}
