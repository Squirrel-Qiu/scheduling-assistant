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

func TestGetJoin(t *testing.T) {
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
	ts.GET("/join", middleware.SessionChecker(), apI.GetJoin)

	req, err := http.NewRequest("GET", "http://localhost/join", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	req.Header.Add("cookie", cookie)

	// before we actually execute our api function, we need to expect required DB actions
	rows1 := sqlmock.NewRows([]string{"rota_id"}).
		AddRow("291255583271555074").
		AddRow("291255583271555078")
	rows2 := sqlmock.NewRows([]string{"rota_id", "title"}).
		AddRow("291255583271555074", "一月份值班表").
		AddRow("291255583271555078", "人事部值班表")

	mock.ExpectQuery("SELECT DISTINCT rota_id FROM free WHERE openid=?").
		WithArgs("oj134ltvn555544444_4abcdefgh").WillReturnRows(rows1)
	mock.ExpectQuery("SELECT rota_id, title FROM rota WHERE rota_id IN \\(\\?,\\?\\)").
		WithArgs(291255583271555074, 291255583271555078).WillReturnRows(rows2)

	// now we execute our request
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
	assert.JSONEq(t, resp.Body.String(), `{"status": 0, "joins": [
		{"rota_id":"291255583271555074", "title":"一月份值班表"},
		{"rota_id":"291255583271555078", "title":"人事部值班表"}]}`)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
