package product

type product struct {
	ProductID  int    `json:"-"`
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	Stock      int    `json:"stock"`
	SellerUUID string `json:"seller_uuid"`
}

// productV2 is the v2 representation of product
type productV2 struct {
	ProductID int    `json:"-"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Brand     string `json:"brand"`
	Stock     int    `json:"stock"`
	Seller    seller `json:"seller"`
}

// seller represents seller used by productV2
type seller struct {
	UUID  string `json:"uuid"`
	Links links  `json:"_links"`
}

type links struct {
	Self self `json:"self"`
}

type self struct {
	HRef string `json:"href"`
}
