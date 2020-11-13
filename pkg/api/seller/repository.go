package seller

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
)

// NewRepository builds a new DB repo for Seller.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Repository is DB repo.
type Repository struct {
	db *sql.DB
}

// FindByUUID is the DB implementation for the product.SellerFinder.
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

// list is the DB implementation for the ManyFinder.
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

	if rows.Err() != nil {
		return nil, fmt.Errorf("seller.Repository.list: failed to read sql.Rows %w", rows.Err())
	}

	return sellers, nil
}

// top is the DB implementation of TopSellerFinder.
//
// Returns the Sellers who are selling products ordered by count of products
// they have for sale from the largest to the smallest number limited by given limit.
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
			log.Log().Err(err).Msg("seller.Repository: failed to close the sql.Rows")
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

	if rows.Err() != nil {
		return nil, fmt.Errorf("seller.Repository.top: failed to read sql.Rows %w", rows.Err())
	}

	return sellers, nil
}
