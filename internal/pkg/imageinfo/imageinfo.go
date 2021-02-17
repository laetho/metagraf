package imageinfo

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	corev1 "k8s.io/api/core/v1"
	"metagraf/internal/pkg/helpers/helpers"
	"metagraf/pkg/metagraf"
	"metagraf/pkg/mgver"
)

// Type alias of Config type from go-containerregistry, for creating our own method sets.
type Info v1.Config

// Returns the oci or docker image config section based on images referenced
// in the metaGraf specification.
func MGImageInfo(mg *metagraf.MetaGraf) (Info, error) {
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
		return Info{}, err
	}
	return config, nil
}

// Returns the oci or docker image config from the imageref argument or an error.
func ImageInfo(imageref string) (Info, error) {

	ref, err := name.ParseReference(imageref)
	if err != nil {
		return Info{},err
	}

	var options []remote.Option

	options = append(options, remote.WithUserAgent("mg v"+ mgver.GitTag))
	options = append(options, remote.WithAuth(authn.Anonymous))

	img, err := remote.Image(ref, options...)
	if err != nil {
		return Info{}, err
	}

	config, _ := img.ConfigFile()
	return Info(config.Config), nil
}


// Turns the OCI or Docker image volume information into a
// slice of corev1.Volume{}
func (info Info) ImageVolumes(nameprefix string) []corev1.Volume {
	var out []corev1.Volume

	// Volumes & VolumeMounts from base image into podspec
	for k := range info.Volumes {
		// Volume Definitions
		Volume := corev1.Volume{
			Name: nameprefix + helpers.PathToIdentifier(k),
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}
		out = append(out, Volume)
	}
	return out
}

// Turns the OCI or Docker image volume information into a
// slice of corev1.VolumeMount{}
func (info Info) ImageVolumeMounts(nameprefix string) []corev1.VolumeMount {
	var out []corev1.VolumeMount

	for k := range info.Volumes {
		VolumeMount := corev1.VolumeMount{
			MountPath: k,
			Name:      nameprefix + helpers.PathToIdentifier(k),
		}
		out = append(out, VolumeMount)
	}
	return nil
}