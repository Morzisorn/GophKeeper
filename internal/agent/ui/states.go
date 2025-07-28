package ui

type state int

const (
	stateMenuLoggedOut state = iota
	stateMenuLoggedIn
	stateMasterPassword
	stateItemsList
	stateItemTypeSelection
	stateViewItemsByType
	stateAddItem
	stateAddItemName
	stateAddItemData
	//stateAddCredentialLogin
	stateAddCredentialPassword
	stateAddItemError
	stateAddItemSuccess
	stateEditItem
	stateDeleteItem
	stateSignUpLogin
	stateSignUpPassword
	stateSignInLogin
	stateSignInPassword
	stateItemDetails
	stateDecryptError
	stateConfirmDelete
	stateProcessing
	stateSuccess
	stateError
	stateFinished
	stateDeleteSuccess
	stateDeleteError
	stateEditItemName
	stateEditItemData
	stateEditCredentialLogin
	stateEditCredentialPassword
	stateEditTextContent
	stateEditBinaryData
	stateEditCardNumber
	stateEditCardExpiry
	stateEditCardCVV
	stateEditCardHolder
	stateEditSuccess
	stateEditError
	stateMetadataList
	stateAddMetadataKey
	stateAddMetadataValue
	stateEditMetadataValue
	stateConfirmDeleteMetadata
	stateMetadataSuccess
	stateMetadataError
	stateAddCardExpiry
	stateAddCardCVV
	stateAddCardHolder
	stateConfirmLogout
	stateLogoutSuccess
	stateLogoutError
)

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
