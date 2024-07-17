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

func ConnectDatabase() {
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
		log.Fatal("Falha ao conectar ao banco de dados: ", err)
	}

	err = engine.Ping()
	if err != nil {
		log.Fatal("Falha ao pingar o banco de dados: ", err)
	}

	engine.ShowSQL(true)
	Repo = engine

	fmt.Println("Conexão com o banco de dados estabelecida com sucesso")
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
		new(PropusedDinamic),
		new(ControlDinamic),
		new(RelevanceDinamic),
	}

	for _, table := range tables {
		if err := engine.Sync2(table); err != nil {
			return fmt.Errorf("falha ao migrar a tabela %T: %v", table, err)
		}
	}
	return nil
}
