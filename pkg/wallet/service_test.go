package wallet

import (
	"reflect"
	"testing"

	"github.com/umedjj/wallet/pkg/types"
)

func TestService_FindAccountByID(t *testing.T) {
	
		service := Service{
			accounts: []*types.Account{
				{ ID: 22, Phone: "9999999999", Balance: 1000,},
				{ ID: 32, Phone: "8888888888", Balance: 1000,},		
			},
		}
	
		expected := &types.Account{
			ID: 22, Phone: "999999999", Balance: 1000,	
		}
	
		result, _ := service.FindAccountByID(10)
	
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Invalid Result: Excpected: %v, actual: %v ", expected, result)
		}
	}

