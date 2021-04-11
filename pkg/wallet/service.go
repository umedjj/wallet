package wallet

import(
	"github.com/umedjj/wallet/pkg/types"
	"errors"
)

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	}

var ErrAccountNotFound =  errors.New("account not found")


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
