package database

import (
	"database/sql"
	"github.com/normegil/dionysos/internal/model"
)

type UserDAO struct {
	DB *sql.DB
}

func (d UserDAO) Load(username string) (model.User, error) {

}
