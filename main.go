package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/gidyon/file-handlers/static"
	"net/http"
	"os"
	"strings"
)

const (
	siteCert = "certs/cert.pem"
	siteKey  = "certs/key.pem"
)

func main() {
	certFile := flag.String("cert", siteCert, "Path to public key")
	keyFile := flag.String("key", siteKey, "Path to private key")
	port := flag.String("port", ":443", "Port to serve files")
	root := flag.String("root", "static", "Root directory")
	dirs := flag.String("dirs", ".", "allowed directories relative to root")
	env := flag.Bool("env", false, "Whether to read parameters from env variables")
	insecure := flag.Bool("insecure", false, "Use insecure server")

	flag.Parse()

	if *env {
		*certFile = setIfEmpty(os.Getenv("TLS_CERT_FILE"), *certFile)
		*keyFile = setIfEmpty(os.Getenv("TLS_KEY_FILE"), *keyFile)
		*port = setIfEmpty(os.Getenv("PORT"), *port)
		*root = setIfEmpty(os.Getenv("ROOT_DIR"), *root)
		*dirs = setIfEmpty(os.Getenv("ALLOWED_DIRS"), *dirs)
	}

	_ = dirs
	// static file server
	staticHandler, err := static.NewHandler(&static.ServerOptions{
		RootDir: *root,
		Index:   "static/index.html",
		// AllowedDirs: strings.Split(*dirs, ","),
		NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/index.html")
		}),
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	http.Handle("/", staticHandler)

	logrus.Infof("server started on %v\n", *port)

	*port = ":" + strings.TrimPrefix(*port, ":")

	if *insecure {
		logrus.Fatalln(http.ListenAndServe(*port, nil))
	} else {
		logrus.Fatalln(http.ListenAndServeTLS(*port, *certFile, *keyFile, nil))
	}
}

func setIfEmpty(strCurrent, strFinal string) string {
	if strCurrent == "" && strFinal == "" {
		return ""
	}
	if strCurrent == "" {
		return strFinal
	}
	return strCurrent
}
