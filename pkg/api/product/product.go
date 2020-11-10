package product

type product struct {
	ProductID  int    `json:"-"`
	UUID       string `json:"uuid"`
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	Stock      int    `json:"stock"`
	SellerUUID string `json:"seller_uuid"`
}

type productV2 struct {
	ProductID int    `json:"-"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Brand     string `json:"brand"`
	Stock     int    `json:"stock"`
	Seller    seller `json:"seller"`
}

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
