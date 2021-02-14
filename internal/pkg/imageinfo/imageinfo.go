package imageinfo

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"metagraf/pkg/metagraf"
)


// Returns the oci or docker image config section based on images referenced
// in the metaGraf specification.
func MGImageInfo(mg *metagraf.MetaGraf) (v1.Config, error) {
	var imageref string
	if len(mg.Spec.BaseRunImage) > 0 {
		imageref = mg.Spec.BaseRunImage
	} else if len(mg.Spec.BuildImage) > 0 {
		imageref = mg.Spec.BuildImage
	} else if len(mg.Spec.Image) > 0 {
		imageref = mg.Spec.Image
	}

	config, err := ImageInfo(imageref)
	if err != nil {
		return v1.Config{}, err
	}
	return config, nil
}


func ImageInfo(imageref string) (v1.Config, error) {

	ref, err := name.ParseReference(imageref)
	if err != nil {
		return v1.Config{},err
	}

	var options []remote.Option

	options = append(options, remote.WithUserAgent("mg"))
	options = append(options, remote.WithAuth(authn.Anonymous))

	img, err := remote.Image(ref, options...)
	if err != nil {
		return v1.Config{}, err
	}

	config, _ := img.ConfigFile()
	return config.Config, nil
}
