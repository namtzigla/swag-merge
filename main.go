package main

import (
	"flag"
	"fmt"
	"os"

	"io/ioutil"

	"encoding/json"

	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/imdario/mergo"
)

var (
	verbose = flag.Bool("v", false, "Verbose")
	outfile = flag.String("out", "api.swagger.json", "Output file")
	desc    = flag.String("desc", "", "New value for info.description field")
	title   = flag.String("title", "", "New value for info.title field")
	version = flag.String("version", "", "New value for info.version field")
)

func check(dst, src map[string]interface{}) error {
	if len(dst) == 0 {
		return nil
	}
	if !reflect.DeepEqual(dst["schemes"], src["schemes"]) {
		return fmt.Errorf("Supported schemes are not the same")
	}
	if !reflect.DeepEqual(dst["basePath"], src["basePath"]) {
		return fmt.Errorf("BasePath's not the same")
	}
	if !reflect.DeepEqual(dst["host"], src["host"]) {
		return fmt.Errorf("Hosts not the same")
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [Options] File1 [File2 [File3 ...]]\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	dst := map[string]interface{}{}
	flag.Parse()
	// read files
	for fileNo := range flag.Args() {
		log.WithField("File", flag.Arg(fileNo)).Info("Processing file")
		buf, err := ioutil.ReadFile(flag.Arg(fileNo))
		if err != nil {
			log.WithError(err).Fatal("Fatal error reading file!")
		}
		var m map[string]interface{}
		err = json.Unmarshal(buf, &m)
		if err != nil {
			log.WithError(err).Fatal("Fatal error parsing json!")
		}
		err = check(dst, m)

		if err != nil {
			log.WithError(err).Fatal("Fatal error while checking merge preconditions!")
		}
		err = mergo.Map(&dst, m)
		if err != nil {
			log.WithError(err).Fatal("Fatal error merging json structures!")
		}
	}
	if *desc != "" || *title != "" || *version != "" {
		if *desc != "" {
			dst["info"].(map[string]interface{})["description"] = *desc
		}
		if *title != "" {
			dst["info"].(map[string]interface{})["title"] = *title
		}
		if *version != "" {
			dst["info"].(map[string]interface{})["version"] = *version
		}

	}
	buf, err := json.MarshalIndent(dst, " ", " ")
	if err != nil {
		log.WithError(err).Fatal("Fatal error converting strucures back to json!")
	}
	ioutil.WriteFile(*outfile, buf, 0644)
	log.Info("Done!")
}
