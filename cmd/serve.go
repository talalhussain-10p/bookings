package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
	"github.com/pkg/errors"

	"github.com/bookings/pkg/param"
	"github.com/bookings/pkg/router"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	listenAddr string
	corsHost   []string
	appEnvProd bool
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run http server",
	Long:  `this command serves the backend server.`,
	RunE:  run,
}

func init() {
	serveCmd.Flags().StringSliceVar(&corsHost, "cors-host", []string{"*"}, "cors allowed host, comma separated")
	serveCmd.Flags().StringVarP(&listenAddr, "listen-addr", "a", ":9111", "listen address")
	serveCmd.Flags().BoolVarP(&appEnvProd, "app-env", "e", false, "production environment enabler")

	rootCmd.AddCommand(serveCmd)

	viper.BindPFlags(serveCmd.Flags())
}

func run(cmd *cobra.Command, args []string) error {
	logrus.NewEntry(logrus.StandardLogger())
	p := &param.Param{
		AppENV: appEnvProd,
	}
	return serve(p)
}

func serve(p *param.Param) error {
	chErr := make(chan error)

	session := sessionManager()

	p.Session = session

	go func() {
		logrus.Info(fmt.Sprintf("listening on %s", listenAddr))
		chErr <- serveHTTP(p)
	}()
	return <-chErr
}

func sessionManager() *scs.SessionManager {
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	if appEnvProd {
		session.Cookie.Secure = appEnvProd
	}
	return session
}

func serveHTTP(p *param.Param) error {
	gin.SetMode(gin.ReleaseMode)
	r := chi.NewRouter()

	p.CSRFHandler = nosurf.New(r)
	cors.New(cors.Options{
		AllowedOrigins: corsHost,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
	})
	param.Inject(r, p)
	router.HandlerHTTP(r)

	return errors.Wrap(http.ListenAndServe(listenAddr, r), "unable to run server")
}
