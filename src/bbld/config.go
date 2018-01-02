package bbld

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const configFileName = "bookBuilder.yaml"

var cfg AppCfg

type AppCfg struct {
	AudioBook        AudioBook
	Fb2Parsing       Fb2Parsing
	WritingBookPaths map[string]string
}

type Fb2Parsing struct {
	BookTitleIsLastWords bool
}

type AudioBook struct {
	CoverFileMask       string
	CoverWriteFirstOnly bool
	Fb2FileMask         string
	Fb2Delete           bool
	TxtFileMask         string
	ZipFileMask         string
	ZipRename           bool
	RewriteBook         bool
}

func getConfigFile() string {
	dir, _ := filepath.Abs("./")
	flName := filepath.Join(dir, configFileName)
	if FileExists(flName) {
		return flName
	} else {
		log.Println("Not found config file " + flName)
	}

	dir, _ = os.Executable()
	dir, _ = filepath.Split(dir)
	flName = filepath.Join(dir, configFileName)
	if FileExists(flName) {
		return flName
	} else {
		log.Println("Not found config file " + flName)
	}

	log.Panic("Not found config file " + configFileName)
	return ""
}

func LoadConfig() {
	b, err := ioutil.ReadFile(getConfigFile())
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	err = yaml.Unmarshal([]byte(b), &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

}

func SaveConfig() {
	d, err := yaml.Marshal(&cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
	ioutil.WriteFile(getConfigFile(), d, 0644)
}

func (me AppCfg) GetWritingBookPath() string {
	hostname, _ := os.Hostname()
	if dir, ok := cfg.WritingBookPaths[hostname]; ok {
		return dir
	} else {
		return ""
	}
}

func init() {
	cfg = AppCfg{}
	LoadConfig()
	hostname, _ := os.Hostname()
	if cfg.WritingBookPaths == nil {
		cfg.WritingBookPaths = make(map[string]string)
	}
	if _, ok := cfg.WritingBookPaths[hostname]; !ok {
		cfg.WritingBookPaths[hostname] = "c:\\temp\\12"
	}
	SaveConfig()
}
