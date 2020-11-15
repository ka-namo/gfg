package product

import (
	"net/http"

	sellerAPI "coding-challenge-go/pkg/api/seller"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	listPageSize = 10
	versionV1    = "v1"
	versionV2    = "v2"
)

// SellerFinder is a Finder for Seller.
type SellerFinder interface {
	FindByUUID(uuid string) (*sellerAPI.Seller, error)
}

// FinderByUUID is a Finder for Product by UUID.
type FinderByUUID interface {
	findByUUID(uuid string) (*product, error)
}

// ManyFinder is a Finder for many Products with paging.
type ManyFinder interface {
	list(offset int, limit int) ([]*product, error)
}

// Updater is a updater which updates the Product to repository.
type Updater interface {
	update(product *product) error
}

// Inserter inserts the Product to underlying repository.
type Inserter interface {
	insert(product *product) (*product, error)
}

// Deletes the Product from underlying repository
type Deleter interface {
	delete(product *product) error
}

// controller is an HTTP controller handles HTTP requests for Product APIs.
type controller struct {
	deleter          Deleter
	updater          Updater
	inserter         Inserter
	finderByUUID     FinderByUUID
	finder           ManyFinder
	sellerRepository SellerFinder
	emailProvider    StockChangedNotifier
	smsProvider      StockChangedNotifier
}

// NewController builds the Product controller.
func NewController(
	deleter Deleter,
	updater Updater,
	inserter Inserter,
	finderByUUID FinderByUUID,
	finder ManyFinder,
	sellerRepository SellerFinder,
	emailProvider StockChangedNotifier,
	smsProvider StockChangedNotifier,
) *controller {
	return &controller{
		deleter:          deleter,
		updater:          updater,
		inserter:         inserter,
		finderByUUID:     finderByUUID,
		finder:           finder,
		sellerRepository: sellerRepository,
		emailProvider:    emailProvider,
		smsProvider:      smsProvider,
	}
}

// List returns many products as per page and number of results.
func (pc *controller) List(c *gin.Context) {
	request := &struct {
		Page int `form:"page,default=1"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := pc.finder.list((request.Page-1)*listPageSize, listPageSize)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product list"})
		return
	}

	productsJson, err := marshalJSON(c, products)
	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal products")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal products"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", productsJson)
}

// Get returns the Product by id.
func (pc *controller) Get(c *gin.Context) {
	request := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.finderByUUID.findByUUID(request.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	jsonData, err := marshalJSON(c, product)
	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
}

// Post creates and returns the Product.
func (pc *controller) Post(c *gin.Context) {
	request := &struct {
		Name   string `form:"name"`
		Brand  string `form:"brand"`
		Stock  int    `form:"stock"`
		Seller string `form:"seller"`
	}{}

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seller, err := pc.sellerRepository.FindByUUID(request.Seller)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query seller by UUID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query seller by UUID"})
		return
	}

	if seller == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Seller is not found"})
		return
	}

	product := &product{
		// NOTE - removing UUID generation from controller, as it is a responsibility of
		// repository, and will make the controller also testable.
		Name:       request.Name,
		Brand:      request.Brand,
		Stock:      request.Stock,
		SellerUUID: seller.UUID,
	}

	product, err = pc.inserter.insert(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to insert product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to insert product"})
		return
	}

	jsonData, err := marshalJSON(c, product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
}

// Put updates the Product.
func (pc *controller) Put(c *gin.Context) {
	queryRequest := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(queryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.finderByUUID.findByUUID(queryRequest.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	request := &struct {
		Name  string `form:"name"`
		Brand string `form:"brand"`
		Stock int    `form:"stock"`
	}{}

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldStock := product.Stock

	product.Name = request.Name
	product.Brand = request.Brand
	product.Stock = request.Stock

	err = pc.updater.update(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to insert product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to insert product"})
		return
	}

	if oldStock != product.Stock {
		seller, err := pc.sellerRepository.FindByUUID(product.SellerUUID)

		if err != nil {
			log.Error().Err(err).Msg("Fail to query seller by UUID")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query seller by UUID"})
			return
		}

		// Note - The StockChanged signature seems to me wrong, it was expecting product name and it was sending
		// email, so i changed it to incorporate the logging the correct information.
		if pc.emailProvider != nil {
			pc.emailProvider.StockChanged(seller.UUID, seller.Email, oldStock, product.Stock, product.Name)
		}

		if pc.smsProvider != nil {
			pc.smsProvider.StockChanged(seller.UUID, seller.Phone, oldStock, product.Stock, product.Name)
		}
	}

	jsonData, err := marshalJSON(c, product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal product"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonData)
}

// Delete deletes the product.
func (pc *controller) Delete(c *gin.Context) {
	request := &struct {
		UUID string `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := pc.finderByUUID.findByUUID(request.UUID)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query product by uuid")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query product by uuid"})
		return
	}

	if product == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not found"})
		return
	}

	err = pc.deleter.delete(product)

	if err != nil {
		log.Error().Err(err).Msg("Fail to delete product")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
