package mongodbpersistence

import (
	"errors"
	"github.com/zibilal/logwrapper"

	"gopkg.in/mgo.v2"
)

type MongodbPersistence struct {
	dbSession *mgo.Session
	dbName    string
}

func NewMongodbPersistence(dbSession *mgo.Session, dbName string) *MongodbPersistence {
	p := new(MongodbPersistence)
	p.dbSession = dbSession
	p.dbName = dbName
	return p
}

func (m *MongodbPersistence) Store(name string, data interface{}) error {

	sCopy := m.dbSession.Copy()

	defer sCopy.Close()

	return sCopy.DB(m.dbName).C(name).Insert(data)
}

func (m *MongodbPersistence) Update(query interface{}, name string, data interface{}) error {
	sCopy := m.dbSession.Copy()
	defer sCopy.Close()

	return sCopy.DB(m.dbName).C(name).Update(query, data)
}

func (m *MongodbPersistence) Fetch(query interface{}, name string, output interface{}) error {

	sCopy := m.dbSession.Copy()

	defer sCopy.Close()

	queryMap, ok := query.(map[string]interface{})
	if !ok {
		logwrapper.Error("[MongodbPersistence.Fetch] Query ", query)
		return errors.New("[MongodbPersistence.Fetch]unexpected type")
	}

	err := sCopy.DB(m.dbName).C(name).Find(queryMap).All(output)
	if err != nil {
		logger.Error("[MongodbPersistence.Fetch] Query: ", query, "Error: ", err)
		return err
	}

	return nil
}

// Ping runs a trivial ping command just to get in touch with the server.
func (m *MongodbPersistence) Ping() error {
	return m.dbSession.Ping()
}

// ExtendMongodbPersistence for support mongodb aggregate query
type ExtendMongodbPersistence struct {
	MongodbPersistence
}

func NewExtendMongodbPersistence(dbSession *mgo.Session, dbName string) *ExtendMongodbPersistence {
	p := new(ExtendMongodbPersistence)
	p.dbSession = dbSession
	p.dbName = dbName
	return p
}

func (m *ExtendMongodbPersistence) Aggregate(query interface{}, name string, output interface{}) error {

	sCopy := m.dbSession.Copy()

	defer sCopy.Close()

	queryMap, ok := query.([]map[string]interface{})
	if !ok {
		return errors.New("[MongodbPersistence.Pipe]unexpected type")
	}

	err := sCopy.DB(m.dbName).C(name).Pipe(queryMap).All(output)
	if err != nil {
		return err
	}

	return nil
}
