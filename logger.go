package txx

// Logger defines methods to log unsuspected errors, which might not
// be returned.
type Logger interface {
	CloseConn(err error)
	Rollback(err error)
}

type nopLogger struct{}

func (nopLogger) CloseConn(error) {}
func (nopLogger) Rollback(error)  {}
