package errors

type ApiError struct {
   Err error
   Message string
   Code int
}

func (e ApiError) Error() string {
    return e.Message
}
