package product

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// marshalJSON marshals products as per the API version.
// this layer represents views, which differs for different version.
func marshalJSON(c *gin.Context, products interface{}) ([]byte, error) {
	switch c.MustGet("version").(string) {
	case versionV1:
		return json.Marshal(products)
	case versionV2:
		return json.Marshal(hydrateProductsToV2(c, products))
	}

	return nil, errors.New("invalid API version")
}

// hydrateProductsToV2 transforms product/products to its version V2.
func hydrateProductsToV2(c *gin.Context, products interface{}) interface{} {
	switch v := products.(type) {
	case []*product:
		productsV2 := make([]productV2, 0)
		for _, p := range v {
			productsV2 = append(productsV2, productV2{
				UUID:  p.UUID,
				Name:  p.Name,
				Brand: p.Brand,
				Stock: p.Stock,
				Seller: seller{
					UUID: p.SellerUUID,
					Links: links{
						self{
							HRef: fmt.Sprintf("%s/%s/%s", c.Request.Host, "sellers", p.SellerUUID),
						}},
				},
			})
		}

		return productsV2
	case *product:
		return productV2{
			UUID:  v.UUID,
			Name:  v.Name,
			Brand: v.Brand,
			Stock: v.Stock,
			Seller: seller{
				UUID: v.SellerUUID,
				Links: links{
					self{
						HRef: fmt.Sprintf("%s/%s/%s", c.Request.Host, "sellers", v.SellerUUID),
					}},
			},
		}
	}

	return nil
}
