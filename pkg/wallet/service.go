package wallet

import(
	"github.com/umedjj/wallet/pkg/types"
	"errors"
	"github.com/google/uuid"
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	}

	var ErrPhoneRegistered = errors.New("Phone already Registered")
	var ErrAmountMustBePositive= errors.New("Amount must be greater than 0")
	var ErrAccountNotFound= errors.New("Account not found")
	var ErrNotEnoughBalance = errors.New("Balance not enough")


func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount<=0{
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID==accountID{
			account =acc
			break
		}			
	}
	if account==nil{
		return nil, ErrAccountNotFound
	}
	if account.Balance<amount{
		return nil, ErrNotEnoughBalance
	}
	account.Balance-=amount
	paymentID:=uuid.New().String()
	payment:=&types.Payment{
		ID:       	paymentID,
		AccountID:	accountID,
		Amount:   	amount,
		Category: 	category,
		Status:   	types.PaymentStatusInProgress,
	}

	s.payments = append(s.payments, payment)
	return payment, nil
}



func (s *Service) FindAccountByID(accountID int64)(*types.Account, error)  {
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID ==  accountID{       // платёж найден
			account = acc	
		}
	}	

	if account == nil {    // платёж не найден
	return nil, ErrAccountNotFound
	}
	
	return account, nil
}
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrAccountNotFound
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil

}
