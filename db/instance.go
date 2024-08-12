package db

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"xorm.io/xorm"
)

var Repo *xorm.Engine

func ConnectDatabaseXorm(schemaName string) (*xorm.Engine, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || dbName == "" {
		log.Fatal("Variáveis de ambiente do banco de dados não configuradas corretamente")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco de dados: %v", err)
	}

	_, err = engine.Exec(fmt.Sprintf("USE %s", schemaName))
	if err != nil {
		return nil, fmt.Errorf("falha ao usar o schema %s: %v", schemaName, err)
	}

	err = engine.Ping()
	if err != nil {
		return nil, fmt.Errorf("falha ao pingar o banco de dados: %v", err)
	}

	engine.ShowSQL(true)
	Repo = engine

	fmt.Println("Conexão com o banco de dados estabelecida com sucesso")
	return engine, nil
}

func Migrate(engine *xorm.Engine) error {
	tables := []interface{}{
		new(AssetInventory),
		new(ThreatEventCatalog),
		new(ThreatEventAssets),
		new(Frequency),
		new(LinkThreat),
		new(LossHigh),
		new(RiskCalculation),
		new(ControlLibrary),
		new(RiskController),
		new(Implements),
		new(AggregatedStrength),
		new(Propused),
		new(Control),
		new(Relevance),
		new(LossHighTotal),
		new(LossExceedance),
		new(LossHighGranular),
		new(RiskAssessment),
	}

	for _, table := range tables {
		if err := engine.Sync2(table); err != nil {
			return fmt.Errorf("falha ao migrar a tabela %T: %v", table, err)
		}
	}
	return nil
}

func CreateSchemaAndMigrate(schemaName string) error {
	engine, err := ConnectDatabaseXorm("")
	if err != nil {
		return err
	}

	_, err = engine.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName))
	if err != nil {
		return fmt.Errorf("erro ao criar o schema %s: %v", schemaName, err)
	}

	engine, err = ConnectDatabaseXorm(schemaName)
	if err != nil {
		return err
	}

	err = Migrate(engine)
	if err != nil {
		return fmt.Errorf("erro ao rodar as migrations no schema %s: %v", schemaName, err)
	}

	return nil
}
