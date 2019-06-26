package models

import (
    "database/sql"
)

type Tx interface {
    Exec(query string, args)
}
