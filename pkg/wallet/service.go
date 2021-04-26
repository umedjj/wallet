package wallet

import(
	"errors"
	"log"
	"os"
	//"path/filepath"
	"strconv"
	"strings"
	"github.com/google/uuid"
	"github.com/umedjj/wallet/pkg/types"
	"io"
	"sync"
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

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error){
	payment, err :=s.FindPaymentByID(paymentID)
	if err != nil {
		return nil,err
	}
	favoriteID := uuid.New().String()
	favorite := &types.Favorite{
		ID: 		favoriteID,
		AccountID: 	payment.AccountID,
		Name: 		name,
		Amount: 	payment.Amount,
		Category: 	payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite,nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error){
	var favorite *types.Favorite
	for _, fav := range s.favorites{
		if fav.ID == favoriteID {
			favorite=fav
			break
		}
	}
	if favorite==nil {
		return nil, ErrFavoriteNotFound				
	}
	new_paymentID := uuid.New().String()
	new_payment := &types.Payment{
		ID: 		new_paymentID,
		AccountID: 	favorite.AccountID,
		Amount: 	favorite.Amount,
		Category: 	favorite.Category,
		Status: 	types.PaymentStatusInProgress,
	}
	account, account_err := s.FindAccountByID(new_payment.AccountID)
	if account_err != nil {
		return nil, account_err
	}
	account.Balance-=new_payment.Amount
	s.payments = append(s.payments, new_payment)
	return new_payment, nil
	
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

func(s *Service) ExportToFile(path string) error{
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
	}
	defer func ()  {
		err := file.Close()
		if err != nil {
			log.Print(err)
			return
		}
	}()
	result := ""
	for _, account:= range s.accounts {
		collect := strconv.FormatInt(account.ID, 10) + ";" + string(account.Phone) + ";" +strconv.FormatInt(int64(account.Balance),10) + "|"
		result = result + collect
	}
	_, errr := file.Write([]byte(result))
	if err != nil {
		log.Print(errr)
	}
	return errr

}

func (s *Service) ImportFromFile(path string) error{
	file,err := os.Open(path)
	if err != nil{
		log.Print(err)
		return err
	}
	defer func(){
		err := file.Close()
		if err != nil {
			log.Print(err)
			return
		}
	} ()

	content := make([]byte,0)
	buf := make([]byte,4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		content = append(content, buf[:read]...)
	}
	data := strings.Split(string(content), "|")
	for _, accounts := range data {
		if len(accounts)>1{
		account := strings.Split(accounts, ";")
		id, err := strconv.ParseInt(account[0],10,64)
		if err != nil {
			log.Print(err)
		}
		balance,err := strconv.ParseInt(account[2],10,64)
		if err != nil {
			log.Print(err)
		}
		accountt := &types.Account{
			ID: id,
			Phone: types.Phone(account[1]),
			Balance: types.Money(balance),
		}
		s.accounts = append(s.accounts, accountt)
		}
	}
	return err
}

func (s *Service) Export(dir string) error {
	// err := os.Chdir(dir)
	// if err !=nil{
	// 	log.Print(err)
	// 	return err
	// }

	result_accounts := ""
	for _, account:= range s.accounts {
		collect := strconv.FormatInt(account.ID, 10) + ";" + string(account.Phone) + ";" +strconv.FormatInt(int64(account.Balance),10) + "\r\n"
		result_accounts += collect
	}
	if len(result_accounts) > 0 {
		account_adress:=dir+"/accounts.dump"
		accounts_file, err1 := os.Create(account_adress)
		if err1 != nil {
			log.Print(err1)
		return err1
		}
		_, account_error := accounts_file.Write([]byte(result_accounts))
		if account_error != nil {
			log.Print(account_error)
		}
		defer func ()  {
			err := accounts_file.Close()
			if err != nil {
				log.Print(err)
				return
			}
		}()
	}	

	result_payments := ""
	for _, payment:= range s.payments {
		collect := string(payment.ID) + ";" + strconv.FormatInt(payment.AccountID, 10) + ";" +strconv.FormatInt(int64(payment.Amount),10) + ";" +string(payment.Category)+ ";" +string(payment.Status) + "\r\n"
		result_payments += collect
	}
	if len(result_payments) > 0 {
		payment_adress:=dir+"/payments.dump"
		payments_file, err2 := os.Create(payment_adress)
		if err2 != nil {
			log.Print(err2)
		return err2
		}
		_, payment_error := payments_file.Write([]byte(result_payments))
		if payment_error != nil {
			log.Print(payment_error)
		}
		defer func ()  {
			err := payments_file.Close()
			if err != nil {
				log.Print(err)
				return
			}
		}()
	}

	result_favorites := ""
	for _, favorite:= range s.favorites {
		collect := string(favorite.ID) + ";" + strconv.FormatInt(favorite.AccountID,10) + ";" + string(favorite.Name) + ";" +strconv.FormatInt(int64(favorite.Amount),10) + ";" + string(favorite.Category)+ "\r\n"
		result_favorites += collect
	}
	if len(result_favorites) > 0 {
		favorite_adress:=dir+"/favorites.dump"
		favorite_file, err3 := os.Create(favorite_adress)
		if err3 != nil {
			log.Print(err3)
		return err3
		}
		_, favorites_error := favorite_file.Write([]byte(result_favorites))
		if favorites_error != nil {
			log.Print(favorites_error)
		}
		defer func ()  {
			err := favorite_file.Close()
			if err != nil {
				log.Print(err)
				return
			}
		}()
	}
	return nil	
}


func (s *Service) Import(dir string) error {
	// err := os.Chdir(dir)
	// if err != nil {
	// 	log.Print(err)
	// 	return err
	// }

	//accounts
	account_adress:=dir+"/accounts.dump"
	file_accounts,err := os.Open(account_adress)
	if err != nil{
		log.Print(err)
		return err
	}
	defer func(){
		err := file_accounts.Close()
		if err != nil {
			log.Print(err)
			return
		}
	} ()
	content_accounts := make([]byte,0)
	buf := make([]byte,4)
	for {
		read_accounts, err := file_accounts.Read(buf)
		if err == io.EOF {
			content_accounts = append(content_accounts, buf[:read_accounts]...)
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		content_accounts = append(content_accounts, buf[:read_accounts]...)
	}
	data := strings.Split(string(content_accounts), "\r\n")
	for _, accounts := range data {
		if len(accounts)>1{
		account := strings.Split(accounts, ";")
		id, err := strconv.ParseInt(account[0],10,64)
		if err != nil {
			log.Print(err)
		}
		balance,err := strconv.ParseInt(account[2],10,64)
		if err != nil {
			log.Print(err)
		}
		accountt := &types.Account{
			ID: id,
			Phone: types.Phone(account[1]),
			Balance: types.Money(balance),
		}
		s.accounts = append(s.accounts, accountt)
		}
	}

	//payments
	payment_adress:=dir+"/payments.dump"
	file_payments,err := os.Open(payment_adress)
	if err != nil{
		log.Print(err)
		return err
	}
	defer func(){
		err := file_payments.Close()
		if err != nil {
			log.Print(err)
			return
		}
	} ()
	content_payments := make([]byte,0)
	buff := make([]byte,4)
	for {
		read_payment, err := file_payments.Read(buff)
		if err == io.EOF {
			content_payments = append(content_payments, buff[:read_payment]...)
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		content_payments = append(content_payments, buff[:read_payment]...)
	}
	dataa := strings.Split(string(content_payments), "\r\n")
	for _, payments := range dataa {
		if len(payments)>1{
		payment := strings.Split(payments, ";")
		id_account, err := strconv.ParseInt(payment[1],10,64)
		if err != nil {
			log.Print(err)
		}
		amount,err := strconv.ParseInt(payment[2],10,64)
		if err != nil {
			log.Print(err)
		}
		paymentt := &types.Payment{
			ID: payment[0],
			AccountID: id_account,
			Amount: types.Money(amount),
			Category: types.PaymentCategory(payment[3]),
			Status: types.PaymentStatus(payment[4]),
		}
		s.payments = append(s.payments, paymentt)
		}
	}


	//favorites
	favorite_adress:=dir+"/favorites.dump"
	file_favorites,err := os.Open(favorite_adress)
	if err == nil{

	defer func(){
		err := file_favorites.Close()
		if err != nil {
			log.Print(err)
			return
		}
	} ()
	content_favorites := make([]byte,0)
	bufff := make([]byte,4)
	for {
		read_favorite, err := file_favorites.Read(bufff)
		if err == io.EOF {
			content_favorites = append(content_favorites, bufff[:read_favorite]...)
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		content_favorites = append(content_favorites, bufff[:read_favorite]...)
	}
	dataaa := strings.Split(string(content_favorites), "\r\n")
	for _, favorites := range dataaa {
		if len(favorites)>1{
		favorite := strings.Split(favorites, ";")
		id_account, err := strconv.ParseInt(favorite[1],10,64)
		if err != nil {
			log.Print(err)
		}
		amount,err := strconv.ParseInt(favorite[3],10,64)
		if err != nil {
			log.Print(err)
		}
		favoritee := &types.Favorite{
			ID: favorite[0],
			AccountID: id_account,
			Name: favorite[2],
			Amount: types.Money(amount),
			Category: types.PaymentCategory(favorite[4]),
		}
		s.favorites = append(s.favorites,favoritee)
		}
	}
	}
	return nil
}

func (s *Service) ExportAccountHistory(accountID int64) ([]types.Payment, error) {
	var result []types.Payment
	i:=1
	for _, payment:= range s.payments{
		if payment.AccountID == accountID {
			result = append(result, *payment)
			i++
		}
	}
	if i==1{
		return result, ErrAccountNotFound
	} 
	return result, nil
}


func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	result:= ""
	if len(payments)> 0 && len(payments) <=records {
		for _, payment:= range payments {
			collect := string(payment.ID) + ";" + strconv.FormatInt(payment.AccountID, 10) + ";" +strconv.FormatInt(int64(payment.Amount),10) + ";" +string(payment.Category)+ ";" +string(payment.Status) + "\r\n"
			result += collect
		}
		adres:=dir+"/payments.dump"
		file, err := os.Create(adres)
		if err != nil {
			log.Print(err)
		}
		defer file.Close()
		_, file_error := file.Write([]byte(result))
		if file_error != nil {
			log.Print(file_error)
		}
		return nil 
	} else {
	i:=1 //helps for count number files
	j:=0 //helps for comparison with records 
	for _, payment:= range payments {
		collect := string(payment.ID) + ";" + strconv.FormatInt(payment.AccountID, 10) + ";" +strconv.FormatInt(int64(payment.Amount),10) + ";" +string(payment.Category)+ ";" +string(payment.Status) + "\r\n"
		result += collect
			if j==0{
				adress:=dir+"/payments"+strconv.Itoa(i)+".dump" 
				file, err := os.Create(adress)
				if err != nil {
					log.Print(err)
				}
				_, file_error := file.Write([]byte(result))
				if file_error != nil {
					log.Print(file_error)
				}
				file.Close()
				}
				j++
			if j== records {
				i++
				j = 0
				result=""
			}
	}
	return nil
}
}


func (s *Service) SumPayments(goroutines int) types.Money {
	wg:=sync.WaitGroup{}
	mu:=sync.Mutex{}
	summ:=int64(0)
	count:=0
	i:=0
	if goroutines == 0{
		count =len(s.payments)
	} else {
		count = int(len(s.payments)/goroutines)
	}
	for i=0; i<=goroutines-1; i++{
		wg.Add(1)
		go func(number int){
			defer wg.Done()
			val:=int64(0)
			payments:= s.payments[count*number:(count)*(number+1)]
			for _, payment:=range payments{
				val += int64(payment.Amount)
			}
			mu.Lock()
			summ+=val
			mu.Unlock()
		}(i)
	}
	wg.Add(1)
	go func(){
		defer wg.Done()
		val:=int64(0)
		payments:= s.payments[i*count:]
		for _, payment:= range payments{
			val+=int64(payment.Amount)
		}
		mu.Lock()
		summ+=int64(val)
		mu.Unlock()
	}()
	wg.Wait()
	return types.Money(summ)
}


func (s *Service) FilterPayments(accountID int64, goroutines int) ([]types.Payment, error){
	wg:=sync.WaitGroup{}
	mu:=sync.Mutex{}
	var result []types.Payment
	count:=0
	i:=0
	if goroutines == 0{
		count =len(s.payments)
	} else {
		count = int(len(s.payments)/goroutines)
	}
	for i=0; i<=goroutines-1; i++{
		wg.Add(1)
		go func(number int){
			defer wg.Done()
			var pay []types.Payment
			payments:= s.payments[count*number:(count)*(number+1)]
			for _, payment:=range payments{
				if payment.AccountID == accountID{
					pay=append(pay, types.Payment{
						ID: payment.ID,
						AccountID: payment.AccountID,
						Amount: payment.Amount,
						Category: payment.Category,
						Status: payment.Status,
					})				
				}
			}
			mu.Lock()
			result=append(result, pay...)
			mu.Unlock()
		}(i)
	}
	wg.Add(1)
	go func(){
		defer wg.Done()
		var pay []types.Payment
		payments:= s.payments[i*count:]
		for _, payment:= range payments{
			if payment.AccountID == accountID{
				pay=append(pay, types.Payment{
					ID: payment.ID,
					AccountID: payment.AccountID,
					Amount: payment.Amount,
					Category: payment.Category,
					Status: payment.Status,
				})
		}
		}
		mu.Lock()
		result=append(result, pay...)
		mu.Unlock()
	}()
	wg.Wait()
	if len(result) == 0{
		return nil,ErrAccountNotFound
	}
	return result,nil
}

func (s *Service) FilterPaymentsByFn(filter func(payment types.Payment) bool, goroutines int,
	) ([]types.Payment, error){
		wg:=sync.WaitGroup{}
		mu:=sync.Mutex{}
		var result []types.Payment
		count:=0
		i:=0
		if goroutines == 0{
			count =len(s.payments)
		} else {
			count = int(len(s.payments)/goroutines)
		}
		for i=0; i<=goroutines-1; i++{
			wg.Add(1)
			go func(number int){
				defer wg.Done()
				var pay []types.Payment
				payments:= s.payments[count*number:(count)*(number+1)]
				for _, payment:=range payments{
						paym := types.Payment{
							ID: payment.ID,
							AccountID: payment.AccountID,
							Amount: payment.Amount,
							Category: payment.Category,
							Status: payment.Status,
						}
						if filter(paym) {
							pay = append(pay, paym)
						}				
				}
				mu.Lock()
				result=append(result, pay...)
				mu.Unlock()
			}(i)
		}
		wg.Add(1)
		go func(){
			defer wg.Done()
			var pay []types.Payment
			payments:= s.payments[i*count:]
			for _, payment:= range payments{
					paym:=types.Payment{
						ID: payment.ID,
						AccountID: payment.AccountID,
						Amount: payment.Amount,
						Category: payment.Category,
						Status: payment.Status,
					}
					if filter(paym) {
						pay = append(pay, paym)
					}	
			}
			mu.Lock()
			result=append(result, pay...)
			mu.Unlock()
		}()
		wg.Wait()
		if len(result) == 0{
			return nil,ErrAccountNotFound
		}
		return result,nil
}