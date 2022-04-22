package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	echolog "github.com/labstack/gommon/log"
	"github.com/rs/zerolog/log"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/authorization"
	"github.com/tradeface/suggest_service/pkg/elastic"
	"github.com/tradeface/suggest_service/pkg/middleware"
	"github.com/tradeface/suggest_service/pkg/mongo"
	"github.com/tradeface/suggest_service/pkg/server"
	"github.com/tradeface/suggest_service/pkg/store"
)

//TODO: config cli/dockersecrets

//  https://pkg.go.dev/github.com/golang-jwt/jwt/v4

const (
	// APPNAME contains the name of the program
	APPNAME = "suggest_service"
	// APPVERSION contains the version of the program
	APPVERSION = "0.0.2"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

//tmp to get token
func login(c echo.Context) error {

	// Set custom claims
	claims := &authorization.AuthClaims{
		"Jon Snow",
		"123456",
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func main() {

	cfg, err := conf.NewDefaultConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	dbconn, err := mongo.NewClient(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	esconn, err := elastic.NewElastic(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to es")
	}

	stores, err := store.New(dbconn, esconn, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	auth := &authorization.AuthChecker{}

	srv, err := server.NewServer(cfg, stores, auth)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bind api")
	}

	e := echo.New()
	// shut up
	e.Logger.SetOutput(ioutil.Discard)
	e.Logger.SetLevel(echolog.OFF)

	// tmp Login route
	e.GET("/login", login)

	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	e.Use(middleware.JWTWithConfig(&middleware.JWTConfig{}, auth))

	srv.RegisterHandlers(e)

	log.Info().Str("addr", cfg.Addr).Msg("starting http listener")
	err = e.Start(cfg.Addr)
	log.Fatal().Err(err).Msg("Server failed")
}
