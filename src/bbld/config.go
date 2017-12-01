package bbld

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var cfg AppCfg

type AppCfg struct {
	AudioBook AudioBook
}

type AudioBook struct {
	Fb2FileMask string
	Fb2Delete   bool
	ZipFileMask string
	ZipRename   bool
}

func LoadConfig() {
	b, err := ioutil.ReadFile("bookBuilder.yaml")
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
	ioutil.WriteFile("bookBuilder.yaml", d, 0644)
}

func init() {
	cfg = AppCfg{}
	LoadConfig()
	SaveConfig()
}

//
//sweaters := Inventory{"wool", 17}
//tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
//if err != nil { panic(err) }
//err = tmpl.Execute(os.Stdout, sweaters)
//if err != nil { panic(err) }
