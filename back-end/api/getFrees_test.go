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

func TestGetFrees(t *testing.T) {
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
	ts.GET("/rota/:rotaId", middleware.SessionChecker(), apI.GetFrees)

	req, err := http.NewRequest("GET", "http://localhost/rota/291255583271555078", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	req.Header.Add("cookie", cookie)

	// before we actually execute our api function, we need to expect required DB actions
	rows := sqlmock.NewRows([]string{"free_id"}).AddRow(0)
	rows2 := sqlmock.NewRows([]string{"free_id"}).AddRow(0).AddRow(4).AddRow(17)

	mock.ExpectQuery("SELECT free_id FROM free WHERE openid=\\? AND rota_id=\\?").
		WithArgs("oj134ltvn555544444_4abcdefgh", 291255583271555078).WillReturnRows(rows)
	mock.ExpectQuery("SELECT free_id FROM free WHERE openid=\\? AND rota_id=\\?").
		WithArgs("oj134ltvn555544444_4abcdefgh", 291255583271555078).WillReturnRows(rows2)

	// now we execute our request
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
	assert.JSONEq(t, resp.Body.String(), `{"status": 0, "frees": [0, 4, 17]}`)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
