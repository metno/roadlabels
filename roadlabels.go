package main

import (
	"embed"
	"flag"
	"fmt"

	"html/template"
	"log"
	"net/http"
	"os"

	//"text/template"
	"time"

	frostclient "github.com/metno/frostclient-roadweather"
	"github.com/metno/roadlabels/pkg/db"
	"github.com/metno/roadlabels/pkg/handlers"
	authandlers "github.com/myggen/wwwauth/pkg/handlers"
)

var appHome = ""
var buildTime = ""
var version = ""

var (
	//go:embed templates/**
	templateFiles embed.FS

	//go:embed css/** js/**
	staticFiles embed.FS

	templates map[string]*template.Template
	appRoot   = "roadlabels"
)

func init() {
	if buildTime == "" {
		buildTime = time.Now().UTC().Format("20060102T1504Z")
	}
	if version == "" {
		version = "go run"
	}

}

type StoreHandler interface {
	ObjectExists(string, string) (bool, error)
	string
}

type S3Handler struct {
	Prefix string // S3 Bucket
}

type FSHandler struct {
	Prefix string // image root directory "/lustre/storeB/project/metproduction/products/webcams"
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/roadlabels", http.StatusMovedPermanently)
}

var dbPath *string
var userDBPath *string
var s3accessKey *string
var s3secretKey *string
var s3endpoint string
var h2s3accessKey string
var h2s3secretKey string

var noCacheHeaders = map[string]string{
	"Expires":         "Mon,11 Nov 2019 08:36:00 GMT",
	"Cache-Control":   "no-cache, private, max-age=0, no-store, must-revalidate",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		fn(w, r, r.URL.Path)
	}
}

var frostObses []frostclient.ObsRoadweather

// Frost data examples
var class2FrostObses = make(map[string][]frostclient.ObsRoadweather)

func main() {

	s3endpoint = "rgw.met.no"
	dbPath = flag.String("db-path", "/var/lib/roadlabels/roadcams.db", "path-to-sqlite-db")
	userDBPath = flag.String("userdb-path", "/var/lib/roadlabels/users.db", "path-to-sqlite-userdb")
	s3accessKey = flag.String("access-key", "", "S3 Access Key")
	s3secretKey = flag.String("secret-key", "", "S3 Secret Key")
	flag.Parse()

	if appHome == "" {
		appHome = os.Getenv("PWD")
	}
	log.Printf("AppHome: %s", appHome)
	log.Printf("Version: %s", version)
	log.Printf("BuildTime: %s", buildTime)
	log.Printf("Userdb: %s", *userDBPath)

	if *s3secretKey == "" {
		if os.Getenv("S3SecretKey") != "" {
			log.Printf("arg s3secretKey is empty but found value in environment S3SecretKey")
			h2s3secretKey = os.Getenv("S3SecretKey")
		}
	} else {
		h2s3secretKey = *s3secretKey
	}
	if h2s3secretKey == "" {
		log.Printf("h2s3secretKey is empty")
		flag.Usage()
		return
	}
	if *s3accessKey == "" {
		if os.Getenv("S3AccessKey") != "" {
			log.Printf("arg s3accessKey is empty found value from environment S3AccessKey")
			h2s3accessKey = os.Getenv("S3AccessKey")
		}
	} else {
		h2s3accessKey = *s3accessKey
	}
	if h2s3accessKey == "" {
		log.Printf("h2s3accessKey is empty")
		flag.Usage()
		return
	}

	if *dbPath == "" {
		flag.Usage()
		return
	}

	if *userDBPath == "" {
		flag.Usage()
		return
	}

	db.DBFILE = *dbPath
	var err error

	class2FrostObses, err = frostclient.GetObsMapForLabelApp()
	if err != nil {
		log.Fatalf("%v,", err)
	}

	t := time.Now().UTC()
	var port = 25260
	var portStr = fmt.Sprintf(":%d", port)

	log.Printf("Starting at %s port: %d\n", t.Format("20060102T150405Z"), port)
	log.Printf("Using db %s\n", db.DBFILE)

	// authandlers for authentication
	//auth-authandlers
	authandlers.Files = templateFiles
	authandlers.TemplatesDir = "templates"
	authandlers.AppRoot = appRoot
	err = authandlers.LoadTemplates()
	if err != nil {
		log.Printf("authandlers.LoadTemplates: %v", err)
		os.Exit(1)
	}
	templates = authandlers.Templates
	authandlers.DbFile = *userDBPath
	authandlers.AppRoot = appRoot
	authandlers.SetAuthHandlers()

	handlers.BuildTime = buildTime
	handlers.Version = version
	handlers.AppRoot = appRoot
	handlers.Templates = templates
	handlers.S3endpoint = s3endpoint
	handlers.H2s3accessKey = h2s3accessKey
	handlers.H2s3secretKey = h2s3secretKey
	handlers.FrostObses = frostObses
	handlers.DbFile = *dbPath
	handlers.Class2FrostObses = class2FrostObses
	handlers.UserDBPath = *userDBPath

	http.HandleFunc("/roadlabels/camlist", makeHandler(handlers.CamlistHandler))
	http.HandleFunc("/roadlabels/allcams", makeHandler(handlers.AllCamsHandler))
	http.HandleFunc("/roadlabels/thumbs", makeHandler(handlers.ThumbsPageHandler))
	http.HandleFunc("/roadlabels/labeledthumb", makeHandler(handlers.LabelThumbHandler))
	http.HandleFunc("/roadlabels/inputlabel", makeHandler(handlers.InputLabelHandler))
	http.HandleFunc("/roadlabels/labeledimage", makeHandler(handlers.LabelImageHandler))
	http.HandleFunc("/roadlabels/showlabels", makeHandler(handlers.ShowLabelsHandler))
	http.HandleFunc("/roadlabels/frost_based_examples", makeHandler(handlers.FrostObsHandler))
	http.HandleFunc("/roadlabels/imagepagingapi", makeHandler(handlers.ImagePagingHandler))
	http.HandleFunc("/roadlabels/iceimagepagingapi", makeHandler(handlers.IceImagePagingHandler))
	http.HandleFunc("/roadlabels", makeHandler(handlers.FrontPageHandler))

	// Serve the label dat abase for backup and training. Nothing secret here . Could be read access for everyone
	http.HandleFunc("/roadlabels/db", handlers.DbDownloadHandler)
	http.HandleFunc("/roadlabels/userdb", handlers.UserDBDownloadHandler)

	// Serve javascript files from embedded fs at /roadlabels/js
	http.FileServer(http.FS(staticFiles))
	http.Handle("/roadlabels/js/", http.StripPrefix("/roadlabels/", http.FileServer(http.FS(staticFiles))))
	http.Handle("/roadlabels/css/", http.StripPrefix("/roadlabels/", http.FileServer(http.FS(staticFiles))))

	// TODO: Read from objectstore
	http.Handle("/roadlabels/nodata/", http.StripPrefix("/roadlabels/nodata/", http.FileServer(http.Dir("/lustre/storeB/project/metproduction/products/webcams"))))
	http.HandleFunc("/", redirect)

	log.Fatal(http.ListenAndServe(portStr, nil))
}
