package kaniko

import (
	"context"
	"encoding/json"
	"io"

	"github.com/ghodss/yaml"
	"github.com/laetho/metagraf/internal/pkg/k8sclient"
	log "k8s.io/klog"

	"github.com/laetho/metagraf/pkg/metagraf"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	errNoDockerfile = "metaGraf specifies no Dockerfile, aborting build"
)

type Generator interface {
	Create(obj interface{}) error
	Delete(obj interface{}) error
}

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

	// --cache Indicate if we want to use caching at all.
	Cache bool

	// Uses --cache-dir to set local directory as cache for base images. In our context
	// this will be a Volume in the metaGraf specifiation that becomes a PersistantVolume.
	CacheDir string

	// Option for skipping cert verification on image push to a registry.
	SkipTLSVerify bool

	// Option for skipping cert verification on image pulls from a registry.
	SkipTLSVerifyPull bool

	StdIn     bool
	StdInOnce bool
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

	if len(g.MetaGraf.Spec.Image) > 0 {
		g.Options.DestinationArg = g.MetaGraf.Spec.Image
	}

	return g
}

func (g *KanikoPodGenerator) Generate(name string) corev1.Pod {

	g.Resource = corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kaniko-" + g.MetaGraf.Name("", ""),
			Namespace: g.Options.Namespace,
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers:    g.kanikoPod("kaniko-" + g.MetaGraf.Name("", "")),
			Volumes:       append(g.MetaGraf.BuildSecretsToVolumes(), g.MetaGraf.Volumes()...),
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
	if g.Options.SkipTLSVerifyPull {
		args = append(args, "--skip-tls-verify-pull")
	}
	if g.Options.SkipTLSVerify {
		args = append(args, "--skip-tls-verify")
	}
	if g.Options.Cache {
		args = append(args, "--cache=true")
	}
	if g.Options.Cache && len(g.Options.CacheDir) > 0 {
		args = append(args, "--cache-dir="+g.Options.CacheDir)
		// Also cache RUN and copy layers.
		args = append(args, "--cache-copy-layers")
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
		VolumeMounts:             append(g.MetaGraf.BuildSecretsToVolumeMounts(),g.MetaGraf.VolumesToVolumeMounts()...),
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: corev1.TerminationMessageReadFile,
	}
	return append(containers, c)
}

func (g *KanikoPodGenerator) Create(obj corev1.Pod) error {
	client := k8sclient.GetCoreClient().Pods(g.Options.Namespace)

	result, err := client.Create(context.TODO(), &obj, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	log.Infof("Created Kaniko Build Pod: %v(%v)", result.Name, obj.Name)

	return nil
}

func(g *KanikoPodGenerator) Delete(obj corev1.Pod ) error {
	client := k8sclient.GetCoreClient().Pods(g.Options.Namespace)

	err := client.Delete(context.TODO(), obj.Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	log.Infof("Delete Kaniko Build Pod: %v", obj.Name)

	return nil
}


// Return a reference to a io.ReadCloser based on a Pod Log Request or an error.
func (g *KanikoPodGenerator) LogsReader(obj corev1.Pod) (*io.ReadCloser, error) {
	client := k8sclient.GetCoreClient().Pods(g.Options.Namespace)

	podLogOptions :=  corev1.PodLogOptions{
		Container:                    obj.Name,
		Follow:                       true,
		TailLines:                    nil,
	}

	podLogReq := client.GetLogs(obj.Name,&podLogOptions)
	stream, err := podLogReq.Stream(context.TODO())
	if err != nil {
		return nil, err
	}

	return &stream, nil
}

func (g *KanikoPodGenerator) ToYaml() ([]byte, error) {
	b, err := MarshalToYaml(g.Resource)
	return b, err
}

func (g *KanikoPodGenerator) ToJson() ([]byte, error) {
	b, err := MarshalToJson(g.Resource)
	return b, err
}

func MarshalToYaml(obj interface{}) ([]byte, error) {
	y, err := yaml.Marshal(obj)
	if err != nil {
		return []byte{}, err
	}
	return y, nil
}

func MarshalToJson(obj interface{}) ([]byte, error) {
	j, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return []byte{}, err
	}
	return j, nil
}
