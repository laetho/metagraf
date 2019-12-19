package metagraf

import (
	"encoding/json"
	log "k8s.io/klog"
	"io/ioutil"
	"os"
)

func Store(filepath string, mg *MetaGraf) {
	b, err := json.MarshalIndent(mg, "","  ")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(filepath, b, 0644 )
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
