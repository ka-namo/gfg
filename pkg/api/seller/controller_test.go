package seller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ManyFinder is an autogenerated mock type for the ManyFinder type
type ManyFinderMock struct {
	mock.Mock
}

// list provides a mock function with given fields:
func (_m *ManyFinderMock) list() ([]*Seller, error) {
	ret := _m.Called()

	var r0 []*Seller
	if rf, ok := ret.Get(0).(func() []*Seller); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Seller)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TopSellerFinder is an autogenerated mock type for the TopSellerFinder type
type TopSellerFinderMock struct {
	mock.Mock
}

// top provides a mock function with given fields: limit
func (_m *TopSellerFinderMock) top(limit int) ([]*Seller, error) {
	ret := _m.Called(limit)

	var r0 []*Seller
	if rf, ok := ret.Get(0).(func(int) []*Seller); ok {
		r0 = rf(limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*Seller)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func Test_controller_List(t *testing.T) {
	type fields struct {
		finder    ManyFinder
		topFinder TopSellerFinder
	}
	tests := []struct {
		name      string
		fields    fields
		expStatus int
		path      string
		expBody   string
	}{
		{
			name: "v1: Returns 200OK",
			fields: fields{
				finder: func() ManyFinder {
					m := new(ManyFinderMock)
					sellers := []*Seller{
						{
							SellerID: 1,
							UUID:     "fd1574eb-920b-4677-b7e0-4768a5e504c0",
							Name:     "shawn",
							Email:    "s@example.com",
							Phone:    "123-23-23",
						},
						{
							SellerID: 2,
							UUID:     "c943dc0a-98bb-47b4-9d1d-056b95d3f064",
							Name:     "peter",
							Email:    "p@example.com",
							Phone:    "456-23-23",
						},
					}
					m.On("list").Return(sellers, nil)
					return m
				}(),
			},
			path:      "/api/v1/sellers",
			expStatus: http.StatusOK,
			expBody:   `[{"uuid":"fd1574eb-920b-4677-b7e0-4768a5e504c0","name":"shawn","email":"s@example.com","phone":"123-23-23"},{"uuid":"c943dc0a-98bb-47b4-9d1d-056b95d3f064","name":"peter","email":"p@example.com","phone":"456-23-23"}]`,
		},
		{
			name: "v1: Returns 500",
			fields: fields{
				finder: func() ManyFinder {
					m := new(ManyFinderMock)
					m.On("list").Return(nil, errors.New("any error from repo"))
					return m
				}(),
			},
			path:      "/api/v1/sellers",
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error": "Fail to query seller list"}`,
		},
	}
	for _, tt := range tests {

		pc := NewController(tt.fields.finder, tt.fields.topFinder)
		r := gin.Default()

		r.GET("/api/v1/sellers", pc.List)

		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			assert.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expStatus, w.Code)
			assert.JSONEq(t, tt.expBody, w.Body.String())
		})
	}
}

func Test_controller_Top10(t *testing.T) {
	type fields struct {
		finder    ManyFinder
		topFinder TopSellerFinder
	}
	tests := []struct {
		name      string
		fields    fields
		expStatus int
		path      string
		expBody   string
	}{
		{
			name: "v2: Returns 200OK",
			fields: fields{
				topFinder: func() TopSellerFinder {
					m := new(TopSellerFinderMock)
					sellers := []*Seller{
						{
							SellerID: 1,
							UUID:     "fd1574eb-920b-4677-b7e0-4768a5e504c0",
							Name:     "shawn",
							Email:    "s@example.com",
							Phone:    "123-23-23",
						},
						{
							SellerID: 2,
							UUID:     "c943dc0a-98bb-47b4-9d1d-056b95d3f064",
							Name:     "peter",
							Email:    "p@example.com",
							Phone:    "456-23-23",
						},
					}
					m.On("top", 10).Return(sellers, nil)
					return m
				}(),
			},
			path:      "/api/v2/sellers/top10",
			expStatus: http.StatusOK,
			expBody:   `[{"uuid":"fd1574eb-920b-4677-b7e0-4768a5e504c0","name":"shawn","email":"s@example.com","phone":"123-23-23"},{"uuid":"c943dc0a-98bb-47b4-9d1d-056b95d3f064","name":"peter","email":"p@example.com","phone":"456-23-23"}]`,
		},
		{
			name: "v2: Returns 500",
			fields: fields{
				topFinder: func() TopSellerFinder {
					m := new(TopSellerFinderMock)
					m.On("top", 10).Return(nil, errors.New("any error from repo"))
					return m
				}(),
			},
			path:      "/api/v2/sellers/top10",
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error": "Fail to query seller list"}`,
		},
	}
	for _, tt := range tests {

		sc := NewController(tt.fields.finder, tt.fields.topFinder)
		r := gin.Default()

		r.GET("/api/v2/sellers/top10", sc.Top10)

		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			assert.NoError(t, err)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expStatus, w.Code)
			assert.JSONEq(t, tt.expBody, w.Body.String())
		})
	}

}
