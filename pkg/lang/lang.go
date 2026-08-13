package lang

const (
	ErrorUnexpected   = "unexpected error"
	MessageUnexpected = "an unexpected error occurred"

	ErrorUserNotFound   = "user not found"
	MessageUserNotFound = "user not found"

	ErrorWrongPassword   = "wrong password"
	MessageWrongPassword = "incorrect password"

	ErrorInvalidToken   = "invalid token"
	MessageInvalidToken = "invalid or expired token"

	ErrorMissingToken   = "missing token"
	MessageMissingToken = "authorization token is required"

	ErrorInvalidUsername   = "username is empty"
	MessageInvalidUsername = "username cannot be empty"

	ErrorInvalidPassword   = "password is empty"
	MessageInvalidPassword = "password cannot be empty"

	ErrorUserAlreadyExists   = "user already exists"
	MessageUserAlreadyExists = "this username is already taken"

	ErrorInventoryNotFound   = "inventory not found"
	MessageInventoryNotFound = "inventory item not found"

	ErrorMarketNotFound   = "market not found"
	MessageMarketNotFound = "market listing not found"

	ErrorMaterialNotFound   = "material not found"
	MessageMaterialNotFound = "material not found"

	ErrorMixNotFound   = "mix not found"
	MessageMixNotFound = "mix not found"

	MessageSetForSellSuccessfully = "item listed for sale successfully"

	MessageBuySuccessfully = "purchase completed successfully"

	MessageMixSuccessfully = "materials are mixing"

	MessagePickNewMaterial = "new material discovered! name it before picking"

	MessagePickButRemaining = "mixing not complete yet, please wait"

	MessagePickSuccessfully = "material picked successfully"

	MessageMaterialNamedSuccessfully = "material named successfully"

	MessageRegisterSuccessfully = "registration successful"

	MessageLoginSuccessfully = "login successful"
)
