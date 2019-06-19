package authn

import (
	"github.com/fresh-from-the-farm/authn/internal/model"
)

func GetAccessToken(username, password string) (tkn model.AccessToken, isFound bool, err error) {

	account, err := GetAccount(username)
	if err != nil {
		return tkn, true, err
	}

	if account.Username == "" {
		return tkn, false, nil
	}

	if account.Password == password {
		tkn = model.AccessToken{AccessToken: "evfajfvajfqhvfevfiqvfeyefvi"}
		return tkn, true, nil
	}

	return tkn, true, nil
}
