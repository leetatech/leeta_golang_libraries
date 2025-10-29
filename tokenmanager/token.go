package jwtmiddleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/leetatech/leeta_golang_libraries/errs"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/metadata"
)

// UserClaims represents the JWT claims for a user, including standard claims and custom fields.
type UserClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Phone  string `json:"phone"`
}

// Manager handles JWT operations using RSA public and private keys.
type Manager struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// TokenManager defines the interface for JWT token parsing and user claims extraction.
type TokenManager interface {
	ParseToken(signedTokenString string) (*UserClaims, error)
	putClaimsOnContext(ctx context.Context, claims *UserClaims) (context.Context, error)
	ExtractUserClaims(ctx context.Context) (*UserClaims, error)
}

var _ TokenManager = &Manager{}

var AuthenticatedUserMetadataKey = "AuthenticatedUser"

// New creates a new Manager instance by parsing the provided RSA public and private keys.
func New(publicKey, privateKey string) (*Manager, error) {
	return generateKey(publicKey, privateKey)
}

// generateKey parses and validates the provided RSA public and private keys, returning a Manager.
func generateKey(publicKey, privateKey string) (*Manager, error) {
	pubPEM, err := ensurePEMFormat(publicKey, "PUBLIC KEY")
	if err != nil {
		return nil, fmt.Errorf("failed to process public key: %w", err)
	}

	privPEM, err := ensurePEMFormat(privateKey, "PRIVATE KEY")
	if err != nil {
		return nil, fmt.Errorf("failed to process private key: %w", err)
	}

	tokenGeneratorPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	tokenGeneratorPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privPEM))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &Manager{
		publicKey:  tokenGeneratorPublicKey,
		privateKey: tokenGeneratorPrivateKey,
	}, nil
}

// ensurePEMFormat ensures the provided key string is in PEM format, adding headers if necessary.
func ensurePEMFormat(key string, blockType string) (string, error) {
	key = strings.TrimSpace(key)

	if strings.HasPrefix(key, "-----BEGIN") {
		// Already PEM, return as is
		return key, nil
	}

	// Clean possible whitespace or line breaks
	cleaned := regexp.MustCompile(`\s+`).ReplaceAllString(key, "")
	decoded, err := base64.StdEncoding.DecodeString(cleaned)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	return buildPEM(blockType, decoded), nil
}

// buildPEM constructs a PEM-formatted string from DER bytes and the given block type.
func buildPEM(blockType string, derBytes []byte) string {
	b64 := base64.StdEncoding.EncodeToString(derBytes)

	var pemLines []string
	for i := 0; i < len(b64); i += 64 {
		end := i + 64
		if end > len(b64) {
			end = len(b64)
		}
		pemLines = append(pemLines, b64[i:end])
	}

	return fmt.Sprintf("-----BEGIN %s-----\n%s\n-----END %s-----\n",
		blockType,
		strings.Join(pemLines, "\n"),
		blockType,
	)
}

// GenerateTokenWithExpiration generates a signed JWT token with a 24-hour expiration using the provided claims.
func (handler *Manager) GenerateTokenWithExpiration(claims *UserClaims, expiresAt time.Time) (string, error) {
	claims.ExpiresAt = jwt.NewNumericDate(expiresAt)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(handler.privateKey)
}

// GenerateAuthenticationToken sets user details and generates a signed JWT token with expiration.
func (handler *Manager) GenerateAuthenticationToken(phone, userID string, expiresAt time.Time) (string, error) {
	claims := UserClaims{
		Phone:  phone,
		UserID: userID,
	}
	return handler.GenerateTokenWithExpiration(&claims, expiresAt)
}

// Valid checks if the token's claims are still valid (not expired).
func (claims *UserClaims) Valid() error {
	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		return fmt.Errorf("invalid expiration time: %w", err)
	}
	if expirationTime == nil {
		return errors.New("expiration time is not set")
	}

	if time.Now().After(expirationTime.Time) {
		return errors.New("token has expired")
	}
	return nil
}

