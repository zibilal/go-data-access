package bootstrap

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zibilal/go-data-access"
	"github.com/zibilal/go-data-access/apptest/service"
	"github.com/zibilal/go-data-access/apptest/service/simple-service"
	"github.com/zibilal/go-data-access/mongoconnector"
	"github.com/zibilal/go-data-access/persistence/mongodbpersistence"
	"sync"
	"time"
)

var (
	bootstrapped *Bootstrap
	once         sync.Once
)

func GetBootstrapped() *Bootstrap {
	once.Do(func() {
		var err error
		bootstrapped, err = Bootstrapping(context.Background())
		if err != nil {
			panic(err)
		}
	})

	return bootstrapped
}

type Bootstrap struct {
	DbConnector connector.Connector
	ServiceMap  map[string]service.Service
}

func Bootstrapping(ctx context.Context) (*Bootstrap, error) {
	bootstrap := new(Bootstrap)
	bootstrap.ServiceMap = make(map[string]service.Service)

	// create mongodb client
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	urlString := "mongodb://localhost:27027,localhost:27037,localhost:27047/?replicaSet=rs0"
	client, err := mongo.Connect(ctx, urlString)
	if err != nil {
		panic(err)
	}

	connectionContext := mongoconnector.NewMongodbConnectorContext(client)

	bootstrap.DbConnector, err = mongoconnector.NewMongodbConnector()
	if err != nil {
		str := fmt.Sprintf("%s:%s", "failed creating MongodbConnector", err.Error())
		panic(str)
	}

	err = bootstrap.DbConnector.Connect(ctx)
	if err != nil {
		str := fmt.Sprintf("%s:%s", "failed connect", err.Error())
		panic(str)
	}

	dataPersistence := mongodbpersistence.NewMongoPersistence(connectionContext, "simpledb")
	bootstrap.ServiceMap["simpleservice"] = simple_service.NewSimpleService(dataPersistence)
	bootstrap.ServiceMap["fetchservice"] = simple_service.NewFetchService(dataPersistence)
	bootstrap.ServiceMap["updateservice"] = simple_service.NewUpdateService(dataPersistence)

	return bootstrap, nil
}
