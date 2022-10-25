// Package httpsrv has a http server for dnscontroller
package httpsrv

import (
	"net/http"
	"os"
	"text/template"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	gindump "github.com/tpkeeper/gin-dump"
	"github.com/volatiletech/sqlboiler/v4/boil"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.hollow.sh/toolbox/ginjwt"
	"go.hollow.sh/toolbox/version"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zapio"

	v1router "go.hollow.sh/dnscontroller/pkg/api/v1/router"
)

// Server contains the HTTP server configuration
type DB struct {
	Driver *sqlx.DB
	Debug  bool
}
type Server struct {
	Logger         *zap.SugaredLogger
	Listen         string
	Debug          bool
	DB             DB
	AuthConfig     ginjwt.AuthConfig
	TrustedProxies []string
	TemplateFields map[string]template.Template
}

var (
	readTimeout  = 10 * time.Second
	writeTimeout = 20 * time.Second
)

func (s *Server) setup() *gin.Engine {
	var (
		authMW *ginjwt.Middleware
		err    error
	)

	authMW, err = ginjwt.NewAuthMiddleware(s.AuthConfig)
	if err != nil {
		s.Logger.Fatal("failed to initialize auth middleware", "error", err)
	}

	// Setup default gin router
	r := gin.New()

	// Set the trusted proxies, if they were specified by config
	if len(s.TrustedProxies) > 0 {
		err = r.SetTrustedProxies(s.TrustedProxies)
		if err != nil {
			s.Logger.Fatal("failed to set gin trusted proxies", "error", err)
		}
	}

	p := ginprometheus.NewPrometheus("gin")

	// Remove any params from the URL string to keep the number of labels down
	p.ReqCntURLLabelMappingFn = func(c *gin.Context) string {
		return c.FullPath()
	}

	p.Use(r)

	logF := func(c *gin.Context) []zapcore.Field {
		return []zapcore.Field{
			zap.String("jwt_subject", ginjwt.GetSubject(c)),
			zap.String("jwt_user", ginjwt.GetUser(c)),
		}
	}
	loggerWithContext := s.Logger.With(zap.String("component", "httpsrv"))
	r.Use(ginzap.GinzapWithConfig(loggerWithContext.Desugar(), &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		TraceID:    true,
		Context:    logF,
	}))
	r.Use(ginzap.RecoveryWithZap(loggerWithContext.Desugar(), true))

	tp := otel.GetTracerProvider()
	if tp != nil {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}

		r.Use(otelgin.Middleware(hostname, otelgin.WithTracerProvider(tp)))
	}

	// Version endpoint returns build information
	r.GET("/version", s.version)

	// Health endpoints
	r.GET("/healthz", s.livenessCheck)
	r.GET("/healthz/liveness", s.livenessCheck)
	r.GET("/healthz/readiness", s.readinessCheck)

	v1Rtr := v1router.New(authMW, s.DB.Driver, s.Logger)

	// Host our latest version of the API under / in addition to /api/v*
	latest := r.Group("/")
	{
		v1Rtr.Routes(latest)
		if s.Debug {
			r.Use(gindump.Dump())
		}
	}

	v1 := r.Group(v1router.V1URI)
	{
		v1Rtr.Routes(v1)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "invalid request - route not found"})
	})

	return r
}

// NewServer returns a configured server
func (s *Server) NewServer() *http.Server {
	if !s.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	return &http.Server{
		Handler:      s.setup(),
		Addr:         s.Listen,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

// Run will start the server listening on the specified address
func (s *Server) Run() error {
	gin.SetMode(gin.ReleaseMode)

	boil.DebugMode = s.DB.Debug

	debugWriter := &zapio.Writer{Log: s.Logger.Desugar()}
	defer debugWriter.Close()

	boil.DebugWriter = debugWriter

	return s.setup().Run(s.Listen)
}

// livenessCheck ensures that the server is up and responding
func (s *Server) livenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

// readinessCheck ensures that the server is up and that we are able to process
// requests. Currently our only dependency is the DB so we just ensure that it
// is responding.
func (s *Server) readinessCheck(c *gin.Context) {
	if err := s.DB.Driver.PingContext(c.Request.Context()); err != nil {
		s.Logger.Errorw("readiness check db ping failed", "err", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "UP",
	})
}

// version returns the metadataservice build information
func (s *Server) version(c *gin.Context) {
	c.JSON(http.StatusOK, version.String())
}
