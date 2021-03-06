package versions

import (
	"database/sql"
	"fmt"
	"github.com/normegil/dionysos/internal/dao/database"
	"github.com/rs/zerolog/log"
)

type Syncer struct {
	db       *sql.DB
	versions Versions
}

func NewSyncer(db *sql.DB) (*Syncer, error) {
	v := NewVersions(db)
	return &Syncer{db: db, versions: v}, nil
}

func (m Syncer) UpgradeAll() error {
	versionDAO := database.VersionDAO{DB: m.db}
	currentVersion, err := versionDAO.CurrentVersion()
	if err != nil {
		return err
	}

	highestVersion := m.versions.HighestVersion()
	if currentVersion > highestVersion {
		return fmt.Errorf("database version '%d' is higher than current server version '%d'", currentVersion, highestVersion)
	} else if currentVersion == highestVersion {
		return nil
	}

	log.Info().Int("from", currentVersion).Int("to", highestVersion).Msg("Upgrade database")
	for i := currentVersion + 1; i <= highestVersion; i++ {
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("upgrading version '%d': opening transaction: %w", i, err)
		}
		log.Info().Int("version", i).Msg("Applying new database version")
		if err := m.versions[i].Apply(); nil != err {
			newErr := fmt.Errorf("applying version '%d': %w", i, err)
			if nil := tx.Rollback(); err != nil {
				return fmt.Errorf("error during transaction rollback: %w (original: %s)", err, newErr.Error())
			}
			return newErr
		}
		if err := versionDAO.Insert(i); nil != err {
			newErr := fmt.Errorf("update version to '%d': %w", i, err)
			if nil := tx.Rollback(); err != nil {
				return fmt.Errorf("error during transaction rollback: %w (original: %s)", err, newErr.Error())
			}
			return newErr
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("upgrading version '%d': committing transaction: %w", i, err)
		}
	}
	log.Info().Int("from", currentVersion).Int("to", highestVersion).Msg("Database upgraded")
	return nil
}
