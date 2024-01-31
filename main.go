package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
)

// CalcOption will be my enum to track what I'm calculating
type CalcOption int

const (
	Payment CalcOption = iota
	Principal
	Periods
)

func main() {
	payment, principal, periods, interest := parseArguments()

	calcOption := whatCalcWe(payment, principal, periods, interest)

	switch calcOption {
	case Payment:
		paymentResult := getPayment(*principal, *periods, *interest)
		fmt.Println(paymentResult)
	case Principal:
		principalResult := getPrincipal(*payment, *periods, *interest)
		fmt.Println(principalResult)
	case Periods:
		periodsResult := getPeriods(*principal, *payment, *interest)
		fmt.Println(periodsResult)
	}

}

func getPeriods(principal int, payment int, interest int) int {
	i := getMonthlyInterestRate(interest)
	a := float64(payment) // annuity payment
	p := float64(principal)

	n := math.Log(a/(a-i*p)) / math.Log(1+i)

	return int(math.Ceil(n))
}

func getPrincipal(payment int, periods int, interest int) int {
	a := float64(payment)
	n := float64(periods)
	i := getMonthlyInterestRate(interest)

	p := a / ((i * math.Pow(1+i, n)) / (math.Pow(1+i, n) - 1))

	return int(math.Ceil(p))
}

func getPayment(principal int, periods int, interest int) int {
	p := float64(principal)
	i := getMonthlyInterestRate(interest)
	n := float64(periods)

	payment := p * (i * math.Pow(1+i, n)) / (math.Pow(1+i, n) - 1)
	return int(math.Ceil(payment))

}

func getMonthlyInterestRate(interest int) float64 {
	return float64(interest) / (12.0 * 100)
}

// whatCalcWe is a function to find which flag is not unset from default -1.
// It will return my enum, to use in a switch statement.
func whatCalcWe(payment *int, principal *int) CalcOption {
	if *payment < 0 {
		return Payment
	}
	if *principal < 0 {
		return Principal
	}
	return Periods
}

func parseArguments() (*int, *int, *int, *int) {

	err := checkArgsNumber()

	if err != nil {
		log.Fatal(err)
	}

	payment := flag.Int("payment", -1, "payment amount")
	principal := flag.Int("principal", -1, "loan principal")
	periods := flag.Int("periods", -1, "number of months needed to repay the loan")
	interest := flag.Int("interest", -1, "loan interest")

	flag.Parse()
	return payment, principal, periods, interest
}

func checkArgsNumber() error {
	if len(os.Args[1:]) != 6 {
		return errors.New("expecting exactly three flags set")
	}
	return nil
}
