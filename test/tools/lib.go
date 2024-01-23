package toolstester

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	"github.com/petri-nets/workflow/flow"
	"github.com/petri-nets/workflow/wfmod"
	"gopkg.in/yaml.v2"
)

func GetCurrentPath() string {
	currentPath, err := os.Getwd()
	if err != nil {
		panic("get current path error: " + err.Error())
	}
	return currentPath
}

func ParseTestYaml(yamlPath string) {
	yamlFile, err := os.ReadFile(yamlPath)

	if err != nil {
		panic("config file read err: " + err.Error())
	}

	pipeline := flow.Pipeline{}

	err = yaml.Unmarshal(yamlFile, &pipeline)
	if err != nil {
		panic("yaml config parse error:" + err.Error())
	}

	fmt.Println(pipeline.AppID)

	parser := flow.NewParser(&pipeline)

	parser.Start()
}

func TruncateData(db *gorm.DB) {
	db.Exec("truncate table `wf_workflows`;")
	db.Exec("truncate table `wf_places`;")
	db.Exec("truncate table `wf_transitions`;")
	db.Exec("truncate table `wf_arcs`;")
	db.Exec("truncate table `wf_cases`;")
	db.Exec("truncate table `wf_tokens`;")
	db.Exec("truncate table `wf_workitems`;")
}

func FindTestYamlFiles() {
	files, err := filepath.Glob(GetCurrentPath() + "/test/data/**/*.yaml")
	for _, file := range files {
		ParseTestYaml(file)
		// fileName := filepath.Base(file)
		// nameNoExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

		flow.StartWorkflow(1, 1, wfmod.WfContextType{"workspace": "ww01"}, "edenzou")
	}
	fmt.Println(files, err)
}
