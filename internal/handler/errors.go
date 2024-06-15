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

	ErrInvalidArticleID       = "invalid article ID format"
	ErrInvalidArticleData     = "invalid article data"
	ErrForbiddenModify        = "forbidden to modify this article"
	ErrFailedUpdateArticle    = "failed to update article"
	ErrInvalidID              = "invalid ID"
	ErrFailedDeleteArticle    = "failed to delete article"
	ErrFailedRetrieveArticles = "failed to retrieve articles"
	ErrFailedCreateArticle    = "failed to create article"
	ErrArticleNotFound        = "article not found"
)
