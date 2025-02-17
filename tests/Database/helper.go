package Database

import (
	"app/internal/database"
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"time"
)

type TestDatabase struct {
	DbAddress string
	container testcontainers.Container
	DbService *database.MongoDBClient
}

func SetupTestDatabase() *TestDatabase {
	log.Println("setting up test database")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	container, dbAddr, dbService, err := createMongoContainer(ctx)
	if err != nil {
		log.Fatalf("failed to setup test: %v", err)
	}

	return &TestDatabase{
		container: container,
		DbAddress: dbAddr,
		DbService: dbService,
	}
}

func (tdb *TestDatabase) TearDown() {
	log.Println("tearing down test database")
	_ = tdb.container.Terminate(context.Background())
}

func (tdb *TestDatabase) ClearDatabase() {
	log.Println("clearing database")
	_ = tdb.DbService.Clear()
}

func createMongoContainer(ctx context.Context) (testcontainers.Container, string, *database.MongoDBClient, error) {
	var port = "27017/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{port},
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, "", nil, fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return container, "", nil, fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("mongo container ready and running at port: ", p.Port())
	uri := fmt.Sprintf("mongodb://localhost:%s", p.Port())

	database.URI = uri
	database.DbName = "StravaTestDb"
	serviceDb := database.GetDbClient()

	return container, uri, serviceDb, nil
}
