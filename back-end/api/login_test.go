package api

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"

	"schedule/api/internal"
	"schedule/dbb"
	"schedule/wechatid"
)

func TestLogin(t *testing.T) {
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

	login(t, dbInstance, mock, ts, apI)
}

type testOpenid string

func (t testOpenid) GetOpenId(code *wechatid.Code) (string, error) {
	return string(t), nil
}

func login(t *testing.T, dbInstance dbb.DBApi, mock sqlmock.Sqlmock, ts *gin.Engine, apI *internal.Implement) string {
	// init session store
	sessionAuthKey := make([]byte, 32)
	if _, err := rand.Read(sessionAuthKey); err != nil {
		log.Fatalf("%+v", xerrors.Errorf("generate session auth key failed: %w", err))
	}
	store := memstore.NewStore(sessionAuthKey)

	ts.Use(sessions.Sessions("cookie", store))

	// create apI with mocked db, request and response to test
	ts.POST("/login", apI.Login)

	oParams := `{"appid": "wx8f1d4744161ecede", "secret": "c57b47d2748c0b2e27f0fa1e05c50e26", "js_code": "0432q6Ga1oyNpz071BGa1hsCoq32q6Gw", "grant_type": "authorization_code"}`
	req, err := http.NewRequest("POST", "http://localhost/login", bytes.NewBufferString(oParams))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// before we actually execute our api function, we need to expect required DB actions
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT 1 FROM person WHERE openid=?").WithArgs("oj134ltvn555544444_4abcdefgh").
		WillReturnError(sql.ErrNoRows)

	mock.ExpectExec("INSERT INTO person \\(openid\\) VALUES \\(\\?\\)").WithArgs("oj134ltvn555544444_4abcdefgh").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// now we execute our request
	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
	assert.JSONEq(t, resp.Body.String(), `{"status": 0}`)
	cookie := resp.Header().Get("Set-Cookie")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

	return cookie
}
