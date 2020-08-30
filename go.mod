module metagraf

go 1.14

require (
	contrib.go.opencensus.io/exporter/ocagent v0.6.0 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.1.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.0 // indirect
	github.com/argoproj/argo-cd v1.7.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/coreos/prometheus-operator v0.41.1
	github.com/fvbommel/util v0.0.0-20180919145318-efcd4e0f9787 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/go-containerregistry v0.0.0-20200212224832-c629a66d7231 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kubernetes-sigs/application v0.8.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/oam-dev/oam-go-sdk v0.0.0-20200710215707-508886568e3f
	github.com/onsi/gomega v1.8.1 // indirect
	github.com/openshift/api v0.0.0-20200825174227-962ddb6aceab
	github.com/openshift/client-go v0.0.0-20200729195840-c2b1adc6bed6
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.1
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/text v0.3.3 // indirect
	gopkg.in/ini.v1 v1.51.1 // indirect
	istio.io/api v0.0.0-20200208020912-9564cdd03c96
	k8s.io/api v0.19.0-rc.2
	k8s.io/apimachinery v0.19.0-rc.2
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	knative.dev/pkg v0.0.0-20200214073140-d8b36f359325 // indirect
	knative.dev/serving v0.12.1 // indirect
	sigs.k8s.io/application v0.8.1
	sigs.k8s.io/controller-runtime v0.4.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.18.8
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.8
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.8
	k8s.io/apiserver => k8s.io/apiserver v0.18.8
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.8
	k8s.io/client-go => k8s.io/client-go v0.18.8
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.8
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.8
	k8s.io/code-generator => k8s.io/code-generator v0.18.8
	k8s.io/component-base => k8s.io/component-base v0.18.8
	k8s.io/cri-api => k8s.io/cri-api v0.18.8
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.8
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.8
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.8
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.8
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.8
	k8s.io/kubectl => k8s.io/kubectl v0.18.8
	k8s.io/kubelet => k8s.io/kubelet v0.18.8
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.8
	k8s.io/metrics => k8s.io/metrics v0.18.8
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.8
)
