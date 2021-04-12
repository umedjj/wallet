package wallet

import (

	"testing"
	"github.com/umedjj/wallet/pkg/types"
	"fmt"
)



type testService struct {
	*Service
}
	
type testAccount struct {
	phone types.Phone
	balance types.Money
	payments []struct {
	amount types.Money
	category types.PaymentCategory
}
}
	
var defaultTestAccount=testAccount {
	phone: "+7999999999",
	balance: 100,
	payments: []struct{
	amount types.Money
	category types.PaymentCategory
	}{{100, "auto"},
	},
}
	
func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can`t register account, erro = %v", err)
	}
	
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can`t deposit account, error = %v", err)
	}
	
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can`t make payment, error = %v", err)
		}
	}
	
	return account, payments, nil
}
func TestService_FindAccountByID_possitive(t *testing.T) {
	svc := &Service{}
	account,err := svc.RegisterAccount("+79888888888")
	if err != nil {
		fmt.Println(err)
		return
	}
	accounts, err := svc.FindAccountByID(account.ID)
	if err != nil{
		if account != accounts {
			t.Error(err)
		}
	}
}

func TestService_FindAccountByID_negative(t *testing.T)  {
	svc := &Service{}
	account,err := svc.RegisterAccount("+79888888888")
	if err != nil {
		fmt.Println(err)
		return
	}
	accounts, err := svc.FindAccountByID(account.ID+1)
	if err != nil{
		if err != ErrAccountNotFound{
			t.Error(accounts)
		}
	}
}

func TestService_Reject_found(t *testing.T) {
	svc := &Service{}

	account,err := svc.RegisterAccount("+79888888888")
	if err != nil {
		fmt.Println(err)
		return
	}
	
	errr := svc.Deposit(account.ID,100)
	if errr != nil {
		fmt.Println(errr)
		return
	}
	payment, er := svc.Pay(account.ID, 10, "auto")
	if er != nil {
		fmt.Println(er)
	}
	errrr := svc.Reject(payment.ID)
	if errrr != nil {
		fmt.Println(errrr)
	}
}


func TestService_Reject_notfound(t *testing.T) {
	svc := &Service{}
	
	err:= svc.Reject("1")
	if err == nil {
		t.Error(ErrPaymentNotFound)
		return 
	}

}


func TestService_Repeat_found(t *testing.T) {
	svc := &Service{}
	account,err := svc.RegisterAccount("+79888888888")
	if err != nil {
		fmt.Println(err)
		return
	}
	errr := svc.Deposit(account.ID,100)
	if errr != nil {
		fmt.Println(errr)
		return
	}
	payment, er := svc.Pay(account.ID, 10, "auto")
	if er != nil {
		fmt.Println(er)
	}
	_,errrr := svc.Repeat(payment.ID)
	if errrr != nil {
		fmt.Println(errrr)
	}
}

