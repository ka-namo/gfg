package product

import (
	"database/sql"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_list(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		offset, limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*product
		wantErr bool
	}{
		{
			name: "Returns all the Products",
			fields: fields{
				db: func() *sql.DB {
					db, m, _ := sqlmock.New()
					rows := sqlmock.NewRows([]string{"id_product", "name", "brand", "stock", "uuid", "uuid"}).
						AddRow(1, "shoes", "nike", 10, "c943dc0a-98bb-47b4-9d1d-056b95d3f064", "e943dc0a-98bb-47b4-9d1d-056b95d3f064").
						AddRow(2, "shirt", "nike", 20, "d943dc0a-98bb-47b4-9d1d-056b95d3f064", "f943dc0a-98bb-47b4-9d1d-056b95d3f064")
					m.ExpectQuery("SELECT").WillReturnRows(rows)

					return db
				}()},
			args: args{offset: 0, limit: 10},
			want: []*product{
				{ProductID: 1, UUID: "e943dc0a-98bb-47b4-9d1d-056b95d3f064", Name: "shoes", Brand: "nike", Stock: 10, SellerUUID: "c943dc0a-98bb-47b4-9d1d-056b95d3f064"},
				{ProductID: 2, UUID: "f943dc0a-98bb-47b4-9d1d-056b95d3f064", Name: "shirt", Brand: "nike", Stock: 20, SellerUUID: "d943dc0a-98bb-47b4-9d1d-056b95d3f064"},
			},
			wantErr: false,
		},
		{
			name: "Returns error",
			fields: fields{
				db: func() *sql.DB {
					db, m, _ := sqlmock.New()
					rows := sqlmock.NewRows([]string{"id_product", "name", "brand", "stock", "uuid", "uuid"}).
						AddRow(1, "shoes", "nike", 10, "c943dc0a-98bb-47b4-9d1d-056b95d3f064", "e943dc0a-98bb-47b4-9d1d-056b95d3f064").
						RowError(0, errors.New("any sql error"))
					m.ExpectQuery("SELECT").WillReturnRows(rows)

					return db
				}()},
			args:    args{offset: 0, limit: 10},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				db: tt.fields.db,
			}

			defer r.db.Close()

			got, err := r.list(tt.args.offset, tt.args.limit)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestRepository_findByUUID(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *product
		wantErr bool
	}{
		{
			name: "Returns all the Products",
			fields: fields{
				db: func() *sql.DB {
					db, m, _ := sqlmock.New()
					rows := sqlmock.NewRows([]string{"id_product", "name", "brand", "stock", "uuid", "uuid"}).
						AddRow(1, "shoes", "nike", 10, "c943dc0a-98bb-47b4-9d1d-056b95d3f064", "e943dc0a-98bb-47b4-9d1d-056b95d3f064")
					m.ExpectQuery("SELECT").WillReturnRows(rows)

					return db
				}()},
			args:    args{uuid: "c943dc0a-98bb-47b4-9d1d-056b95d3f064"},
			want:    &product{ProductID: 1, UUID: "e943dc0a-98bb-47b4-9d1d-056b95d3f064", Name: "shoes", Brand: "nike", Stock: 10, SellerUUID: "c943dc0a-98bb-47b4-9d1d-056b95d3f064"},
			wantErr: false,
		},
		{
			name: "Returns error",
			fields: fields{
				db: func() *sql.DB {
					db, m, _ := sqlmock.New()
					rows := sqlmock.NewRows([]string{"id_product", "name", "brand", "stock", "uuid", "uuid"}).
						AddRow(1, "shoes", "nike", 10, "c943dc0a-98bb-47b4-9d1d-056b95d3f064", "e943dc0a-98bb-47b4-9d1d-056b95d3f064").
						RowError(0, errors.New("any sql error"))
					m.ExpectQuery("SELECT").WillReturnRows(rows)

					return db
				}()},
			args:    args{uuid: "c943dc0a-98bb-47b4-9d1d-056b95d3f064"},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				db: tt.fields.db,
			}

			defer r.db.Close()

			got, err := r.findByUUID(tt.args.uuid)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestRepository_delete(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		p *product
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "deletes the Product",
			fields: fields{
				db: func() *sql.DB {
					db, m, _ := sqlmock.New()
					rows := sqlmock.NewRows([]string{"id_product", "name", "brand", "stock", "uuid", "uuid"}).
						AddRow(1, "shoes", "nike", 10, "c943dc0a-98bb-47b4-9d1d-056b95d3f064", "e943dc0a-98bb-47b4-9d1d-056b95d3f064")
					m.ExpectQuery("DELETE").WillReturnRows(rows)

					return db
				}()},
			args:    args{p: &product{ProductID: 1, UUID: "e943dc0a-98bb-47b4-9d1d-056b95d3f064", Name: "shoes", Brand: "nike", Stock: 10, SellerUUID: "c943dc0a-98bb-47b4-9d1d-056b95d3f064"}},
			wantErr: false,
		},
		{
			name: "Returns error",
			fields: fields{
				db: func() *sql.DB {
					db, m, _ := sqlmock.New()
					rows := sqlmock.NewRows([]string{"id_product", "name", "brand", "stock", "uuid", "uuid"}).
						AddRow(1, "shoes", "nike", 10, "c943dc0a-98bb-47b4-9d1d-056b95d3f064", "e943dc0a-98bb-47b4-9d1d-056b95d3f064").
						RowError(0, errors.New("any sql error"))
					m.ExpectQuery("DELETE").WillReturnRows(rows)

					return db
				}()},
			args:    args{p: &product{ProductID: 1, UUID: "e943dc0a-98bb-47b4-9d1d-056b95d3f064", Name: "shoes", Brand: "nike", Stock: 10, SellerUUID: "c943dc0a-98bb-47b4-9d1d-056b95d3f064"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				db: tt.fields.db,
			}

			defer r.db.Close()

			err := r.delete(tt.args.p)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}

}

func TestRepository_update(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		p *product
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "updates the Product, stocks changed",
			fields: fields{
				db: func() *sql.DB {
					db, m, _ := sqlmock.New()
					rows := sqlmock.NewRows([]string{"id_product", "name", "brand", "stock", "uuid", "uuid"}).
						AddRow(1, "shoes", "nike", 20, "c943dc0a-98bb-47b4-9d1d-056b95d3f064", "e943dc0a-98bb-47b4-9d1d-056b95d3f064")
					m.ExpectQuery("UPDATE").WillReturnRows(rows)

					return db
				}()},
			args:    args{p: &product{ProductID: 1, UUID: "e943dc0a-98bb-47b4-9d1d-056b95d3f064", Name: "shoes", Brand: "nike", Stock: 20, SellerUUID: "c943dc0a-98bb-47b4-9d1d-056b95d3f064"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				db: tt.fields.db,
			}

			defer r.db.Close()

			err := r.update(tt.args.p)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
