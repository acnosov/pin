package auth

import (
	"context"
	"errors"
	"fmt"
	api "github.com/aibotsoft/gen/epinapi"
	"github.com/aibotsoft/pin/pkg/config"
	"github.com/aibotsoft/pin/pkg/store"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var UnauthorizedError = errors.New("401 Unauthorized")

const checkLoginPeriod = time.Minute * 5

type Auth struct {
	cfg     *config.Config
	log     *zap.SugaredLogger
	store   *store.Store
	eClient *api.APIClient
	//conf    *config_client.ConfClient
	Account store.Account
	token   store.Token
}

func New(cfg *config.Config, log *zap.SugaredLogger, store *store.Store) *Auth {
	account, err := store.GetAccount(context.Background())
	if err != nil {
		log.Panic(err)
	}
	clientConfig := api.NewConfiguration()
	clientConfig.Host = cfg.PinEsportHost
	tr := &http.Transport{TLSHandshakeTimeout: 0 * time.Second, IdleConnTimeout: 0 * time.Second}
	clientConfig.HTTPClient = &http.Client{Transport: tr}
	clientConfig.Debug = cfg.PinDebug
	client := api.NewAPIClient(clientConfig)

	a := &Auth{cfg: cfg, log: log, store: store, Account: account, eClient: client}
	//err = a.Login(context.Background())
	//if err != nil {
	//	log.Error(err)
	//}
	//go a.AuthJob()
	a.token, err = a.store.LoadToken(context.Background())
	if err != nil {
		log.Info(err)
	}
	return a
}
func (a *Auth) AuthJob() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		err := a.AuthRound(ctx)
		if err != nil {
			a.log.Info(err)
		}
		cancel()
		time.Sleep(checkLoginPeriod)
	}
}

func (a *Auth) AuthRound(ctx context.Context) error {
	err := a.CheckLogin(ctx)
	switch err {
	case nil:
		return nil
	case UnauthorizedError:
		a.log.Info("session_expired_begin_login")
		err := a.Login(ctx)
		if err != nil {
			return err
		}
		err = a.store.SaveToken(ctx, a.token)
		if err != nil {
			return err
		}
		return nil
	default:
		return err
	}
}
func (a *Auth) Login(ctx context.Context) error {
	resp, _, err := a.eClient.ClientApi.Login(a.AuthLogin(ctx)).LoginRequest(api.LoginRequest{
		Username: a.Account.Username, Password: a.Account.Password, TrustCode: a.token.TrustCode}).Execute()
	if err != nil {
		apiErr := err.(api.GenericOpenAPIError)
		a.log.Infow("err", "err", string(apiErr.Body()))
		return err
	}
	a.token.Session = resp.GetToken()
	a.token.TrustCode = resp.GetTrustCode()
	a.log.Infow("resp", "", resp)
	return nil
}
func (a *Auth) CheckLogin(ctx context.Context) error {
	a.log.Infow("begin_check_login")
	_, r, err := a.eClient.ClientApi.CheckLogin(a.Auth(ctx), a.token.Session).Execute()
	if err != nil {
		if r != nil && r.StatusCode == 401 {
			return UnauthorizedError
		}
		return fmt.Errorf("check_login_error: %s", err)
	}
	return nil
}

func (a *Auth) Auth(ctx context.Context) context.Context {
	keyMap := map[string]api.APIKey{
		"x-api-key":     {Key: a.cfg.PinApiKey},
		"x-session":     {Key: a.token.Session},
		"x-device-uuid": {Key: a.cfg.PinApiDevice},
	}
	return context.WithValue(ctx, api.ContextAPIKeys, keyMap)
}
func (a *Auth) AuthLogin(ctx context.Context) context.Context {
	keyMap := map[string]api.APIKey{
		"x-api-key": {Key: a.cfg.PinApiKey},
		//"x-session":     {Key: a.token.Session},
		"x-device-uuid": {Key: a.cfg.PinApiDevice},
	}
	return context.WithValue(ctx, api.ContextAPIKeys, keyMap)
}
