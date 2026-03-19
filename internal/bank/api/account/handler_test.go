package account_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/romangurevitch/go-training/internal/bank/api/account"
	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/service/mocks"
)

const testSecret = "test-secret"

// testToken issues a signed JWT for use in test Authorization headers.
func testToken(t *testing.T, sub, scope string) string {
	t.Helper()
	claims := middleware.Claims{
		Scope: scope,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testSecret))
	require.NoError(t, err)
	return signed
}

// setupRouter builds a minimal Gin engine with the account routes wired up.
func setupRouter(svc *mocks.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := account.New(svc)

	v1 := r.Group("/v1/accounts")
	v1.Use(middleware.JWTMiddleware(testSecret))
	{
		v1.GET("/:id", middleware.RequireScope("accounts:read"), h.GetAccount)
		v1.POST("", middleware.RequireScope("accounts:write"), h.CreateAccount)
	}
	return r
}

var testAccount = &domain.Account{
	ID:      "ACC-001",
	Owner:   "alice",
	Balance: 10000, // 100.00
	Status:  domain.StatusOpen,
}

func TestGetAccount(t *testing.T) {
	type fields struct {
		svc func(t *testing.T) *mocks.Service
	}
	tests := []struct {
		name     string
		fields   fields
		id       string
		scope    string
		wantCode int
		wantBody any
	}{
		{
			name: "success — returns 200 with account",
			fields: fields{
				svc: func(t *testing.T) *mocks.Service {
					m := mocks.NewService(t)
					m.EXPECT().GetAccount(mock.Anything, "ACC-001").Return(testAccount, nil).Once()
					return m
				},
			},
			id:       "ACC-001",
			scope:    "accounts:read",
			wantCode: http.StatusOK,
			wantBody: map[string]any{"id": "ACC-001", "owner": "alice"},
		},
		{
			name: "not found — returns 404",
			fields: fields{
				svc: func(t *testing.T) *mocks.Service {
					m := mocks.NewService(t)
					m.EXPECT().GetAccount(mock.Anything, "MISSING").Return(nil, domain.ErrAccountNotFound).Once()
					return m
				},
			},
			id:       "MISSING",
			scope:    "accounts:read",
			wantCode: http.StatusNotFound,
		},
		{
			name: "wrong scope — returns 403",
			fields: fields{
				svc: func(t *testing.T) *mocks.Service {
					return mocks.NewService(t) // no calls expected
				},
			},
			id:       "ACC-001",
			scope:    "accounts:write", // missing accounts:read
			wantCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.fields.svc(t))
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/v1/accounts/"+tt.id, nil)
			req.Header.Set("Authorization", "Bearer "+testToken(t, "alice", tt.scope))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			if tt.wantBody != nil {
				var got map[string]any
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
				for k, v := range tt.wantBody.(map[string]any) {
					assert.Equal(t, v, got[k])
				}
			}
		})
	}
}

func TestCreateAccount(t *testing.T) {
	type fields struct {
		svc func(t *testing.T) *mocks.Service
	}
	tests := []struct {
		name     string
		fields   fields
		body     any
		scope    string
		wantCode int
	}{
		{
			name: "success — returns 201",
			fields: fields{
				svc: func(t *testing.T) *mocks.Service {
					m := mocks.NewService(t)
					m.EXPECT().CreateAccount(mock.Anything, "alice").Return(testAccount, nil).Once()
					return m
				},
			},
			body:     map[string]string{"owner": "alice"},
			scope:    "accounts:write",
			wantCode: http.StatusCreated,
		},
		{
			name: "missing owner — returns 400",
			fields: fields{
				svc: func(t *testing.T) *mocks.Service {
					return mocks.NewService(t)
				},
			},
			body:     map[string]string{},
			scope:    "accounts:write",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.fields.svc(t))
			bodyBytes, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/accounts", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken(t, "alice", tt.scope))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
