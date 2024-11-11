package simulator

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"ride_sharing_api/app/assert"
	"ride_sharing_api/app/common"
	"strings"

	"golang.org/x/oauth2"
)

type simulatorWrapper struct {
	base                            Simulator
	overrideLogOutput               *io.Writer
	overrideHttpGet                 *func(url string) (*http.Response, error)
	overrideHttpRedirect            *func(w http.ResponseWriter, r *http.Request, url string, code int)
	overrideOauthGoogleCodeExchange *func(ctx context.Context, cfg *oauth2.Config, code string) (*oauth2.Token, error)
	overrideDbName                  *string
}

// builder functions

func FromBase(base Simulator) *simulatorWrapper {
	return &simulatorWrapper{base: base}
}

func (b *simulatorWrapper) LogTo(w io.Writer) {
	b.overrideLogOutput = &w
}

func (b *simulatorWrapper) WithDb(name string) {
	b.overrideDbName = &name
}

func (b *simulatorWrapper) WithHttpGet(get func(url string) (*http.Response, error)) {
	b.overrideHttpGet = &get
}

func (b *simulatorWrapper) WithHttpRedirect(redirect func(w http.ResponseWriter, r *http.Request, url string, code int)) {
	b.overrideHttpRedirect = &redirect
}

func (b *simulatorWrapper) WithOauthGoogleCodeExchange(exchange func(ctx context.Context, cfg *oauth2.Config, code string) (*oauth2.Token, error)) {
	b.overrideOauthGoogleCodeExchange = &exchange
}

// Pre-create config

func (b *simulatorWrapper) AlwaysPassGoogleOauth() {
	b.WithHttpRedirect(func(w http.ResponseWriter, r *http.Request, url string, code int) {
		if strings.HasPrefix(url, "https://accounts.google.com/o/oauth2/auth") {
			state := strings.Split(url, "state=")[1]
			state = strings.Split(state, "&")[0]

			callbackUrl := b.GetEnvRequired(common.ENV_HOST_ADDR) + "/auth/google/callback?state=" + state
			b.base.HttpRedirect(w, r, callbackUrl, code)
			return
		}

		b.base.HttpRedirect(w, r, url, code)
	})

	b.WithHttpGet(func(url string) (*http.Response, error) {
		if strings.HasPrefix(url, "https://www.googleapis.com/oauth2/v2/userinfo?access_token=") {
			id := "fake-id"
			email := "fake@gmail.com"
			name := "Ronald McDonald"
			verifiedEmail := true

			profile := common.GoogleProfile{
				Id:            &id,
				Email:         &email,
				Name:          &name,
				VerifiedEmail: &verifiedEmail,
			}

			data, err := json.Marshal(profile)
			assert.True(err == nil)

			res := http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(data)),
			}

			return &res, nil
		}

		return b.base.HttpGet(url)
	})

	b.WithOauthGoogleCodeExchange(func(ctx context.Context, cfg *oauth2.Config, code string) (*oauth2.Token, error) {
		return &oauth2.Token{
			AccessToken:  "fake-access-token",
			RefreshToken: "fake-refresh-token",
		}, nil
	})
}

// `Simulator` implementation

func (b *simulatorWrapper) FsStat(name string) (fs.FileInfo, error) {
	return b.base.FsStat(name)
}

func (b *simulatorWrapper) FsCreate(name string) (File, error) {
	return b.base.FsCreate(name)
}

func (b *simulatorWrapper) HttpNewServerMux() HTTPMux {
	return &HTTPMuxRealWorld{inner: http.NewServeMux()}
}

func (b *simulatorWrapper) HttpListenAndServe(handler http.Handler, addr string) error {
	return b.base.HttpListenAndServe(handler, addr)
}

func (b *simulatorWrapper) HttpGet(url string) (resp *http.Response, err error) {
	if b.overrideHttpGet != nil {
		return (*b.overrideHttpGet)(url)
	}
	return b.base.HttpGet(url)
}

func (b *simulatorWrapper) HttpRedirect(w http.ResponseWriter, r *http.Request, url string, code int) {
	if b.overrideHttpRedirect != nil {
		(*b.overrideHttpRedirect)(w, r, url, code)
		return
	}
	b.base.HttpRedirect(w, r, url, code)
}

func (b *simulatorWrapper) OauthGoogleExchangeCode(ctx context.Context, cfg *oauth2.Config, code string) (*oauth2.Token, error) {
	if b.overrideOauthGoogleCodeExchange != nil {
		return (*b.overrideOauthGoogleCodeExchange)(ctx, cfg, code)
	}
	return b.base.OauthGoogleExchangeCode(ctx, cfg, code)
}

func (b *simulatorWrapper) LogOutput() io.Writer {
	if b.overrideLogOutput != nil {
		return *b.overrideLogOutput
	}
	return b.base.LogOutput()
}

func (b *simulatorWrapper) GetEnv(key string) string {
	return b.base.GetEnv(key)
}

func (b *simulatorWrapper) GetEnvRequired(key string) string {
	return b.base.GetEnvRequired(key)
}

func (b *simulatorWrapper) SqlOpen(driverName string, dataSourceName string) (DB, error) {
	return b.base.SqlOpen(driverName, dataSourceName)
}

func (b *simulatorWrapper) DbName() string {
	if b.overrideDbName != nil {
		return *b.overrideDbName
	}
	return b.base.DbName()
}

func (b *simulatorWrapper) RandCrypto(bytes []byte) {
	b.base.RandCrypto(bytes)
}
