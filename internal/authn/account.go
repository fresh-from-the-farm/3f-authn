package authn

import (
	"errors"
	"github.com/fresh-from-the-farm/authn/internal/db"
	"github.com/fresh-from-the-farm/authn/internal/model"
	"regexp"
	"time"
)

func GetAccount(username string) (model.Account, error) {
	a, err := db.GetAccount(username)
	return a, err
}

func UpdateAccount(account model.Account, update model.Account) error {

	errs := validateAccountUpdate(update)
	if len(errs) > 0 {
		return errs[0]
	}

	t := time.Now().UTC()
	update.Updated = &t
	upd, err := db.UpdateAccount(account, update)
	if err != nil {
		return err
	}
	if upd.Username != update.Username {
		return errors.New("USERNAME MISMATCH")
	}
	return nil
}

func PurgeAccount(account model.Account) error {
	account.Active = false
	return UpdateAccount(account, account)
}

func CreateAccounts(accounts []model.Account) (processed []model.AccountOperationResponse) {

	for _, a := range accounts {
		var p model.AccountOperationResponse
		errs := validateNewAccount(a)
		if len(errs) != 0 {
			p.Account = a
			p.Status = `FAILED`
			p.Message = `NOT A VALID INFO`
			for _, err := range errs {
				p.Errors = append(p.Errors, err.Error())
			}
			processed = append(processed, p)
			continue
		}

		ax, err := GetAccount(a.Username)
		if err != nil && (len(ax.Username) > 0 || !ax.ID.IsZero()) {
			p.Account = a
			p.Status = `FAILED`
			p.Message = `USERNAME NOT AVAILABLE OR ACCOUNT EXISTS`
			p.Errors = append(p.Errors, err.Error())
			processed = append(processed, p)
			continue
		}

		t := time.Now().UTC()
		a.Created = &t
		a.Updated = &t
		a.Active = true
		err = db.AddAccount(a)

		if err != nil {
			p.Account = a
			p.Status = `FAILED`
			p.Message = `USERNAME NOT VALID`
			p.Errors = append(p.Errors, err.Error())
			processed = append(processed, p)
			continue
		}

		p.Account = a
		p.Status = `SUCCESS`
		p.Message = `ACCOUNT CREATED`
		processed = append(processed, p)
	}

	return processed
}

func validateNewAccount(account model.Account) (errs []error) {

	//^(?=.{8,20}$)(?![_.])(?!.*[_.]{2})[a-zA-Z0-9._]+(?<![_.])$
	//└─────┬────┘└───┬──┘└─────┬─────┘└─────┬─────┘ └───┬───┘
	//│         │         │            │           no _ or . at the end
	//│         │         │            │
	//│         │         │            allowed characters
	//│         │         │
	//│         │         no __ or _. or ._ or .. inside
	//│         │
	//│         no _ or . at the beginning
	//│
	//username is 8-20 characters long
	if account.Username == "" || !regexp.MustCompile(`^[a-z0-9_-]{3,16}$`).MatchString(account.Username) {
		errs = append(errs, errors.New("INVALID USER NAME"))
	}

	_, err := GetAccount(account.Username)
	if err == nil {
		errs = append(errs, errors.New("USERNAME NOT AVAILABLE"))
	}

	if account.FirstName == "" || !regexp.MustCompile(`^([a-z]+[,.]?[ ]?|[a-z]+['-]?)+$`).MatchString(account.FirstName) {
		errs = append(errs, errors.New("INVALID FIRST NAME"))
	}

	if account.LastName == "" || !regexp.MustCompile(`^([a-z]+[,.]?[ ]?|[a-z]+['-]?)+$`).MatchString(account.LastName) {
		errs = append(errs, errors.New("INVALID LAST NAME"))
	}

	if account.Email == "" || !regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$").MatchString(account.Email) {
		errs = append(errs, errors.New("INVALID EMAIL ADDRESS"))
	}

	if account.Password == "" || !regexp.MustCompile(`^[a-zA-Z]\w{3,14}$`).MatchString(account.Password) {
		errs = append(errs, errors.New("INVALID PASSWORD /n The password's first character must be a letter, it must contain at least 4 characters and no more than 15 characters and no characters other than letters, numbers and the underscore may be used"))
	}

	return errs
}

func validateAccountUpdate(account model.Account) (errs []error) {

	//^(?=.{8,20}$)(?![_.])(?!.*[_.]{2})[a-zA-Z0-9._]+(?<![_.])$
	//└─────┬────┘└───┬──┘└─────┬─────┘└─────┬─────┘ └───┬───┘
	//│         │         │            │           no _ or . at the end
	//│         │         │            │
	//│         │         │            allowed characters
	//│         │         │
	//│         │         no __ or _. or ._ or .. inside
	//│         │
	//│         no _ or . at the beginning
	//│
	//username is 8-20 characters long
	if account.Username == "" || !regexp.MustCompile(`^[a-z0-9_-]{3,16}$`).MatchString(account.Username) {
		errs = append(errs, errors.New("INVALID USER NAME"))
	}

	if account.FirstName == "" || !regexp.MustCompile(`^([a-z]+[,.]?[ ]?|[a-z]+['-]?)+$`).MatchString(account.FirstName) {
		errs = append(errs, errors.New("INVALID FIRST NAME"))
	}

	if account.LastName == "" || !regexp.MustCompile(`^([a-z]+[,.]?[ ]?|[a-z]+['-]?)+$`).MatchString(account.LastName) {
		errs = append(errs, errors.New("INVALID LAST NAME"))
	}

	if account.Email == "" || !regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$").MatchString(account.Email) {
		errs = append(errs, errors.New("INVALID EMAIL ADDRESS"))
	}

	if account.Password == "" || !regexp.MustCompile(`^[a-zA-Z]\w{3,14}$`).MatchString(account.Password) {
		errs = append(errs, errors.New("INVALID PASSWORD /n The password's first character must be a letter, it must contain at least 4 characters and no more than 15 characters and no characters other than letters, numbers and the underscore may be used"))
	}

	return errs
}
