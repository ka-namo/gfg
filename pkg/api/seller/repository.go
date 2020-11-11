package seller

import (
	"database/sql"

	"github.com/rs/zerolog/log"
)

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *sql.DB
}

func (r *Repository) FindByUUID(uuid string) (*Seller, error) {
	rows, err := r.db.Query("SELECT id_seller, name, email, phone, uuid FROM seller WHERE uuid = ?", uuid)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	seller := &Seller{}

	err = rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)

	if err != nil {
		return nil, err
	}

	return seller, nil
}

func (r *Repository) list() ([]*Seller, error) {
	rows, err := r.db.Query("SELECT id_seller, name, email, phone, uuid FROM seller")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var sellers []*Seller

	for rows.Next() {
		seller := &Seller{}

		err := rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)
		if err != nil {
			return nil, err
		}

		sellers = append(sellers, seller)
	}

	return sellers, nil
}

func (r *Repository) top(limit int) ([]*Seller, error) {
	query := "SELECT id_seller, name, email, phone, uuid FROM seller" +
		" WHERE id_seller IN (SELECT fk_seller FROM product GROUP BY fk_seller ORDER BY COUNT(*) DESC)" +
		" LIMIT ?"

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Log().Err(err).Msg("sellerRepository: failed to close the sql.Rows")
		}
	}()

	var sellers []*Seller

	for rows.Next() {
		seller := new(Seller)

		err := rows.Scan(&seller.SellerID, &seller.Name, &seller.Email, &seller.Phone, &seller.UUID)
		if err != nil {
			return nil, err
		}

		sellers = append(sellers, seller)
	}

	return sellers, nil
}
