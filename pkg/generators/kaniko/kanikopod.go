package kaniko

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
	log "k8s.io/klog"

	"github.com/laetho/metagraf/pkg/metagraf"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	errNoDockerfile = "metaGraf specifies no Dockerfile, aborting build"
)

// Generator for the Application type
type KanikoPodGenerator struct {
	MetaGraf   metagraf.MetaGraf
	Properties metagraf.MGProperties
	Options    KanikoPodOptions
	Resource   corev1.Pod
}

type KanikoPodOption func(options *KanikoPodOptions)

type KanikoPodOptions struct {
	// Namespace for generated Pod
	Namespace string

	// Kaniko container image reference
	Image string

	// --destination Argument for Kaniko executor.
	DestinationArg string

	// --dockerfile Argument for Kaniko executor.
	// Path to Dockerfile in provided --context
	DockerfileArg string

	// --context Argument for Kaniko executor.
	// Se https://github.com/GoogleContainerTools/kaniko#kaniko-build-contexts for possible
	// usage scenarios.
	ContextArg string

}

// Instance of ApplicationOptions that can be used for propagating flags and addressable outside of pacage.
var KanikoPodOpts KanikoPodOptions

func init() {
	KanikoPodOpts.Image = "gcr.io/kaniko-project/executor:debug"
	KanikoPodOpts.DockerfileArg = "Dockerfile"
}

// Creates a new Options based on defaults from Opts and runs
// the functional Option methods against it from options.
func NewOptions(options ...KanikoPodOption) KanikoPodOptions {
	opts := KanikoPodOpts
	o := &opts
	for _, opt := range options {
		opt(o)
	}
	return *o
}

// Factory for creating a new KanikoPodGenerator
func NewKanikoPodGenerator(mg metagraf.MetaGraf, prop metagraf.MGProperties, options KanikoPodOptions) KanikoPodGenerator {
	g := KanikoPodGenerator{
		MetaGraf:   mg,
		Properties: prop,
		Options:    options,
	}

	if len(g.MetaGraf.Spec.Dockerfile) <= 0 {
		log.Fatal(errNoDockerfile)
	}
	g.Options.DockerfileArg = mg.Spec.Dockerfile
	g.Options.ContextArg = mg.Spec.Repository

	return g
}

func (g *KanikoPodGenerator) GenerateKanikoPod(name string) corev1.Pod {

	g.Resource = corev1.Pod{
		TypeMeta:   metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kaniko-"+g.MetaGraf.Name("",""),
			Namespace: g.Options.Namespace,
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: g.kanikoPod("kaniko-"+g.MetaGraf.Name("","")),
		},
	}

	return g.Resource
}

func (g *KanikoPodGenerator) podArgs() []string {
	args := []string{}

	if len(g.Options.DockerfileArg) > 0 {
		args = append(args, "--dockerfile="+g.Options.DockerfileArg)
	}
	if len(g.Options.DestinationArg) > 0 {
		args = append(args, "--destination="+g.Options.DestinationArg)
	}
	if len(g.Options.ContextArg) > 0 {
		args = append(args, "--context="+g.Options.ContextArg)
	}

	return args
}

func (g *KanikoPodGenerator) kanikoPod(name string) []corev1.Container {
	var containers []corev1.Container
	c := corev1.Container{
		Name:                     name,
		Image:                    g.Options.Image,
		Args:                     g.podArgs(),
		Env:                      g.MetaGraf.KubernetesBuildVars(),
		VolumeMounts:             nil,
		VolumeDevices:            nil,
	}
	return append(containers,c)
}

func (g *KanikoPodGenerator) MarshalToYaml() {
	MarshalToYaml(g.Resource)
}

func (g *KanikoPodGenerator) MarshalToJson() {
	MarshalToJson(g.Resource)
}

func MarshalToYaml(obj interface{}) {
	y, err := yaml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(y))
}

func MarshalToJson(obj interface{}) {
	j, err := json.MarshalIndent(obj,"", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}