package postgres

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"msig/common/baseconf/types"
)

const migrationsDir = "./repos/migrations"

type DAO struct {
	db *gorm.DB
}

func New(params types.DBParams) (*gorm.DB, error) {
	d, err := newConn(params)
	if err != nil {
		return nil, err
	}

	sqlDB, err := d.DB()
	if err != nil {
		return nil, err
	}

	err = makeMigration(sqlDB, migrationsDir, params.Database, params.Schema)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func newConn(params types.DBParams) (*gorm.DB, error) {

	sqlDB, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", params.User, params.Password, params.Host, params.Port, params.Database))
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	if strings.Contains(defaultTableName, "views") {
	//		defaultTableName = defaultTableName[:len(defaultTableName)-1]
	//	}
	//	//TODO check connect to db
	//	return defaultTableName
	//}

	return db, nil
}

func makeMigration(conn *sql.DB, migrationDir, dbName, schemaName string) (err error) {
	driver, err := postgresMigrate.WithInstance(conn, &postgresMigrate.Config{DatabaseName: dbName, SchemaName: schemaName})
	if err != nil {
		return err
	}

	err = MakeMigration(driver, migrationDir, dbName)
	if err != nil {
		return err
	}

	return nil
}

func MakeMigration(driver database.Driver, migrationDir, dbName string) error {

	mg, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		dbName, driver)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance: %s", err.Error())
	}

	if err := mg.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}
	return nil
}
