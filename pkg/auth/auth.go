package auth

type ConfigureCredentials struct {
	ClientID     string
	ClientSecret string
}

type ConfigureResult struct {
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
	NeedsConfigure() bool
	IsAuthenticated() bool
	ConfigureCredentialsChan() chan ConfigureCredentials
	ConfigureResultChan() chan ConfigureResult
}
