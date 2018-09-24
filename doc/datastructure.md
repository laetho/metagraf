# metaGraf Datastructure

The metaGraf datastructure is inspired by a kubernetes resource (kind). This document
shows examples in JSON. Long term it might become a CRD itself.

Examples in JSON are stubs of the complete spec. For complete examples take a look
at the examples provided in the repository.

Follows the Kubernetes metadata specification.
```json
{
  "kind": "MetaGraf",
  "version": "v1alpha1",
  "metadata": {},
  "spec" : {}
}
```

## Metadata

Follows the Kubernetes metadata specification.

```json
{
  "metadata": {
    "name":"ComponentName",
    "annotations": {
      "example.com/my_annotation": "my value",
      "example.com/another_annotation" : "another value"
    },
    "labels": {
      "app" : "ComponentName"
    }
  }
}
```


### Labels

```json
{
    "labels": {
      "label": "labelvalue"
    }
}
```

For all resources belonging to a component it's advisable to label it with the component name:

```json
{
    "labels" : {
      "component" : "component name"
    }
}
```


### Annotations

While the **Spec** structure is rigid, you can add custom organizational or solution information 
about a component using annotations, that your tooling may glean knowledge from or react on. 

```json
{
    "annotations": {
      "myapp.example.com/application-type": "frontend",
      "myapp.example.com/maintainer": "John Doe <john@example.com>"
    }
}
```


## Spec

```json
{
 "spec": {
    "version": "4.2.3",
    "type": "service",
    "description": "Some kind of software component.",
    "repository" : "Repository URL",
    "branch": "master",
    "repsecref": "Referense to name of secret to pull source code",
    "buildimage" : "URL to image on registry",
    "baserunimage" : "URL to image on registry",
    "resources": [],
    "environment": {},
    "config" : []
    }
}
```

| Field | - | Description |
|-------|---|-------------|
| version | required | Need to be a valid SemVer specification version. Might be reduced to major, minor, patch during evaluations and comparisons.|
| type | required | There are currently two component types at the moment: service and datastore.|
| description | required | A textual description of the software component.|
| repository | optional | Repository URL.|
| branch | optional | Branch name.| 
|repsecref|optional|Reference to secret for accessing the repository.|
|buildimage|optional|Registry url to the image we build from/on. If no buildimage is provided we won't attempt to build anything.|
|baserunimage|optional|Registry url to the runtime image. Used for binary builds or finished components.|
|resources|optional|An array of Resource records. See Resources section.|
|environment|optional|Environment variables. See Environment section.|
|config|optional|Array of Config records. See Config section.|

* If only repository url is provided it indicates a build from source and Dockerfile.
* If only repository and buildimage is provided it indicates a s2i build image or similar.
* If buildimage and baserunimage is provided it indicates a binary build with 
* If only a baserunimage is provided it indicates instrumentation of a prebuilt component. 
The scenarios here needs work.

### Resources

The resources section in the file describes a needed or optional attached resource.

```go
type Resource struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	External    bool    `json:"external"`
	User        string  `json:"user,omitempty"`
	Secret      string  `json:"secret,omitempty"`
	SecretType  string  `json:"secrettype,omitempty"`
	Semop       string  `json:"semop"`
	Semver      string  `json:"semver"`
	Required    bool    `json:"required"`
	Url         string  `json:"url,omitempty"`
}
```


There are currently two types (Type) of resources:

* *service* generally a http rest service
* *datasource* a database or similar backend with a connection string



#### Implicit Secrets

If the User field is filled out and the SecretRef field is empty we, have an implicit
secret. A convention for creating or accessing the secret must be created.

In an 

#### Explicit Secrets

If the SecretRef field is filled out it means there is a explicit secret related to this 
 

### Environment

```json
{
    "environment": {
      "build": [],
      "local": [],
      "external": []
      }
}
```

This section of the specification is split in two local and external.
* `build` Build level environment variables
* `local` Environment variables that needs to be set locally. Example: Where to 
get centrally managed config.
* `external` Environment variables that come from some configuration
management solution. This would just reference a unique key and downstream 
processing should produce the desired value or be overridden explicitly in 
a deployment.


#### Build


#### Local

Environment variables local to the component.

#### External

The external environment variables are split in two sections and introduces 
the concept of which pre-existing variables it **consumes** and what variables
it **introduces**.

* `introduces` Variables introduced by component
* `consumes` Variables consumed by component introduced by other components.

### Config

This section is the spec is for defining traditional configuration variables 
group by `filename`. Options is an array of ConfigParam.

```json
{
    "config": [
      {
        "name": "name_of_file",
        "type": "parameters",
        "description": "",
        "options": []
      }
    ]
}
```

ConfigParam example: 
```json
  {
    "name": "parameter.name",
    "required": "true",
    "description": "Description of parameter",
    "type": "string",
    "default": "512"
  }
```


## Status  

WIP

## Full example