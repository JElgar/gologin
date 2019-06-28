package models

type ApiError struct {
   E error
   Message string
   Code int
}

func (e ApiError) Error() string {
    return e.Message
}
