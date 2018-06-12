package generators

import (
"metagraf/internal/metagraf"
metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
)

type BuildConfig struct {
	metav1.TypeMeta			`json:",inline"`
	metav1.ObjectMeta		`json:"metadata"`
	BuildConfigStatus		`json:"status"`
}

type Trigger struct {
	Type string 		`json:"type"`
	Map map[string]string `json:",inline"`

}

type BuildConfigStrategy struct {
	Type string		`json:"type"`
	SourceStrategy 	BuildConfigSourceStrategy `json:"sourceStrategy"`
}

type BuildConfigSourceStrategy struct {

}

type BuildConfigSpec struct {
	Triggers []Trigger  `json:"triggers,omitempty"`
	RunPolicy string	`json:"runPolicy"`
	Source map[string]string
	Strategy BuildConfigStrategy `json:"strategy"`
	Output struct {}
	Resources struct{}
	PostCommit struct{}
	NodeSelector string
	SuccessfulBuildsHistoryLimit int
	FailedBuildsHistoryLimit int

}

type BuildConfigStatus struct {
	LastVersion int64	`json:"lastVersion"`
}

func GenBuildConfig( mg metagraf.MetaGraf) {
	fmt.Println("do buildconfig generation")
}
