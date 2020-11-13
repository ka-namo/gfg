package seller

import (
	"database/sql"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_top(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Seller
		wantErr bool
	}{
		{
			name: "Returns top sellers limited by given limit",
			fields: fields{
				db: func() *sql.DB {
					db, m, _ := sqlmock.New()
					rows := sqlmock.NewRows([]string{"id_seller", "name", "email", "phone", "uuid"}).
						AddRow(3, "james", "j@ex.com", "323-23423-3", "c943dc0a-98bb-47b4-9d1d-056b95d3f064").
						AddRow(2, "mark", "m@ex.com", "789-23423-3", "d943dc0a-98bb-47b4-9d1d-056b95d3f064")
					m.ExpectQuery("SELECT").WillReturnRows(rows)

					return db
				}()},
			args: args{limit: 10},
			want: []*Seller{
				{3, "c943dc0a-98bb-47b4-9d1d-056b95d3f064", "james", "j@ex.com", "323-23423-3"},
				{2, "d943dc0a-98bb-47b4-9d1d-056b95d3f064", "mark", "m@ex.com", "789-23423-3"},
			},
			wantErr: false,
		},
		{
			name: "Returns error when error occurred during rows.Next()",
			fields: fields{
				db: func() *sql.DB {
					db, m, _ := sqlmock.New()
					rows := sqlmock.NewRows([]string{"id_seller", "name", "email", "phone", "uuid"}).
						AddRow(3, "james", "j@ex.com", "323-23423-3", "c943dc0a-98bb-47b4-9d1d-056b95d3f064").
						RowError(0, errors.New("sql error"))
					m.ExpectQuery("SELECT").WillReturnRows(rows)

					return db
				}()},
			args:    args{limit: 10},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				db: tt.fields.db,
			}

			defer r.db.Close()

			got, err := r.top(tt.args.limit)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestRepository_list(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		limit int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Seller
		wantErr bool
	}{
		{
			name: "Returns all the Sellers",
			fields: fields{
				db: func() *sql.DB {
					db, m, _ := sqlmock.New()
					rows := sqlmock.NewRows([]string{"id_seller", "name", "email", "phone", "uuid"}).
						AddRow(3, "james", "j@ex.com", "323-23423-3", "c943dc0a-98bb-47b4-9d1d-056b95d3f064").
						AddRow(2, "mark", "m@ex.com", "789-23423-3", "d943dc0a-98bb-47b4-9d1d-056b95d3f064")
					m.ExpectQuery("SELECT").WillReturnRows(rows)

					return db
				}()},
			args: args{limit: 10},
			want: []*Seller{
				{3, "c943dc0a-98bb-47b4-9d1d-056b95d3f064", "james", "j@ex.com", "323-23423-3"},
				{2, "d943dc0a-98bb-47b4-9d1d-056b95d3f064", "mark", "m@ex.com", "789-23423-3"},
			},
			wantErr: false,
		},
		{
			name: "Returns error when error occurred during rows.Next()",
			fields: fields{
				db: func() *sql.DB {
					db, m, _ := sqlmock.New()
					rows := sqlmock.NewRows([]string{"id_seller", "name", "email", "phone", "uuid"}).
						AddRow(3, "james", "j@ex.com", "323-23423-3", "c943dc0a-98bb-47b4-9d1d-056b95d3f064").
						RowError(0, errors.New("sql error"))
					m.ExpectQuery("SELECT").WillReturnRows(rows)

					return db
				}()},
			args:    args{limit: 10},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				db: tt.fields.db,
			}

			defer r.db.Close()

			got, err := r.list()
			assert.Equal(t, tt.wantErr, err != nil)
			assert.EqualValues(t, tt.want, got)
		})
	}
}
