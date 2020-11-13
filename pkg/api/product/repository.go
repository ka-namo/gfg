package product

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// New Repository build new DB repo.
func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

// repository is the new DB repo.
type repository struct {
	db *sql.DB
}

// delete is the DB implementation for the Deleter.
func (r *repository) delete(product *product) error {
	rows, err := r.db.Query("DELETE FROM product WHERE uuid = ?", product.UUID)

	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

// insert is the DB implementation for the Inserter.
// NOTE - as uuid is created in repository now, contract has to be changed to let controller
// know the created product with UUID.
func (r *repository) insert(product *product) (*product, error) {
	product.UUID = uuid.New().String()

	rows, err := r.db.Query(
		"INSERT INTO product (id_product, name, brand, stock, fk_seller, uuid) VALUES(?,?,?,?,(SELECT id_seller FROM seller WHERE uuid = ?),?)",
		product.ProductID, product.Name, product.Brand, product.Stock, product.SellerUUID, product.UUID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return product, nil
}

// update is the DB implementation for the Updater.
func (r *repository) update(product *product) error {
	rows, err := r.db.Query(
		"UPDATE product SET name = ?, brand = ?, stock = ? WHERE uuid = ?",
		product.Name, product.Brand, product.Stock, product.UUID,
	)

	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

// list is the DB implementation for the ManyFinder.
func (r *repository) list(offset int, limit int) ([]*product, error) {
	rows, err := r.db.Query(
		"SELECT p.id_product, p.name, p.brand, p.stock, s.uuid, p.uuid FROM product p "+
			"INNER JOIN seller s ON(s.id_seller = p.fk_seller) LIMIT ? OFFSET ?",
		limit, offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []*product

	for rows.Next() {
		product := &product{}

		err = rows.Scan(&product.ProductID, &product.Name, &product.Brand, &product.Stock, &product.SellerUUID, &product.UUID)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("product.Repository.list: failed to read sql.Rows: %w", rows.Err())
	}

	return products, nil
}

// findByUUID is the DB implementation for the FinderByUUID.
func (r *repository) findByUUID(uuid string) (*product, error) {
	rows, err := r.db.Query(
		"SELECT p.id_product, p.name, p.brand, p.stock, s.uuid, p.uuid FROM product p "+
			"INNER JOIN seller s ON(s.id_seller = p.fk_seller) WHERE p.uuid = ?",
		uuid,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	product := &product{}

	err = rows.Scan(&product.ProductID, &product.Name, &product.Brand, &product.Stock, &product.SellerUUID, &product.UUID)

	if err != nil {
		return nil, err
	}

	return product, nil
}
