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

func TestDeleteRota(t *testing.T) {
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
	ts.DELETE("/delete/:rotaId", middleware.SessionChecker(), apI.DeleteRota)

	req, err := http.NewRequest("DELETE", "http://localhost/delete/291255583271555078", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	req.Header.Add("cookie", cookie)

	// before we actually execute our api function, we need to expect required DB actions
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM free WHERE rota_id=?").
		WithArgs(291255583271555078).WillReturnResult(sqlmock.NewResult(0, 6))

	mock.ExpectExec("DELETE FROM rota WHERE rota_id=\\? AND openid=\\?").
		WithArgs(291255583271555078, "oj134ltvn555544444_4abcdefgh").
		WillReturnResult(sqlmock.NewResult(0, 1))
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
