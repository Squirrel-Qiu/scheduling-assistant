package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"schedule/api/internal"
	"schedule/dbb"
	"schedule/middleware"
)

func TestGetRotas(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbInstance := dbb.InitDB(db)

	ts := gin.New()

	apI := &internal.Implement{
		DB:           dbInstance,
		OpenidGetter: testOpenid("oj134ltvn555544444_4abcdefgh"),
	}

	cookie := login(t, dbInstance, mock, ts, apI)

	// create apI with mocked db, request and response to test
	ts.GET("/rotas", middleware.SessionChecker(), apI.GetRotas)

	req, err := http.NewRequest("GET", "http://localhost/rotas", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	req.Header.Add("cookie", cookie)

	// before we actually execute our api function, we need to expect required DB actions
	rows := sqlmock.NewRows([]string{"rota_id", "title", "shift", "limit_choose", "counter"}).
		AddRow("291255583271555074", "一月份值班表", 2, 2, 6).
		AddRow("291255583271555078", "人事部值班表", 4, 6, 7)

	mock.ExpectQuery("SELECT rota_id, title, shift, limit_choose, counter FROM rota WHERE openid=?").
		WithArgs("oj134ltvn555544444_4abcdefgh").WillReturnError(nil)
	mock.ExpectQuery("SELECT rota_id, title, shift, limit_choose, counter FROM rota WHERE openid=?").
		WithArgs("oj134ltvn555544444_4abcdefgh").WillReturnRows(rows)

	// now we execute our request
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
	assert.JSONEq(t, resp.Body.String(), `{"status": 0, "rotas": [
		{"rota_id": "291255583271555074", "title": "一月份值班表", "shift": 2, "limit_choose": 2, "counter": 6},
        {"rota_id": "291255583271555078", "title": "人事部值班表", "shift": 4, "limit_choose": 6, "counter": 7}]}`)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
