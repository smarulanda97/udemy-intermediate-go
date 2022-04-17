package kernel

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/smarulanda97/app-stripe/internal/utils"
)

type IKernel interface {
	Start(version string)
	CreateLoggers() (*log.Logger, *log.Logger)
	FileServer() http.Handler
}

type Kernel struct {
	Port            int
	ApiUrl          string
	CurrentVersion  string
	Environment     string
	RoutePrefix     string
	ThemePath       string
	PublicFilesPath string
	Stripe          utils.StripeConfig
	Database        utils.DatabaseConfig
}

func (k *Kernel) Start(version string) {
	flag.IntVar(&k.Port, "port", 4000, "Server port to listen on")
	flag.StringVar(&k.ApiUrl, "api", "http://localhost:4001", "URL to api")
	flag.StringVar(&k.PublicFilesPath, "public_files_path", "./static", "Public files path")
	flag.StringVar(&k.ThemePath, "theme_path", "theme/templates", "Path to templates directory")
	flag.StringVar(&k.Environment, "env", "development", "Application environment {development|production}")
	flag.StringVar(&k.Database.Dsn, "dsn", "root:1997@tcp(localhost:3306)/app_db?parseTime=true&tls=false", "DSN")

	flag.Parse()

	k.Stripe.PublicKey = os.Getenv("STRIPE_KEY")
	k.Stripe.SecretKey = os.Getenv("STRIPE_SECRET")
}

func (k *Kernel) CreateLoggers() (*log.Logger, *log.Logger) {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (k *Kernel) FileServer() http.Handler {
	fileServer := http.FileServer(http.Dir(k.PublicFilesPath))
	return http.StripPrefix("/static", fileServer)
}
