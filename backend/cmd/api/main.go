package main

import (
	"log"
	"os"

	"github.com/ledger/backend/internal/config"
	"github.com/ledger/backend/internal/models"
	"github.com/ledger/backend/internal/handlers"
	"github.com/ledger/backend/pkg/middleware"
	"github.com/ledger/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Category{},
		&models.Transaction{},
		&models.Budget{},
		&models.Tag{},
		&models.Device{},
		&models.Family{},
		&models.FamilyMember{},
		&models.FamilyTransaction{},
		&models.AssetChange{},
		&models.Reminder{},
		&models.Notification{},
		&models.AAGroup{},
		&models.AAMember{},
		&models.AASettlement{},
		&models.FinancialGoal{},
		&models.Insurance{},
		&models.Backup{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	r := gin.Default()

	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		authHandler := handlers.NewAuthHandler(db, cfg)
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/refresh", authHandler.RefreshToken)

		protected := api.Group("")
		protected.Use(middleware.Auth(cfg.JWTSecret))
		{
			dashboardHandler := handlers.NewDashboardHandler(db)
			protected.GET("/dashboard/stats", dashboardHandler.GetStats)

			userHandler := handlers.NewUserHandler(db)
			protected.GET("/users/me", userHandler.GetMe)
			protected.PUT("/users/me", userHandler.UpdateMe)
			protected.POST("/users/password", userHandler.ChangePassword)

			accountHandler := handlers.NewAccountHandler(db)
			protected.GET("/accounts", accountHandler.List)
			protected.POST("/accounts", accountHandler.Create)
			protected.GET("/accounts/:id", accountHandler.Get)
			protected.PUT("/accounts/:id", accountHandler.Update)
			protected.DELETE("/accounts/:id", accountHandler.Delete)
			protected.GET("/accounts/total/balance", accountHandler.GetTotalBalance)

			categoryHandler := handlers.NewCategoryHandler(db)
			protected.GET("/categories", categoryHandler.List)
			protected.GET("/categories/tree", categoryHandler.GetTree)
			protected.POST("/categories", categoryHandler.Create)
			protected.PUT("/categories/:id", categoryHandler.Update)
			protected.DELETE("/categories/:id", categoryHandler.Delete)

			transactionHandler := handlers.NewTransactionHandler(db)
			protected.GET("/transactions", transactionHandler.List)
			protected.GET("/transactions/by-date", transactionHandler.ListByDate)
			protected.POST("/transactions", transactionHandler.Create)
			protected.GET("/transactions/:id", transactionHandler.Get)
			protected.PUT("/transactions/:id", transactionHandler.Update)
			protected.DELETE("/transactions/:id", transactionHandler.Delete)
			protected.POST("/transactions/batch-delete", transactionHandler.BatchDelete)

			tagHandler := handlers.NewTagHandler(db)
			protected.GET("/tags", tagHandler.List)
			protected.POST("/tags", tagHandler.Create)
			protected.PUT("/tags/:id", tagHandler.Update)
			protected.DELETE("/tags/:id", tagHandler.Delete)

			budgetHandler := handlers.NewBudgetHandler(db)
			protected.GET("/budgets", budgetHandler.List)
			protected.POST("/budgets", budgetHandler.Create)
			protected.GET("/budgets/:id", budgetHandler.Get)
			protected.PUT("/budgets/:id", budgetHandler.Update)
			protected.DELETE("/budgets/:id", budgetHandler.Delete)
			protected.GET("/budgets/:id/progress", budgetHandler.GetProgress)

			assetHandler := handlers.NewAssetHandler(db)
			protected.GET("/assets", assetHandler.List)
			protected.POST("/assets", assetHandler.Create)
			protected.GET("/assets/:id", assetHandler.Get)
			protected.PUT("/assets/:id", assetHandler.Update)
			protected.DELETE("/assets/:id", assetHandler.Delete)
			protected.GET("/assets/summary", assetHandler.GetSummary)

			reminderHandler := handlers.NewReminderHandler(db)
			protected.GET("/reminders", reminderHandler.List)
			protected.POST("/reminders", reminderHandler.Create)
			protected.PUT("/reminders/:id", reminderHandler.Update)
			protected.DELETE("/reminders/:id", reminderHandler.Delete)

			notificationHandler := handlers.NewNotificationHandler(db)
			protected.GET("/notifications", notificationHandler.List)
			protected.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
			protected.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)
			protected.DELETE("/notifications/:id", notificationHandler.Delete)

			searchHandler := handlers.NewSearchHandler(db)
			protected.GET("/search", searchHandler.Search)

			importHandler := handlers.NewImportHandler(db)
			protected.POST("/import/transactions", importHandler.ImportTransactions)

			aaGroupHandler := handlers.NewAAGroupHandler(db)
			protected.GET("/aa-groups", aaGroupHandler.List)
			protected.POST("/aa-groups", aaGroupHandler.Create)
			protected.DELETE("/aa-groups/:id", aaGroupHandler.Delete)
			protected.POST("/aa-groups/:id/expense", aaGroupHandler.AddExpense)
			protected.GET("/aa-groups/:id/settlements", aaGroupHandler.GetSettlements)

			goalHandler := handlers.NewGoalHandler(db)
			protected.GET("/goals", goalHandler.List)
			protected.POST("/goals", goalHandler.Create)
			protected.PUT("/goals/:id", goalHandler.Update)
			protected.POST("/goals/:id/add-amount", goalHandler.AddAmount)
			protected.DELETE("/goals/:id", goalHandler.Delete)

			insuranceHandler := handlers.NewInsuranceHandler(db)
			protected.GET("/insurances", insuranceHandler.List)
			protected.POST("/insurances", insuranceHandler.Create)
			protected.PUT("/insurances/:id", insuranceHandler.Update)
			protected.DELETE("/insurances/:id", insuranceHandler.Delete)
			protected.GET("/insurances/summary", insuranceHandler.GetSummary)

			netWorthHandler := handlers.NewNetWorthHandler(db)
			protected.GET("/net-worth", netWorthHandler.GetNetWorth)

			backupHandler := handlers.NewBackupHandler(db)
			protected.GET("/backup/export", backupHandler.ExportAll)
			protected.POST("/backup/import", backupHandler.ImportAll)
			protected.GET("/backup/list", backupHandler.List)

			reportHandler := handlers.NewReportHandler(db)
			protected.GET("/reports/summary", reportHandler.Summary)
			protected.GET("/reports/trend", reportHandler.Trend)
			protected.GET("/reports/category", reportHandler.ByCategory)
			protected.GET("/reports/account", reportHandler.ByAccount)
			protected.GET("/reports/merchant", reportHandler.ByMerchant)
			protected.GET("/reports/monthly", reportHandler.MonthlyCompare)
			protected.GET("/reports/export", reportHandler.Export)

			familyHandler := handlers.NewFamilyHandler(db)
			protected.GET("/families", familyHandler.List)
			protected.POST("/families", familyHandler.Create)
			protected.GET("/families/:id", familyHandler.Get)
			protected.DELETE("/families/:id", familyHandler.Delete)
			protected.POST("/families/:id/members", familyHandler.AddMember)
			protected.DELETE("/families/:id/members/:member_id", familyHandler.RemoveMember)
			protected.GET("/families/:id/transactions", familyHandler.GetTransactions)

			familyTxHandler := handlers.NewFamilyTransactionHandler(db)
			protected.POST("/families/:id/transactions", familyTxHandler.Create)

			currencyHandler := handlers.NewCurrencyHandler(db)
			protected.GET("/currencies", currencyHandler.ListCurrencies)
			protected.GET("/currencies/rates", currencyHandler.GetRates)
			protected.POST("/currencies/convert", currencyHandler.Convert)

			// CSV Import
			csvHandler := handlers.NewCSVImportHandler(db)
			protected.POST("/import/csv", csvHandler.ImportCSV)

			// Statistics
			statisticsHandler := handlers.NewStatisticsHandler(db)
			protected.GET("/statistics", statisticsHandler.GetStatistics)

			// File Upload
			uploadHandler := handlers.NewUploadHandler(db)
			protected.POST("/upload", uploadHandler.Upload)

			// Cash Flow Projection
			cashFlowHandler := handlers.NewCashFlowHandler(db)
			protected.GET("/cashflow/projection", cashFlowHandler.GetProjection)

			// Budget Alerts
			budgetAlertHandler := handlers.NewBudgetAlertHandler(db)
			protected.GET("/budgets/alerts", budgetAlertHandler.GetAlerts)
		}
	}

	// Serve uploaded files
	r.Static("/api/v1/uploads", "./uploads")

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
