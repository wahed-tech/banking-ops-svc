package auth

import (
	"context"
	"net/http"
	"strings"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/wahed-tech/banking-ops-svc/conf"
	"github.com/wahed-tech/go-lib/v2/errors"
	"github.com/wahed-tech/go-lib/v2/jwt"
)

const supabaseSigningMethod = "ES256"

const samlMethod = "sso/saml"

var (
	errMissingAuthConf = errors.New("auth config incomplete")
	errNotSSO          = errors.New("token not issued by our Google SAML login")
	errEmailDomain     = errors.New("email not in allowed domain")
	errMissingSubject  = errors.New("token has no subject")
)

// Caller - verified portal user behind a request
type Caller struct {
	UserID string
	Email  string
}

type callerKey struct{}

func CallerFrom(ctx context.Context) (Caller, bool) {
	c, ok := ctx.Value(callerKey{}).(Caller)
	return c, ok
}

type Authenticator struct {
	validator     jwt.Validator
	ssoProviderID string
	emailSuffix   string
}

func NewAuthenticator(keyfunc jwt.KeyFunc, c conf.AuthConf) (*Authenticator, error) {
	if c.Issuer == "" || c.Audience == "" || c.SSOProviderID == "" || c.AllowedDomain == "" {
		return nil, errMissingAuthConf.Here()
	}

	validator := jwt.NewValidator(keyfunc, nil, func(v *jwt.Validator) {
		v.ParserOptions = []jwtv5.ParserOption{
			jwtv5.WithValidMethods([]string{supabaseSigningMethod}),
			jwtv5.WithIssuer(c.Issuer),
			jwtv5.WithAudience(c.Audience),
			jwtv5.WithExpirationRequired(),
		}
	})

	return &Authenticator{
		validator:     validator,
		ssoProviderID: c.SSOProviderID,
		emailSuffix:   "@" + strings.ToLower(c.AllowedDomain),
	}, nil
}

func (a *Authenticator) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		caller, err := a.verifyCaller(r)
		if err != nil {
			errors.Log(errors.Wrap(err, errors.LevelWarning).Prepend("unable to authenticate portal token"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), callerKey{}, caller)))
	})
}

func (a *Authenticator) verifyCaller(r *http.Request) (Caller, error) {
	claims, err := a.validator.Claims(r)
	if err != nil {
		return Caller{}, err
	}

	if !a.loggedInWithSAML(claims) {
		return Caller{}, errNotSSO.Here()
	}

	email, _ := claims["email"].(string)
	email = strings.ToLower(email)
	if !strings.HasSuffix(email, a.emailSuffix) {
		return Caller{}, errEmailDomain.Here()
	}

	sub, _ := claims.GetSubject()
	if sub == "" {
		return Caller{}, errMissingSubject.Here()
	}

	return Caller{UserID: sub, Email: email}, nil
}

func (a *Authenticator) loggedInWithSAML(claims jwtv5.MapClaims) bool {
	entries, _ := claims["amr"].([]any)
	for _, e := range entries {
		entry, _ := e.(map[string]any)
		if entry["method"] == samlMethod && entry["provider"] == a.ssoProviderID {
			return true
		}
	}

	return false
}
