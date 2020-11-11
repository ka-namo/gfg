package seller

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func NewEmailProvider() EmailProvider {
	return EmailProvider{}
}

type EmailProvider struct {
}

func (ep EmailProvider) StockChanged(sellerUUID, sellerEmail string, oldStock, newStock int, product string) {
	log.Print(
		fmt.Sprintf("Email Warning sent to %s (Email: %s): %s Product stock changed", sellerUUID, sellerEmail, product),
	)
}
