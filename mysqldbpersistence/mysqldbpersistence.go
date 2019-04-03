package mysqldbpersistence

import (
	"errors"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/zibilal/sqltyping"
)

type MysqldbPersistence struct {
	dbName       string
	hasCommitted bool
	hasRollback  bool
	dbConnector  *sqlx.DB
	dbTx         *sqlx.Tx
}

func NewMysqldbPersistence(dbName string, dbConnector *sqlx.DB) *MysqldbPersistence {
	p := new(MysqldbPersistence)
	p.dbName = dbName
	p.dbConnector = dbConnector

	return p
}

func (m *MysqldbPersistence) BeginTx() (err error) {
	m.dbTx, err = m.dbConnector.Beginx()
	return
}

func (m *MysqldbPersistence) Commit() error {
	if m.dbTx != nil {
		return m.dbTx.Commit()
	}
	return nil
}

func (m *MysqldbPersistence) Rollback() error {
	if m.dbTx != nil {
		return m.dbTx.Rollback()
	}
	return nil
}

func (m *MysqldbPersistence) Store(name string, data interface{}) error {

	typing := sqltyping.NewSqlTyping(sqltyping.InsertQuery)

	ivalue := reflect.Indirect(reflect.ValueOf(data))
	if ivalue.Kind() != reflect.Struct {
		return errors.New("MysqldbPersistence.Store: unsupported type")
	}

	var tmp interface{}
	if reflect.ValueOf(data).Kind() == reflect.Ptr {
		tmp = reflect.ValueOf(data).Elem().Interface()
	} else {
		tmp = data
	}

	queries, err := typing.Typing(tmp)
	if err != nil {
		return err
	}

	for i := 0; i < len(queries); i++ {
		if m.dbTx != nil {
			result, err := m.dbTx.Exec(queries[i])
			if err != nil {
				return err
			}
			insertedId, err := result.LastInsertId()
			if err == nil && insertedId > 0 {
				fId := ivalue.FieldByName("ID")
				if fId.IsValid() {
					setID(fId, insertedId)
				}
			}
		} else {
			result, err := m.dbConnector.Exec(queries[i])
			if err != nil {
				return err
			}
			insertedId, err := result.LastInsertId()
			if err == nil && insertedId > 0 {
				fId := ivalue.FieldByName("ID")
				setID(fId, insertedId)
			}
		}
	}

	return nil
}

func (m *MysqldbPersistence) Update(query interface{}, name string, data interface{}) error {

	typing := sqltyping.NewSqlTyping(sqltyping.UpdateQuery)

	ivalue := reflect.Indirect(reflect.ValueOf(data))
	if ivalue.Kind() != reflect.Struct {
		return errors.New("MysqldbPersistence.Update: unsupported type, data should be of type struct")
	}

	qValue := reflect.Indirect(reflect.ValueOf(query))
	if qValue.Kind() != reflect.Struct {
		return errors.New("MysqldbPersistence.Update: unsupported type, query should be of type struct")
	}

	queries, err := typing.TypingUpdateWithWhereClause(data, query)

	if err != nil {
		return err
	}

	if m.dbTx != nil {
		_, err := m.dbTx.Exec(queries)
		if err != nil {
			return err
		}
	} else {
		_, err := m.dbConnector.Exec(queries)
		if err != nil {
			return err
		}
	}

	return nil
}

// Fetch fetches the data
func (m *MysqldbPersistence) Fetch(query interface{}, name string, output interface{}) error {
	typing := sqltyping.NewSqlTyping(sqltyping.SelectQuery)
	ivalue := reflect.Indirect(reflect.ValueOf(query))
	if ivalue.Kind() != reflect.Struct {
		return errors.New("MysqldbPersistence.Fetch: unsupported type query")
	}

	var tmp interface{}
	if reflect.ValueOf(query).Kind() == reflect.Ptr {
		tmp = reflect.ValueOf(query).Elem().Interface()
	} else {
		tmp = query
	}

	queries, err := typing.Typing(tmp)
	if err != nil {
		return err
	}

	if len(queries) == 0 {
		return errors.New("MysqldbPersistence.Fetch: invalid state, there are no queries generated")
	}

	if m.dbTx != nil {
		err = m.dbTx.Select(output, queries[0])
		if err != nil {
			return err
		}
	} else {
		err = m.dbConnector.Ping()
		if err != nil {
			return err
		}

		err = m.dbConnector.Select(output, queries[0])
		if err != nil {
			return err
		}
	}

	return nil
}

func setID(fId reflect.Value, insertedId int64) {
	switch {
	case fId.Kind() == reflect.Uint64:
		if fId.IsValid() {
			fId.Set(reflect.ValueOf(uint64(insertedId)))
		}
	case fId.Kind() == reflect.Uint32:
		if fId.IsValid() {
			fId.Set(reflect.ValueOf(uint32(insertedId)))
		}
	case fId.Kind() == reflect.Uint16:
		if fId.IsValid() {
			fId.Set(reflect.ValueOf(uint16(insertedId)))
		}
	case fId.Kind() == reflect.Uint8:
		if fId.IsValid() {
			fId.Set(reflect.ValueOf(uint8(insertedId)))
		}
	case fId.Kind() == reflect.Uint:
		if fId.IsValid() {
			fId.Set(reflect.ValueOf(uint(insertedId)))
		}
	case fId.Kind() == reflect.Int64:
		if fId.IsValid() {
			fId.Set(reflect.ValueOf(int64(insertedId)))
		}
	case fId.Kind() == reflect.Int32:
		if fId.IsValid() {
			fId.Set(reflect.ValueOf(int32(insertedId)))
		}
	case fId.Kind() == reflect.Int16:
		if fId.IsValid() {
			fId.Set(reflect.ValueOf(int16(insertedId)))
		}
	case fId.Kind() == reflect.Int8:
		if fId.IsValid() {
			fId.Set(reflect.ValueOf(int8(insertedId)))
		}
	case fId.Kind() == reflect.Int:
		if fId.IsValid() {
			fId.Set(reflect.ValueOf(int(insertedId)))
		}
	case fId.Type().String() == "sql.NullInt64":
		fId.FieldByName("Int64").Set(reflect.ValueOf(insertedId))
		fId.FieldByName("Valid").Set(reflect.ValueOf(true))
	}
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (m *MysqldbPersistence) Ping() error {
	return m.dbConnector.Ping()
}
