package migrations

import (
	"github.com/jmoiron/sqlx"
	migrate "github.com/phaesoo/sqlx-migrate"
)

var InitTables = migrate.Migration{
	ID:   "20210220220000",
	Name: "init db",
	Migrate: func(tx *sqlx.Tx) error {
		if _, err := tx.Exec(`
		CREATE TABLE auth_key (
			id INT PRIMARY KEY AUTO_INCREMENT,
			access_key VARCHAR(63) UNIQUE,
			secret_key VARCHAR(63) UNIQUE,
			user_uuid VARCHAR(32) UNIQUE
		) ENGINE=INNODB;
		`); err != nil {
			return err
		}

		if _, err := tx.Exec(`
		CREATE TABLE path_permission (
			id INT PRIMARY KEY AUTO_INCREMENT,
			path_pattern VARCHAR(63) UNIQUE
		) ENGINE=INNODB;
		`); err != nil {
			return err
		}

		if _, err := tx.Exec(`
		CREATE TABLE auth_key_path_permissions (
			key_id INT,
			perm_id INT,
			FOREIGN KEY(key_id) REFERENCES auth_key(id),
			FOREIGN KEY(perm_id) REFERENCES path_permission(id)
		) ENGINE=INNODB;
		`); err != nil {
			return err
		}

		return nil
	},
	Rollback: func(tx *sqlx.Tx) error {
		if _, err := tx.Exec("SET FOREIGN_KEY_CHECKS=0;"); err != nil {
			return err
		}

		if _, err := tx.Exec(`
		DROP TABLE auth_key;
		`); err != nil {
			return err
		}
		if _, err := tx.Exec(`
		DROP TABLE path_permission;
		`); err != nil {
			return err
		}
		if _, err := tx.Exec(`
		DROP TABLE auth_key_path_permissions;
		`); err != nil {
			return err
		}

		if _, err := tx.Exec("SET FOREIGN_KEY_CHECKS=1;"); err != nil {
			return err
		}
		return nil
	},
}
