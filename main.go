package main

import(	
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	) 

	type Canditado struct {
		ID             uint      `gorm:"primaryKey" json:"id"`
		Nome           string    `json:"nome"`
		CPF            string    `json:"cpf"`
		DataNasc       time.Time `json:"data_nascimento"`
		NomeMae        string    `json:"nome_mae"`
		NomePai        string    `json:"nome_pai"`
		Endereco       string    `json:"endereco"`
		Cidade         string    `json:"cidade"`
		Uf             string    `json:"uf"`
		Cep            string    `json:"cep"`
		Telefone       string    `json:"telefone"`
		Email          string    `json:"email"`
		CriadoEm       time.Time `gorm:"autoCreateTime" json:"criado_em"`
	}

	var (
		db  *gorm.DB
		err error
	)	

// Inicializa a conexão com o banco
func initDB() {
	dsn := "host=localhost user=postgres password=123456 dbname=meubanco port=5432 sslmode=disable TimeZone=America/Sao_Paulo"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao banco: %v", err)
	}

	// Cria a tabela se não existir
	if err := db.AutoMigrate(&Cadastro{}); err != nil {
		log.Fatalf("❌ Erro ao migrar tabela: %v", err)
	}

	log.Println("✅ Conexão com banco estabelecida e tabela pronta.")
}

// Handler do webhook
func webhookHandler(c *fiber.Ctx) error {
	var novo Cadastro

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

	app.Post("/webhook/cadastro", webhookHandler)

	log.Println("🚀 Servidor iniciado em http://localhost:8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
