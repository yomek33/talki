package handler

const (
	ErrInvalidUserData     = "invalid user data"
	ErrCouldNotCreateUser  = "could not create user"
	ErrInvalidUserToken    = "invalid user token"
	ErrUserNotFound        = "user not found"
	ErrInvalidUserUID      = "invalid user ID"
	ErrCouldNotUpdateUser  = "could not update user"
	ErrCouldNotDeleteUser  = "could not delete user"
	ErrInvalidCredentials  = "invalid credentials"
	TokenExpirationMinutes = 60

	ErrInvalidMaterialID       = "invalid material ID format"
	ErrInvalidMaterialData     = "invalid material data"
	ErrForbiddenModify         = "forbidden to modify this material"
	ErrFailedUpdateMaterial    = "failed to update material"
	ErrInvalidID               = "invalid ID"
	ErrFailedDeleteMaterial    = "failed to delete material"
	ErrFailedRetrieveMaterials = "failed to retrieve materials"
	ErrFailedCreateMaterial    = "failed to create material"
	ErrMaterialNotFound        = "material not found"
)
