package seller

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type SMSProvider struct {
}

func NewSMSProvider() SMSProvider {
	return SMSProvider{}
}

func (sp SMSProvider) StockChanged(sellerUUID string, sellerPhone string, oldStock int, newStock int, product string) {
	log.Print(
		fmt.Sprintf("SMS Warning sent to %s (Phone: %s): %s Product stock changed", sellerUUID, sellerPhone, product),
	)
}
