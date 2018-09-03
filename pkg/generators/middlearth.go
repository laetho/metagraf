package generators

// Middlearth is a Norsk Tipping internal solution for configuration management of
// Java applications on WebSphere Liberty. Much of this complexity should be hidden
// in the process of building the artifact/container image.

import (
	"encoding/json"
	"fmt"
	"github.com/blang/semver"
	"metagraf/pkg/metagraf"
	"strconv"
	"strings"
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
	Features       []string          `json:"features,omitempty"`
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
	ConMgr         []MEConMgr         `json:"conMgr,omitempty"`
	Qcf            []MEQcf            `json:"qcf,omitempty"`
	Queues         []MEQueue          `json:"queues,omitempty"`
	ActivationSpec []MEActivationSpec `json:"activationSpec,omitempty"`
}

type MEConMgr struct {
	Id          string `json:"id"`
	MaxPoolSize int32  `json:"maxPoolSize"`
}

type MEQcf struct {
	JndiName      string `json:"jndiName"`
	Topic         bool   `json:"topic"`
	ConMgr        string `json:"conMgr"`
	UserName      string `json:"userName"`
	ClientId      string `json:"clientId"`
	TransportType string `json:"transportType"`
	HostName      string `json:"hostName"`
	Port          int16  `json:"port"`
	Channel       string `json:"channel"`
	Qmgr          string `json:"qmgr"`
}

type MEQueue struct {
	Id       string `json:"id"`
	JndiName string `json:"jndiName"`
	Qname    string `json:"qname"`
	Qmgr     string `json:"qmgr"`
}

type MEActivationSpec struct {
	EJBPath       string `json:"EJBPath"`
	QRefId        string `json:"qRefId"`
	TransportType string `json:"transportType"`
	Qmgr          string `json:"qmgr"`
	HostName      string `json:"hostname"`
	Port          int16  `json:"port"`
	Channel       string `json:"channel"`
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

	// Defaults managed by generator
	mea.Contextroot = "/" + mea.AppName
	mea.SelftestURI = "/selftest"
	mea.Jolokia = false
	mea.JavaHome = "/opt/ibm/latestsdk8"
	mea.LibertyVersion = "17.0.0.4-open"
	mea.Logging.LogDir = "/var/log/websphere"
	mea.Logging.MaxFiles = 3
	mea.Logging.MaxFileSize = 100

	for k, v := range mg.Metadata.Annotations {
		switch k {
		case "middlearth.norsk-tipping.no/jolokia":
			b, _ := strconv.ParseBool(v)
			mea.Jolokia = b
		case "middlearth.norsk-tipping.no/features":
			s := strings.Split(v, ",")
			for _, v := range s {
				mea.Features = append(mea.Features, strings.TrimSpace(v))
			}
		case "middlearth.norsk-tipping.no/libertyVersion":
			mea.LibertyVersion = v
		}
	}

	b, err := json.Marshal(mea)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

}
