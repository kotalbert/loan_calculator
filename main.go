package main

import (
	"flag"
	"fmt"
	"math"
)

// CalcOption will be my enum to track what I'm calculating
type CalcOption int

const (
	Payment CalcOption = iota
	Principal
	Periods
)

func main() {
	payment, principal, periods, interest, _ := parseArguments()

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

func parseArguments() (*float64, *float64, *float64, *float64, *string) {

	payment := flag.Float64("payment", -1, "payment amount")
	principal := flag.Float64("principal", -1, "loan principal")
	periods := flag.Float64("periods", -1, "number of months needed to repay the loan")
	interest := flag.Float64("interest", -1, "loan interest")
	calcType := flag.String("type", "", "type of calculation, must be either 'annuity' or 'diff'")
	flag.Parse()
	// todo: add validation for --type flag:
	// 	it should be either 'annuity' or 'diff'
	// todo: add validation if --type is diff and --payment is not set
	// todo: add validation if --type is diff and all other flags are set
	// todo: add validation if --type is annuity and exactly 4 of the flags are set
	// todo: add validation if --interest is always provided
	// todo: add validation that all the values are positive
	return payment, principal, periods, interest, calcType
}
