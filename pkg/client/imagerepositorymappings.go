package client

import (
	imageapi "github.com/openshift/origin/pkg/image/api"
)

// ImageRepositoryMappingsNamespacer has mathods to work with ImageRepositoryMapping resources in a namespace
type ImageRepositoryMappingsNamespacer interface {
	ImageRepositoryMappings(namespace string) ImageRepositoryMappingInterface
}

// ImageRepositoryMappingInterface exposes methods on ImageRepositoryMapping resources.
type ImageRepositoryMappingInterface interface {
	Create(mapping *imageapi.ImageRepositoryMapping) error
}

// imageRepositoryMappings implements ImageRepositoryMappingsNamespacer interface
type imageRepositoryMappings struct {
	r  *Client
	ns string
}

// newImageRepositoryMappings returns an imageRepositoryMappings
func newImageRepositoryMappings(c *Client, namespace string) *imageRepositoryMappings {
	return &imageRepositoryMappings{
		r:  c,
		ns: namespace,
	}
}

// Create create a new imagerepository mapping on the server. Returns error if one occurs.
func (c *imageRepositoryMappings) Create(mapping *imageapi.ImageRepositoryMapping) error {
	return c.r.Post().Namespace(c.ns).Path("imageRepositoryMappings").Body(mapping).Do().Error()
}
