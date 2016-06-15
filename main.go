package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/crowdmob/goamz/aws"
	"github.com/spf13/cobra"
)

var configFile string

type configuration struct {
	Port      int    `json:"port"`
	SecretKey string `json:"secret_key"`
	AccessKey string `json:"access_key"`
	Region    string `json:"region"`
	Service   string `json:"service"`
	Endpoint  string `json:"endpoint"`
}

func main() {
	root := cobra.Command{
		Use: "es_proxy",
		Run: start,
	}
	root.PersistentFlags().StringVarP(&configFile, "config", "c", "config.json", "config file to use")

	if err := root.Execute(); err != nil {
		panic(err)
	}
}

func start(cmd *cobra.Command, args []string) {
	config, err := loadFromFile(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	region, ok := aws.Regions[config.Region]
	if !ok {
		keys := make([]string, 0, len(aws.Regions))
		for k := range aws.Regions {
			keys = append(keys, k)
		}

		log.Fatalf("Failed to parse the region: '%s', possible ones are %s", config.Region, strings.Join(keys, ","))
	}

	var auth aws.Auth
	if config.SecretKey == "" && config.AccessKey == "" {
		log.Println("Setting config from env")
		auth, err = aws.EnvAuth()
		if err != nil {
			log.Fatal("Failed to load the auth from the environment")
		}
	} else {
		log.Println("Setting config from file")
		auth = aws.Auth{
			SecretKey: config.SecretKey,
			AccessKey: config.AccessKey,
		}
	}

	signer := aws.NewV4Signer(auth, config.Service, region)

	url, _ := url.Parse(fmt.Sprintf("http://%s", config.Endpoint))
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Director = func(r *http.Request) {
		proxyURL, _ := url.Parse(fmt.Sprintf("http://%s%s", config.Endpoint, r.URL.RequestURI()))
		r.URL = proxyURL
		r.Host = config.Endpoint
		r.Header.Del("Connection")
		r.Header.Del("authorization")
		r.Header.Del("X-Forwarded-For")

		signer.Sign(r)
	}

	log.Printf("Serving on %v", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), proxy))
}

func loadFromFile(filename string) (*configuration, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := new(configuration)
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return validate(config)
}

func validate(config *configuration) (*configuration, error) {
	if config.Endpoint == "" {
		return nil, errors.New("Must provide an endpoint")
	}
	if config.Region == "" {
		return nil, errors.New("Must provide a region")
	}
	if config.Service == "" {
		return nil, errors.New("Must provide a region")
	}
	return config, nil
}
