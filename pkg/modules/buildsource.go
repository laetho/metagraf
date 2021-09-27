package modules

import (
	"github.com/laetho/metagraf/internal/pkg/params"
	"github.com/laetho/metagraf/pkg/metagraf"
	buildv1 "github.com/openshift/api/build/v1"
	corev1 "k8s.io/api/core/v1"
)

// Generate build sources

func ConstructBuildSource(mg *metagraf.MetaGraf) buildv1.BuildSource {
	var buildsource buildv1.BuildSource
	if len(mg.Spec.BaseRunImage) > 0 && len(mg.Spec.Repository) > 0 {
		buildsource = genBinaryBuildSource()
	} else if len(mg.Spec.BuildImage) > 0 && len(mg.Spec.BaseRunImage) < 1 {
		buildsource = genGitBuildSource(mg)
	}
	return buildsource
}

func genBinaryBuildSource() buildv1.BuildSource {
	return buildv1.BuildSource{
		Type:   "Source",
		Binary: &buildv1.BinaryBuildSource{},
	}
}

func genGitBuildSource(mg *metagraf.MetaGraf) buildv1.BuildSource {

	// When using the Git repository as a source without specifying the ref field,
	//OpenShift Enterprise performs a shallow clone (--depth=1 clone).
	//That means only the HEAD (usually the master branch) is downloaded.
	//This results in repositories downloading faster, including the commit history.
	// https://docs.openshift.com/enterprise/3.2/dev_guide/builds.html#source-code

	var branch string
	if len(params.SourceRef) > 0 {
		branch = params.SourceRef
	} else {
		branch = mg.Spec.Branch
	}

	bs := buildv1.BuildSource{
		Type: "Git",
		Git: &buildv1.GitBuildSource{
			URI: mg.Spec.Repository,
			Ref: branch,
		},
	}

	if len(mg.Spec.RepSecRef) > 0 {
		bs.SourceSecret = &corev1.LocalObjectReference{
			Name: mg.Spec.RepSecRef,
		}
	}

	return bs
}
