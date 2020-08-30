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

func TestNewRota(t *testing.T) {
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
		RotaidGetter: testRotaid(291255583271555074),
	}

	cookie := login(t, dbInstance, mock, ts, apI)

	// create apI with mocked db, request and response to test
	ts.POST("/newRota", middleware.SessionChecker(), apI.NewRota)

	oParams2 := `{"title": "一月份值班表", "shift": 2,"limit_choose": 2, "counter": 6}`
	req, err := http.NewRequest("POST", "http://localhost/newRota", bytes.NewBufferString(oParams2))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cookie", cookie)

	// before we actually execute our api function, we need to expect required DB actions
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO rota").
		WithArgs(291255583271555074, "一月份值班表", "oj134ltvn555544444_4abcdefgh", 2, 2, 6).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// now we execute our request
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
	assert.JSONEq(t, resp.Body.String(), `{"status": 0, "rota_id": "291255583271555074"}`)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

type testRotaid int64

func (t testRotaid) GetRotaId() int64 {
	return int64(t)
}
