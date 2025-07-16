package ui

// Состояния приложения - теперь с отдельными состояниями для sign up/sign in
type state int

const (
	stateMenuLoggedOut state = iota
	stateMenuLoggedIn
	stateMasterPassword
	stateItemsList
	stateAddItem
    stateAddItemName        
   // stateAddItemType
	stateAddItemData
	stateAddCredentialLogin
	stateAddCredentialPassword 
	stateEditItem
	stateDeleteItem
	stateSignUpLogin
	stateSignUpPassword
	stateSignInLogin
	stateSignInPassword
	stateProcessing
	stateSuccess
	stateError
	stateFinished
)

// Методы для определения типа состояния
func (s state) IsAuth() bool {
	return s >= stateSignUpLogin && s <= stateSignInPassword
}

func (s state) IsSignUp() bool {
	return s == stateSignUpLogin || s == stateSignUpPassword
}

func (s state) IsSignIn() bool {
	return s == stateSignInLogin || s == stateSignInPassword
}

func (s state) IsLoginInput() bool {
	return s == stateSignUpLogin || s == stateSignInLogin
}

func (s state) IsPasswordInput() bool {
	return s == stateSignUpPassword || s == stateSignInPassword
}

func (s state) IsLoggedIn() bool {
    return s == stateMenuLoggedIn || s == stateItemsList || s == stateAddItem || s == stateEditItem || s == stateDeleteItem
}

func (s state) IsAddItemInput() bool {
    return s >= stateAddItem && s <= stateAddCredentialPassword
}