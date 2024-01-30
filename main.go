package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

// CalcOption will be my enum to track what I'm calculating
type CalcOption int

const (
	Payment CalcOption = iota
	Principal
	Periods
	Interest
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
	case Interest:
		interestResult := getInterest(*principal, *payment, *periods)
		fmt.Println(interestResult)
	}

}

func getInterest(principal int, payment int, periods int) int {
	return -999
}

func getPeriods(principal int, payment int, interest int) int {
	return -999
}

func getPrincipal(payment int, periods int, interest int) int {
	return -999
}

func getPayment(principal int, periods int, interest int) int {
	return -999
}

// whatCalcWe is a function to find which flag is not unset from default -1.
// It will return my enum, to use in a switch statement.
func whatCalcWe(payment *int, principal *int, periods *int, interest *int) CalcOption {
	if *payment < 0 {
		return Payment
	}
	if *principal < 0 {
		return Principal
	}
	if *periods < 0 {
		return Periods
	}
	return Interest
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
