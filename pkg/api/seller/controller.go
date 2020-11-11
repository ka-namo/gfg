package seller

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type FinderByUUID interface {
	FindByUUID(uuid string) (*Seller, error)
}

type ManyFinder interface {
	list() ([]*Seller, error)
}

type TopSellerFinder interface {
	top(limit int) ([]*Seller, error)
}

func NewController(repository *Repository) *controller {
	return &controller{
		repository: repository,
	}
}

type controller struct {
	repository *Repository
}

func (pc *controller) List(c *gin.Context) {
	sellers, err := pc.repository.list()

	if err != nil {
		log.Error().Err(err).Msg("Fail to query seller list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query seller list"})
		return
	}

	sellersJson, err := json.Marshal(sellers)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal sellers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal sellers"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", sellersJson)
}

// Top10 gets the array of maximum 10 sellers ordered by count of products they
// have for sale (count of entries in product table) from the largest to the smallest number.
func (pc *controller) Top10(c *gin.Context) {
	sellers, err := pc.repository.top(10)

	if err != nil {
		log.Error().Err(err).Msg("Fail to query seller list")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to query seller list"})
		return
	}

	sellersJson, err := json.Marshal(sellers)

	if err != nil {
		log.Error().Err(err).Msg("Fail to marshal sellers")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fail to marshal sellers"})
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", sellersJson)
}
