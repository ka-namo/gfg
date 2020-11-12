package api

import (
	"database/sql"

	"coding-challenge-go/pkg/api/middleware"
	"coding-challenge-go/pkg/api/product"
	"coding-challenge-go/pkg/api/seller"
	"coding-challenge-go/pkg/config"

	"github.com/gin-gonic/gin"
)

// CreateAPIEngine creates engine instance that serves API endpoints,
// consider it as a router for incoming requests.
func CreateAPIEngine(db *sql.DB, cfg config.ENVConfig) (*gin.Engine, error) {
	r := gin.New()

	r.Use(middleware.APIVersionResolver)

	v1 := r.Group("api/v1")
	v2 := r.Group("api/v2")

	productRepository := product.NewRepository(db)
	sellerRepository := seller.NewRepository(db)

	var emailProvider, smsProvider product.StockChangedNotifier

	if cfg.NotifySMS {
		smsProvider = seller.NewSMSProvider()
	}

	if cfg.NotifyEmail {
		emailProvider = seller.NewEmailProvider()
	}

	productController := product.NewController(
		productRepository,
		productRepository,
		productRepository,
		productRepository,
		productRepository,
		sellerRepository,
		emailProvider,
		smsProvider,
	)

	v1.GET("products", productController.List)
	v1.GET("product", productController.Get)
	v1.POST("product", productController.Post)
	v1.PUT("product", productController.Put)
	v1.DELETE("product", productController.Delete)
	sellerController := seller.NewController(sellerRepository, sellerRepository)
	v1.GET("sellers", sellerController.List)

	// The decision of having same Controller and having different view as per
	// different API version is - as this is the minimal change and easy to maintain
	// and very less duplication. Having Different controllers could be useful, if
	// we have major changes and previous versions are going to be removed in short time,
	// but this is not clear from the requirement, so I am assuming we will maintain 2 versions.
	v2.GET("products", productController.List)
	v2.GET("product", productController.Get)
	v2.POST("product", productController.Post)
	v2.PUT("product", productController.Put)
	v2.DELETE("product", productController.Delete)
	v2.GET("sellers/top10", sellerController.Top10)

	return r, nil
}
