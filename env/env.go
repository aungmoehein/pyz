package env

import (
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"math/rand"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

// Specification contains all environment variable dependencies
type Specification struct {
	// system configs
	Port                   string        `split_words:"true" default:":8080"`
	AppName                string        `envconfig:"APPNAME" default:"maputin"`
	DefaultLanguage        string        `split_words:"true" default:"en"`
	LogFilePath            string        `split_words:"true" default:""`
	DatabaseURL            string        `split_words:"true" required:"true" `
	DatabaseEngine         string        `split_words:"true" required:"true"`
	DatabaseConnLifetime   time.Duration `split_words:"true" default:"3600s"`
	Timezone               string        `split_words:"true" default:"Asia/Rangoon"`
	AuthorizeTokenLifetime time.Duration `split_words:"true" default:"180s"`
	AccessTokenLifetime    time.Duration `split_words:"true" default:"180s"`
}

// Article contains News variable in article.toml
type Article struct {
	News map[string]New
}

// New contains new vairables in article.toml
type New struct {
	PostedUser    string
	PostedDate    string
	Title         string
	Body          string
	PostImage     string
	PostedUserImg string
}

// Packs contains Packs variable in packs.toml
type Packs struct {
	Packs map[string]Pack
}

// Pack describe pack variable in packs.toml
type Pack struct {
	PackID           string
	PackTitle        string
	PackDescription  string
	PackPartner      string
	PackAvailability bool
	PackImgURL       string
	ContactNo        string
	CreatedAt        string
}

// FAQs contains FAQs variable in FAQs.toml
type FAQs struct {
	FAQs map[string]FAQ
}

// FAQ contains FAQ variable in FAQs.toml
type FAQ struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// Envs hold all enviroment variables needed for the app
type Envs struct {
	*Specification
	*Article
	*Packs
	*FAQs

	ID          string
	Location    *time.Location
	articlePath string
	packsPath   string
	faqsPath    string
}

var envInstance *Envs
var envOnce sync.Once
var environ = GetEnvironment()

var logInstance *log.Logger
var logOnce sync.Once
var logger = GetLogger()

const (
	settings     = `PYZ_SETTINGS`
	dbConnParams = `` +
		`?parseTime=true&` +
		`loc=Asia%2FRangoon&` +
		`sql_mode=TRADITIONAL&` +
		`time_zone=%27Asia%2FRangoon%27`
)

// GetLogger return a new singleton GetLogger object
func GetLogger() *log.Logger {
	logOnce.Do(func() {
		logInstance = log.New()
		logInstance.Level = log.InfoLevel
		logInstance.Formatter = &log.TextFormatter{
			DisableSorting:  true,
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		}

		var logFile *os.File
		var err error

		var filePath = environ.LogFilePath
		var fileFlags = os.O_CREATE | os.O_APPEND | os.O_WRONLY

		if logFile, err = os.OpenFile(filePath, fileFlags, 0644); err != nil {
			logInstance.Warn("Unable to open log file - ", filePath)
			logInstance.Out = os.Stdout
		} else {
			logInstance.Out = logFile
		}
	})

	return logInstance
}

// GetEnvironment return a new singleton Env object
func GetEnvironment() *Envs {
	envOnce.Do(func() {
		var ok bool
		var err error
		var settingPath string

		if settingPath, ok = os.LookupEnv(settings); !ok {
			log.Fatal(settings, ` enviroment variable is not set yet.`)
		}

		config := path.Join(settingPath, "system.env")
		if err = godotenv.Load(config); err != nil {
			log.Fatal(err)
		}

		var spec Specification
		if err = envconfig.Process("", &spec); err != nil {
			log.Fatal(err)
		}

		var article Article
		var articlePath = path.Join(settingPath, "article.toml")
		if _, err = toml.DecodeFile(articlePath, &article); err != nil {
			log.Fatal(err)
		}

		var packs Packs
		var packsPath = path.Join(settingPath, "packs.toml")
		if _, err = toml.DecodeFile(packsPath, &packs); err != nil {
			log.Fatal(err)
		}

		var faqs FAQs
		var faqsPath = path.Join(settingPath, "faq.toml")
		if _, err = toml.DecodeFile(faqsPath, &faqs); err != nil {
			log.Fatal(err)
		}

		var location *time.Location
		if location, err = time.LoadLocation(spec.Timezone); err != nil {
			log.Error("Unable to get timzeone", err)
		}

		spec.DatabaseURL = spec.DatabaseURL + dbConnParams
		envInstance = &Envs{
			ID:            strconv.Itoa(rand.Int()),
			Location:      location,
			Specification: &spec,
			Article:       &article,
			Packs:         &packs,
			FAQs:          &faqs,

			articlePath: articlePath,
		}
	})

	return envInstance
}
