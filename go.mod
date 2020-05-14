module metagraf

go 1.13

require (
	contrib.go.opencensus.io/exporter/ocagent v0.6.0 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.1.0 // indirect
	contrib.go.opencensus.io/exporter/stackdriver v0.13.0 // indirect
	github.com/blang/semver v3.5.1+incompatible
	github.com/coreos/prometheus-operator v0.38.0
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/go-containerregistry v0.0.0-20200212224832-c629a66d7231 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/kubernetes-sigs/application v0.8.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/gomega v1.8.1 // indirect
	github.com/openshift/api v0.0.0-20200116145750-0e2ff1e215dd
	github.com/openshift/client-go v0.0.0-20200116152001-92a2713fa240
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.1
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sys v0.0.0-20200116001909-b77594299b42 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/ini.v1 v1.51.1 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	istio.io/api v0.0.0-20200208020912-9564cdd03c96
	k8s.io/api v0.17.3
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	knative.dev/pkg v0.0.0-20200214073140-d8b36f359325
	knative.dev/serving v0.12.1
	sigs.k8s.io/application v0.8.1
	sigs.k8s.io/controller-runtime v0.4.0 // indirect
)

replace k8s.io/client-go v12.0.0+incompatible => k8s.io/client-go v0.17.3
