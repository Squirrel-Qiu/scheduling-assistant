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
)

func TestGenerate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbInstance := dbb.InitDB(db)

	ts := gin.New()

	apI := &internal.Implement{
		DB: dbInstance,
	}

	// create apI with mocked db, request and response to test
	ts.GET("/generate/:rotaId", apI.Generate)

	req, err := http.NewRequest("GET", "http://localhost/generate/291255583271555078", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}

	// before we actually execute our api function, we need to expect required DB actions
	rows1 := sqlmock.NewRows([]string{"openid"}).AddRow("201701").AddRow("201702").AddRow("201703")
	rows2 := sqlmock.NewRows([]string{"free_id"}).AddRow(14).AddRow(21).AddRow(3)
	rows3 := sqlmock.NewRows([]string{"shift", "counter"}).AddRow(2, 3)
	rows4 := sqlmock.NewRows([]string{"openid"}).AddRow("201701").AddRow("201702").AddRow("201703")
	rows5 := sqlmock.NewRows([]string{"openid", "nick_name"}).
		AddRow("201701", "张三").AddRow("201702", "李四").AddRow("201703", "王明")

	rows6 := sqlmock.NewRows([]string{"openid"}).AddRow("201701")
	rows7 := sqlmock.NewRows([]string{"openid"}).AddRow("201702").AddRow("201703")
	rows8 := sqlmock.NewRows([]string{"openid"}).AddRow("201701").AddRow("201702").AddRow("201703")

	mock.ExpectQuery("SELECT DISTINCT openid from free where rota_id=?").
		WithArgs(291255583271555078).WillReturnRows(rows1)
	mock.ExpectQuery("SELECT free_id from free where rota_id=\\? group by free_id order by count\\(\\*\\)").
		WithArgs(291255583271555078).WillReturnRows(rows2)
	mock.ExpectQuery("SELECT shift,counter FROM rota WHERE rota_id=?").
		WithArgs(291255583271555078).WillReturnRows(rows3)
	mock.ExpectQuery("SELECT DISTINCT openid FROM free WHERE rota_id=?").
		WithArgs(291255583271555078).WillReturnRows(rows4)
	mock.ExpectQuery("SELECT openid, nick_name FROM person WHERE openid IN \\(\\?,\\?,\\?\\)").
		WithArgs("201701", "201702", "201703").WillReturnRows(rows5)

	mock.ExpectQuery("SELECT openid FROM free WHERE rota_id=\\? and free_id=\\?").
		WithArgs(291255583271555078, 14).WillReturnRows(rows6)
	mock.ExpectQuery("SELECT openid FROM free WHERE rota_id=\\? and free_id=\\?").
		WithArgs(291255583271555078, 21).WillReturnRows(rows7)
	mock.ExpectQuery("SELECT openid FROM free WHERE rota_id=\\? and free_id=\\?").
		WithArgs(291255583271555078, 3).WillReturnRows(rows8)


	// now we execute our request
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
	assert.JSONEq(t, resp.Body.String(), `{"status": 0, "interval": [
        {"free_id": 14, "members": ["张三"]},
		{"free_id": 21, "members": ["李四", "王明"]},
		{"free_id": 3, "members": ["张三", "李四", "王明"]}]}`)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
