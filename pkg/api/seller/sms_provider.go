package seller

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// SMSProvider implements the SMS Sending to given phone number.
type SMSProvider struct {
}

// NewSMSProvider builds the SMSProvider with all ite provided dependencies.
func NewSMSProvider() SMSProvider {
	return SMSProvider{}
}

// StockChanged sends an SMS to Seller to their phone, when a Product Stock is changed.
func (sp SMSProvider) StockChanged(sellerUUID string, sellerPhone string, oldStock int, newStock int, product string) {
	log.Print(
		fmt.Sprintf("SMS Warning sent to %s (Phone: %s): %s Product stock changed", sellerUUID, sellerPhone, product),
	)
}
