package main

import (
	"context"
	"financial-backend/internal/usecases/dashboard"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"financial-backend/internal/controllers"
	"financial-backend/internal/events"
	_ "financial-backend/internal/events"
	"financial-backend/internal/gateways"
	budgetRepo "financial-backend/internal/repositories/budget"
	budgetMovementRepo "financial-backend/internal/repositories/budget_movement"
	expenseRepo "financial-backend/internal/repositories/expense"
	incomeRepo "financial-backend/internal/repositories/income"
	budgetUseCase "financial-backend/internal/usecases/budget"
	budgetMovementUseCase "financial-backend/internal/usecases/budget_movement"
	expenseUseCase "financial-backend/internal/usecases/expense"
	incomeUseCase "financial-backend/internal/usecases/income"
	"financial-backend/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Configuração do logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Carrega as configurações do ambiente
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Inicializa a conexão com o banco de dados
	db, err := config.GetDatabase()
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	eventPublisher := config.GetPublisher()

	// Inicializa os repositórios
	expenseRepository := expenseRepo.NewRepository(db)
	incomeRepository := incomeRepo.NewRepository(db)
	budgetRepository := budgetRepo.NewRepository(db)
	budgetMovementRepository := budgetMovementRepo.NewRepository(db)

	// Inicializa os gateways
	expenseGateway := gateways.NewExpenseGateway(expenseRepository)
	incomeGateway := gateways.NewIncomeGateway(incomeRepository)
	budgetGateway := gateways.NewBudgetGateway(budgetRepository)
	budgetMovementGateway := gateways.NewBudgetMovementGateway(budgetMovementRepository)

	// Inicializa os casos de uso
	expenseUC := expenseUseCase.NewUseCase(expenseGateway, budgetGateway, eventPublisher, cfg.DefaultDueDate)
	incomeUC := incomeUseCase.NewUseCase(incomeGateway)
	budgetUC := budgetUseCase.NewUseCase(budgetGateway)
	budgetMovementUC := budgetMovementUseCase.NewBudgetMovementUseCase(budgetMovementGateway, budgetGateway, expenseGateway)
	dashboardUC := dashboard.NewDashBoardUseCase(expenseGateway, incomeGateway, budgetMovementGateway)

	// Inicializa os controllers
	expenseController := controllers.NewExpenseController(expenseUC)
	incomeController := controllers.NewIncomeController(incomeUC)
	budgetController := controllers.NewBudgetController(budgetUC)
	budgetMovementController := controllers.NewBudgetMovementController(budgetMovementUC)
	dashboardController := controllers.NewDashboardController(dashboardUC)

	//register handlers
	eventPublisher.RegisterHandler(events.NewExpenseCreatedHandler(db, budgetMovementUC))

	// Configura o router
	router := gin.Default()

	// Middleware de CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Configura as rotas
	api := router.Group("/api")
	{
		// Registra as rotas de cada controller
		expenseController.RegisterRoutes(api)
		incomeController.RegisterRoutes(api)
		budgetController.RegisterRoutes(api)
		budgetMovementController.RegisterRoutes(api)
		dashboardController.RegisterRoutes(api)
	}

	// Configura o servidor HTTP
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	// Inicia o servidor em uma goroutine separada
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	}()

	// Configura o canal para capturar sinais de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Desligando o servidor...")

	// Contexto com timeout para o shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Tenta realizar o shutdown gracioso
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar o servidor: %v", err)
	}

	log.Println("Servidor encerrado com sucesso")
}
