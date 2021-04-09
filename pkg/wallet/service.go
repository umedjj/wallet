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
//s.RegisterAccount undefined (тип * Service не имеет поля или метода RegisterAccount)

func (s *Service) FindAccountByID(accountID int64)(*types.Account, error)  {
	var account *types.Account
	for _, acc := range s.accounts {

		if acc.ID == accountID {
			account = acc
			
		}
		
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}

	return  account, nil
}

