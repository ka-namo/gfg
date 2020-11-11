package seller

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// NewEmailProvider builds the EmailProvider with all ite provided dependencies.
func NewEmailProvider() EmailProvider {
	return EmailProvider{}
}

// EmailProvider implements the email sending to given email id.
type EmailProvider struct {
}

// StockChanged sends an Email to Seller to their email id, when a Product Stock is changed.
func (ep EmailProvider) StockChanged(sellerUUID, sellerEmail string, oldStock, newStock int, product string) {
	log.Print(
		fmt.Sprintf("Email Warning sent to %s (Email: %s): %s Product stock changed", sellerUUID, sellerEmail, product),
	)
}