// ParseToken parses a signed JWT string and returns the user claims if valid.
func (handler *Manager) ParseToken(signedTokenString string) (*UserClaims, error) {
	t, err := jwt.ParseWithClaims(signedTokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, errors.New("invalid signing algorithm")
		}
		return handler.publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	return t.Claims.(*UserClaims), nil
}

// PutClaimsOnContext put user token on context
func (handler *Manager) putClaimsOnContext(ctx context.Context, claims *UserClaims) (context.Context, error) {
	jsonClaims, err := json.Marshal(claims)
	if err != nil {
		return nil, err
	}

	return metadata.AppendToOutgoingContext(ctx, AuthenticatedUserMetadataKey, string(jsonClaims)), nil
}

// ValidateMiddleware middleware required endpoints: verify claims and put claims on context
func (handler *Manager) ValidateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			handler.validateHeaderToken(authorizationHeader, next, w, r, false)
		} else {
			errMsg := errors.New("no token in authorization header")
			WriteJSONResponse(w, http.StatusUnauthorized, errs.Body(errs.ErrorUnauthorized, errMsg))
			return
		}
	})
}

// ValidateRestrictedAccessMiddleware middleware required endpoints: verify claims
// extensively check if they have superior access to these endpoints
// and put claims on context
func (handler *Manager) ValidateRestrictedAccessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader != "" {
			handler.validateHeaderToken(authorizationHeader, next, w, r, true)
		} else {
			WriteJSONResponse(w, http.StatusUnauthorized, errs.Body(errs.ErrorUnauthorized, errors.New("no token in authorization header")))
			return
		}

	})
}

// validateHeaderToken validates the authorization header, parses the token, and injects claims into the request context.
func (handler *Manager) validateHeaderToken(authorizationHeader string, next http.Handler, w http.ResponseWriter, r *http.Request, isAdminPrivileged bool) {
	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) != 2 {
		// malformed token or missing "Bearer" part
		WriteJSONResponse(w, http.StatusUnauthorized, errs.Body(errs.ErrorUnauthorized, errors.New("malformed token in authorization header")))
		return // return early to prevent further execution
	}

	if bearerToken[1] == "" {
		log.Error().Msg("bearer token is empty")
		WriteJSONResponse(w, http.StatusUnauthorized, errs.Body(errs.ErrorUnauthorized, errors.New("empty token in authorization header")))
		return
	}

	claims, err := handler.ParseToken(bearerToken[1])
	if err != nil {
		log.Error().Msgf("unable to parse token string: %v", err)
		WriteJSONResponse(w, http.StatusUnauthorized, errs.Body(errs.ErrorUnauthorized, err))
		return
	}

	ctx, err := handler.putClaimsOnContext(r.Context(), claims)
	if err != nil {
		log.Error().Msgf("unable to put claims on context: %v", err)
		WriteJSONResponse(w, http.StatusUnauthorized, errs.Body(errs.ErrorUnauthorized, err))
		return
	}
	next.ServeHTTP(w, r.WithContext(ctx))
}

// ExtractUserClaims returns claims from an authenticated user
func (handler *Manager) ExtractUserClaims(ctx context.Context) (*UserClaims, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return nil, errors.New("unable to parse authenticated user")
	}
	jsonClaims := md.Get(AuthenticatedUserMetadataKey)
	if len(jsonClaims) == 0 {
		return nil, errors.New("fail decode authenticated user claims")
	}

	var claims UserClaims
	err := json.Unmarshal([]byte(jsonClaims[0]), &claims)
	if err != nil {
		return nil, errors.New("fail to unmarshall claims")
	}

	return &claims, nil
}

// WriteJSONResponse writes a JSON response with the given status code and response data to the HTTP response writer.
func WriteJSONResponse(w http.ResponseWriter, code int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data := struct {
		Data any `json:"data"`
	}{
		Data: response,
	}

	if response != nil {
		err := json.NewEncoder(w).Encode(&data)
		if err != nil {
			log.Err(err).Msg("fail to encode response")
			return
		}
	}

}
