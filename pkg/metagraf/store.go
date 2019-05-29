package metagraf

import (
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
	"os"
)

func Store(filepath string, mg *MetaGraf) {
	b, err := json.MarshalIndent(mg, "","\t")
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(filepath, b, 0644 )
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
}
