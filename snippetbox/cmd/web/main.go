package main

import (
	"context"
	"crypto/tls"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sotnikea/Go_Learn/tree/main/snippetbox/internal/models"

	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Define an application struct to hold the application-wide dependencies
type application struct {
	logger         *slog.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Define a new command-line flags
	addr := flag.String("addr", "0.0.0.0:4000", "HTTP network address")
	uri := flag.String("uri", "mongodb://web:111@localhost:27017/snippetbox", "Mongo database uri")
	flag.Parse()

	// Initialize a new structured logger, which writes to the standard out stream
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Open database
	database, err := openDB(*uri, "snippetbox")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Get client connected with database
	client := database.Client()

	// Connection pool must closed before the main() function exits
	defer client.Disconnect(context.TODO())

	// Initialize a new template cache...
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a decoder instance...
	formDecoder := form.NewDecoder()

	// Initialize a new session manager. Configure it to use Mongo database as the session store,
	// and set a lifetime of 12 hours
	sessionManager := scs.New()
	sessionManager.Store = mongodbstore.New(database)
	sessionManager.Lifetime = 12 * time.Hour

	// Setting Secure attribute means that the cookie will only be sent by a user's web
	// browser when a HTTPS connection is being used
	sessionManager.Cookie.Secure = true

	// Initialize a new instance of our application struct, containing the dependencies
	app := &application{
		logger:         logger,
		snippets:       &models.SnippetModel{DB: database},
		users:          &models.UserModel{DB: database},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a tls.Config struct to hold curve preferences value, so that only elliptic curves with
	// assembly implementations are used
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Initialize a new http.Server struct
	srv := &http.Server{
		Addr:      *addr,
		Handler:   app.routes(),
		ErrorLog:  slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig: tlsConfig,
		// Add Idle, Read and Write timeouts to the server.
		// IdleTimeout:  time.Minute,
		// ReadTimeout:  5 * time.Second,
		// WriteTimeout: 10 * time.Second,
	}

	// Log the starting server message at Info severity
	logger.Info("starting server", "addr", srv.Addr)

	// Start the HTTPS server. Pass in the paths to the TLS certificate and corresponding private key
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(uri string, dbName string) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return db, nil
}
