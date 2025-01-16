package main

import (
	"errors"
	"fmt"
)

// state interface
type ATMState interface {
	Authenticate(string, string) error
	WithdrawMoney(float64) error
	// Cancel() error
	// Exit() error
}

// Account struct
type Account struct {
	AccountNumber string
	Pin           string
	Balance       float64
}

// context
type ATMContext struct {
	CurrentState  ATMState
	Accounts      map[string]*Account
	AvailableCash float64
	ActiveAccount *Account
}

func (atm *ATMContext) setState(state ATMState) {
	atm.CurrentState = state
}

// ATMContext constructor
func NewATMContext(initialCash float64, accounts []*Account) *ATMContext {
	accountMap := make(map[string]*Account)

	for _, acnt := range accounts {
		accountMap[acnt.AccountNumber] = acnt
	}

	atm := &ATMContext{
		AvailableCash: initialCash,
		Accounts:      accountMap,
	}

	atm.setState(&IdleState{
		ATM: atm,
	})

	return atm
}

func (atm *ATMContext) Authenticate(accountNumber, pin string) error {
	return atm.CurrentState.Authenticate(accountNumber, pin)
}

func (atm *ATMContext) WithdrawMoney(amnt float64) error {
	return atm.CurrentState.WithdrawMoney(amnt)
}

// concrete states

// IdleState
type IdleState struct {
	ATM *ATMContext
}

func (s *IdleState) Authenticate(accountNum, pin string) error {
	account, exists := s.ATM.Accounts[accountNum]
	if !exists {
		return errors.New("account number doesn't exist")
	}

	if account.Pin != pin {
		return errors.New("PIN is not matching")
	}

	// when authentication is successful
	s.ATM.ActiveAccount = account
	s.ATM.setState(&AuthenticatedState{
		ATM: s.ATM,
	})

	fmt.Println("\nAuthentication successful!")

	return nil
}

func (s *IdleState) WithdrawMoney(amount float64) error {
	return errors.New("please authenticate first")
}

// AuthenticatedState
type AuthenticatedState struct {
	ATM *ATMContext
}

func (s *AuthenticatedState) Authenticate(accountNum, pin string) error {
	return errors.New("already authenticated")
}

func (s *AuthenticatedState) WithdrawMoney(amount float64) error {
	// transition into `WithdrawState` for further processing
	s.ATM.setState(&WithdrawState{
		ATM: s.ATM,
	})

	return s.ATM.CurrentState.WithdrawMoney(amount)
}

// WithdrawState
type WithdrawState struct {
	ATM *ATMContext
}

func (s *WithdrawState) Authenticate(accountNum, pin string) error {
	return errors.New("can't authenticate during a transaction")
}

func (s *WithdrawState) WithdrawMoney(amount float64) error {
	account := s.ATM.ActiveAccount

	if s.ATM.AvailableCash < amount {
		return errors.New("atm machine ran out of money... doesn't have sufficient cash")
	}

	if account.Balance < amount {
		return errors.New("insufficient account balance")
	}

	// deduct the `amount` from both ATM and Account
	s.ATM.AvailableCash -= amount
	account.Balance -= amount

	fmt.Printf("%.2f rupees got successfully deducted from your account", amount)
	fmt.Println()

	return nil
}

// ErrorState
// type ErrorState struct {
// 	ATM    *ATMContext
// 	ErrMsg string
// }

// func (s *ErrorState) Authenticate(accountNumber, pin string) error {
// 	return
// }

func main() {
	accountsList := []*Account{
		{
			AccountNumber: "81975433120",
			Pin:           "2311",
			Balance:       20000,
		},
		{
			AccountNumber: "51253524113",
			Pin:           "1234",
			Balance:       50000,
		},
	}

	atm := NewATMContext(50000, accountsList)

	// err := atm.Authenticate("81975433120", "1234")
	// if err != nil {
	// 	fmt.Printf("Authentication failed: %q", err)
	// 	return
	// }

	err := atm.Authenticate("81975433120", "2311")
	if err != nil {
		fmt.Printf("Authentication failed: %q", err)
		return
	}

	err = atm.WithdrawMoney(60000)
	if err != nil {
		fmt.Printf("Withdrawal failed: %q", err)
		fmt.Println()
		return
	}

	err = atm.WithdrawMoney(21000)
	if err != nil {
		fmt.Printf("Withdrawal failed: %q", err)
		fmt.Println()
		return
	}
}
