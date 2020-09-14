package api

import (
	"bytes"
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

func TestChooseFree(t *testing.T) {
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
	ts.POST("/chooseFree/:rotaId", middleware.SessionChecker(), apI.ChooseFree)

	oParams := `{"frees": [0,4,17]}`
	req, err := http.NewRequest("POST", "http://localhost/chooseFree/291255583271555078", bytes.NewBufferString(oParams))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cookie", cookie)

	// before we actually execute our api function, we need to expect required DB actions
	rows1 := sqlmock.NewRows([]string{"nick_name"}).AddRow("王小明")

	mock.ExpectQuery("SELECT nick_name FROM person WHERE openid=?").
		WithArgs("oj134ltvn555544444_4abcdefgh").WillReturnRows(rows1)

	rows2 := sqlmock.NewRows([]string{"limit_choose"}).AddRow(3)

	mock.ExpectQuery("SELECT limit_choose FROM rota WHERE rota_id=?").
		WithArgs(291255583271555078).WillReturnRows(rows2)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM free WHERE openid=\\? AND rota_id=\\?").
		WithArgs("oj134ltvn555544444_4abcdefgh", 291255583271555078).
		WillReturnResult(sqlmock.NewResult(0, 3))

	mock.ExpectExec("INSERT INTO free \\(openid, rota_id, free_id\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs("oj134ltvn555544444_4abcdefgh", 291255583271555078, 0).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("INSERT INTO free \\(openid, rota_id, free_id\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs("oj134ltvn555544444_4abcdefgh", 291255583271555078, 4).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO free \\(openid, rota_id, free_id\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs("oj134ltvn555544444_4abcdefgh", 291255583271555078, 17).
		WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	// now we execute our request
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
	assert.JSONEq(t, resp.Body.String(), `{"status": 0}`)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
