module github.com/laetho/metagraf

go 1.16

require (
	github.com/argoproj/argo-cd v1.8.7
	github.com/blang/semver v3.5.1+incompatible
	github.com/containerd/continuity v0.0.0-20210208174643-50096c924a4e // indirect
	github.com/coreos/prometheus-operator v0.41.1
	github.com/crossplane/crossplane-runtime v0.9.0 // indirect
	github.com/crossplane/oam-kubernetes-runtime v0.0.9
	github.com/fsnotify/fsnotify v1.4.9
	github.com/ghodss/yaml v1.0.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/go-containerregistry v0.4.0
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/kubernetes-sigs/application v0.8.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/openshift/api v0.0.0-20200825174227-962ddb6aceab
	github.com/openshift/client-go v0.0.0-20200729195840-c2b1adc6bed6
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.1
	github.com/tidwall/gjson v1.8.1
	github.com/tidwall/pretty v1.2.0 // indirect
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/sys v0.0.0-20210108172913-0df2131ae363 // indirect
	golang.org/x/term v0.0.0-20201210144234-2321bbc49cbf // indirect
	gopkg.in/yaml.v3 v3.0.0-20200603094226-e3079894b1e8
	istio.io/api v0.0.0-20200208020912-9564cdd03c96
	k8s.io/api v0.19.10
	k8s.io/apimachinery v0.19.10
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.80.0
	k8s.io/kubectl v0.19.10 // indirect
	sigs.k8s.io/application v0.8.1
	sigs.k8s.io/controller-runtime v0.6.2 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.19.10
	k8s.io/api/settings/ => k8s.io/api/settings/ v0.19.10
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.10
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.10
	k8s.io/apiserver => k8s.io/apiserver v0.19.10
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.10
	k8s.io/client-go => k8s.io/client-go v0.19.10
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.10
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.10
	k8s.io/code-generator => k8s.io/code-generator v0.19.10
	k8s.io/component-base => k8s.io/component-base v0.19.10
	k8s.io/cri-api => k8s.io/cri-api v0.19.10
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.10
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.10
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.10
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.10
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.10
	k8s.io/kubectl => k8s.io/kubectl v0.19.10
	k8s.io/kubelet => k8s.io/kubelet v0.19.10
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.10
	k8s.io/metrics => k8s.io/metrics v0.19.10
	k8s.io/node-api => k8s.io/node-api v0.19.10
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.10
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.19.10
	k8s.io/sample-controller => k8s.io/sample-controller v0.19.10
)
