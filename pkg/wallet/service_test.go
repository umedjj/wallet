package wallet

import (

	"testing"
	"reflect"
	"github.com/umedjj/wallet/pkg/types"

)



func TestService_FindAccountByID_success(t *testing.T)  {
		service := Service{
			accounts: []*types.Account{
				{ ID: 10, Phone: "9929888881", Balance: 1000,},
				{ ID: 15, Phone: "9929999991", Balance: 1500,},		
			},
		}
	
		expected := &types.Account{
			ID: 10, Phone: "9929888881", Balance: 1000,	
		}
	
		result, _ := service.FindAccountByID(10)
	
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Invalid Result: Excpected: %v, actual: %v ", expected, result)
		}
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