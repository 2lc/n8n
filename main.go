package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Candidato struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Nome     string    `json:"nome"`
	CPF      string    `json:"cpf"`
	DataNasc time.Time `json:"data_nascimento"`
	NomeMae  string    `json:"nome_mae"`
	NomePai  string    `json:"nome_pai"`
	Endereco string    `json:"endereco"`
	Cidade   string    `json:"cidade"`
	Uf       string    `gorm:"size:2" json:"uf"`
	Cep      string    `json:"cep"`
	Telefone string    `json:"telefone"`
	Email    string    `json:"email"`
	CriadoEm time.Time `gorm:"autoCreateTime" json:"criado_em"`
}

var (
	db  *gorm.DB
	err error
)

// Inicializa a conexão com o banco
func initDB() {

	dsn := "host=192.168.200.201 user=postgres dbname=ALTERDATA password=#abc123# port=5432 sslmode=disable TimeZone=America/Sao_Paulo"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco: %v", err)
	}

	// Cria a tabela se não existir
	if err := db.AutoMigrate(&Candidato{}); err != nil {
		log.Fatalf("❌ Erro ao migrar tabela: %v", err)
	}

	log.Println("✅ Conexão com banco estabelecida e tabela pronta.")
}

// Handler do webhook
func webhookHandler(c *fiber.Ctx) error {
	var novo Candidato

	// Faz o parse do JSON recebido
	if err := c.BodyParser(&novo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"erro":    "JSON inválido",
			"detalhe": err.Error(),
		})
	}

	// Insere no banco
	if err := db.Create(&novo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"erro":    "Falha ao gravar no banco",
			"detalhe": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"mensagem": "Cadastro recebido e salvo com sucesso!",
		"id":       novo.ID,
	})
}

func main() {
	initDB()

	app := fiber.New()

	app.Post("/webhook/candidato", webhookHandler)

	if err := app.Listen(":3080"); err != nil {
		log.Fatal(err)
	}
}
