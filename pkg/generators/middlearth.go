package generators

// Middlearth is a Norsk Tipping internal solution for configuration management of
// Jave applications on WebSphere Liberty

import (
	"encoding/json"
	"fmt"
	semver "github.com/blang/semver"
	"metagraf/pkg/metagraf"
	"strconv"
)

type Middlearth struct {
	AppName        string            `json:"appName"`
	AppExploded    bool              `json:"appExploded,omitempty"`
	Binary         string            `json:"binary"`
	Contextroot    string            `json:"contextroot"`
	SelftestURI    string            `json:"selftestUri"`
	Type           string            `json:"type"`
	Version        string            `json:"version"`
	LibertyVersion string            `json:"libertyVersion"`
	LastEdited     string            `json:"lastEdited,omitempty"`
	JavaHome       string            `json:"javaHome"`
	Jolokia        bool              `json:"jolokia,omitempty"`
	SessionCookie  string            `json:"sessionCookie,omitempty"`
	Sha1           string            `json:"sha1,omitempty"`
	Server         MEServer          `json:"server,omitempty"`
	Logging        MELogging         `json:"logging"`
	Features       map[string]string `json:"featurs,omitempty"`
	Certs          []MECert          `json:"certs,omitempty"`
	Datasources    []MEDataSource    `json:"datasources,omitempty"`
	Jms            []MiddlearthJMS   `json:"jms,omitempty"`
	JvmParams      map[string]string `json:"jvmParams,omitempty"`
	SpecialTags    map[string]string `json:"specialTags,omitempty"`
	Files          []MEFile          `json:"files,omitempty"`
	WasVars        map[string]string `json:"wasVars,omitempty"`
}

type MEServer struct {
	MaxHeap   int32  `json:"maxHeap,omitempty"`
	MinHeap   int32  `json:"minHeap,omitempty"`
	HttpPort  int32  `json:"httpPort,omitempty"`
	HttpsPort int32  `json:"httpsPort,omitempty"`
	Host      string `json:"host,omitempty"`
}

type MELogging struct {
	LogDir      string `json:"logdir,omitempty"`
	MaxFileSize int32  `json:"maxFileSize,omitempty"`
	MaxFiles    int32  `json:"maxFiles,omitempty"`
}

type MECert struct {
	Alias string `json:"alias"`
	Url   string `json:"url"`
}

type MEDataSource struct {
	Name     string `json:"name"`
	MaxCon   int32  `json:"maxCon"`
	MinCon   int32  `json:"minCon"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

type MiddlearthJMS struct {
	conMgr []struct {
		id          string
		maxPoolSize int32
	}
	qcf []struct {
		jndiName      string
		topic         bool
		conMgr        string
		userName      string
		clientId      string
		transportType string
		hostName      string
		port          int16
		channel       string
		qmgr          string
	}
	queues []struct {
		id       string
		jndiName string
		qname    string
		qmgr     string
	}
	activationSpec []struct {
		EJBPath       string
		qRefId        string
		transportType string
		qmgr          string
		hostName      string
		port          int16
		channel       string
	}
}

type MEFile struct {
	Source string `json:"source"`
	Dest   string `json:"dest"`
}

func MiddlearthApp(mg *metagraf.MetaGraf) {

	mea := Middlearth{}

	sv, err := semver.Parse(mg.Spec.Version)
	if err != nil {
		fmt.Println(err)
	}

	mea.AppName = mg.Metadata.Name + "v" + strconv.FormatUint(sv.Major, 10)
	mea.Binary = "http://middlearth.prod01.norsk-tipping.no:8081/getapp/no.norsktipping.services/" + mg.Metadata.Name + "-" + mg.Spec.Version + ".war"
	mea.Contextroot = "/" + mea.AppName
	mea.SelftestURI = "/selftest"

	mea.Jolokia = false

	mea.Logging.LogDir = "/var/log/websphere"
	mea.Logging.MaxFiles = 3
	mea.Logging.MaxFileSize = 100

	b, err := json.Marshal(mea)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

}
