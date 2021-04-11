package wallet

import(
	"github.com/umedjj/wallet/pkg/types"
	"errors"
	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone alredy registered")
var ErrAmmountMustBePositive = errors.New("ammount must be greater then zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance ")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("not found")

type Service struct{
	nextAccountID int64
	accounts 	[]*types.Account
	payments 	[]*types.Payment
	favorites   []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
		return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:		s.nextAccountID,
		Phone: 	phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	
	return account, nil
}





func (s *Service) Deposit(accountID int64, ammount types.Money) error {
	if ammount <= 0 {
		return ErrAmmountMustBePositive
	}
	
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
		account = acc
		break
	}
	}
	
	if account == nil {
		return ErrAccountNotFound
	}
	
	account.Balance += ammount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmmountMustBePositive
	}
	
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
		account = acc
		break
	}
	}
	
	if account == nil {
		return nil, ErrAccountNotFound
	}
	
	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}
	
	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID: 		paymentID,
		AccountID: 	accountID,
		Amount: 	amount,
		Category: 	category,
		Status: 	types.PaymentStatusInProgress,
	}
	
	s.payments = append(s.payments, payment)
	return payment, nil
	
}
	

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts{
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}


func (s *Service) Reject(paymentID string) error  {
	var payment_err *types.Payment
	for _, payment:=range s.payments{
		if payment.ID == paymentID{
			payment_err = payment
		}
	}
		if payment_err == nil {
			return ErrPaymentNotFound
		}
			payment_err.Status = types.PaymentStatusFail
			account, err := s.FindAccountByID(payment_err.AccountID)
			if err != nil{
				return nil
			}
			account.Balance+=payment_err.Amount
			return nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments{
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error)	{
	payment, err :=s.FindPaymentByID(paymentID)
	if err != nil {
		return nil,err
	}
	new_paymentID := uuid.New().String()
	new_payment := &types.Payment{
		ID: 		new_paymentID,
		AccountID: 	payment.AccountID,
		Amount: 	payment.Amount,
		Category: 	payment.Category,
		Status: 	payment.Status,
	}
	account, account_err := s.FindAccountByID(new_payment.AccountID)
	if account_err != nil {
		return nil, account_err
	}
	account.Balance-=new_payment.Amount
	s.payments = append(s.payments, new_payment)
	return new_payment, nil
}