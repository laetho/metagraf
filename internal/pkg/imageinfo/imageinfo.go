package imageinfo

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)



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
