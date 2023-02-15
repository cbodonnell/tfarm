package auth

type LoginCredentials struct {
	Username string
	Password string
}

type LoginResult struct {
	Success bool
	Error   error
}

type ErrInvalidCredentials struct {
	Err error
}

func (e *ErrInvalidCredentials) Error() string {
	return e.Err.Error()
}

type AuthWorker interface {
	NeedsLogin() bool
	IsAuthenticated() bool
	LoginCredentialsChan() chan LoginCredentials
	LoginResultChan() chan LoginResult
}
