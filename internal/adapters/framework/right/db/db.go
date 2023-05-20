package db

import (
    "database/sql"
    "log"
    "time"

    sq "github.com/Masterminds/squirrel"
    _ "github.com/go-sql-driver/mysql"
)

type Adapter struct {
    db *sql.DB
}

func NewAdapter(driverName, dataSourceName string) (*Adapter, error) {
    db, err := sql.Open(driverName, dataSourceName)

    if err != nil {
        log.Fatalf("db connection failure: %v", err)
        return nil, err
    }

    //test connection

    err = db.Ping()

    if err != nil {
        log.Fatalf("db ping failure: %v", err)
    }

    return &Adapter{
        db: db,
    }, err
}

func (dba Adapter)  CloseDbConnection() {
    err := dba.db.Close()
    if err != nil {
        log.Fatalf("db close failure: %v", err)
    }

}

func (dba Adapter) AddToHistory(answer int32, operation string) (error) {
    queryString, args, err := sq.Insert("arith_history").
    Columns("date", "answer", "operation").
    Values(time.Now(), answer, operation).
    ToSql()
    if err != nil {
        return err
    }

    _, err = dba.db.Exec(queryString, args)

    if err != nil {
        return err
    }

    return nil
}
