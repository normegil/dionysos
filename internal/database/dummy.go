package database

import (
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/normegil/dionysos"
	"github.com/normegil/dionysos/internal/configuration"
	"github.com/normegil/dionysos/internal/security"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func InsertDummyData(db *sql.DB) error {
	storageDAO := &StorageDAO{DB: db}
	itemDAO := &ItemDAO{DB: db}
	userDAO := &UserDAO{DB: db}
	if viper.GetBool(configuration.KeyDummyData.Name) {
		log.Info().Msg("insert dummy data")
		const dummyItemNb = 50

		user, err := security.NewUser("admin", "admin")
		if err != nil {
			return fmt.Errorf("inserting 'admin' user: %w", err)
		}
		if err := userDAO.Insert(*user); nil != err {
			return fmt.Errorf("inserting %s: %w", user.Name, err)
		}

		for i := 0; i < dummyItemNb; i++ {
			storage := dionysos.Storage{Name: gofakeit.Company()}
			if err := storageDAO.Insert(storage); nil != err {
				return fmt.Errorf("inserting %+v: %w", storage, err)
			}
		}

		for i := 0; i < dummyItemNb; i++ {
			item := dionysos.Item{Name: gofakeit.BeerName()}
			if err := itemDAO.Insert(item); nil != err {
				return fmt.Errorf("inserting %+v: %w", item, err)
			}
		}
	}
	return nil
}
