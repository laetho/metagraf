# metaGraf Datastructure


The metaGraf datastructure is inspired by a kubernetes resource (kind). Doing this in JSON instead of YAML.


## Metadata

Follows the Kubernetes metadata specification.

### Labels

```
    "labels": {
      "component": "ComponentName"
    }
```




### Annotations

While the **Spec** structure is rigid, you can add custom organizational or solution information 
about a component using annotations, that your tooling may glean knowledge from or react on. 

```
    "annotations": {
      "myapp.example.com/application-type": "frontend",
      "myapp.example.com/maintainer": "John Doe <john@example.com>"
    }
```


## Spec

* Version needs to be a valid SemVer specification version. Vill get reduced to Major, Minor and Patch during evaluations and comparisions. 

### Resources

The resources section in the file describes a needed or optional attached resource.

```$go
type Resource struct {
	Name     	string	`json:"name"`
	Type     	string	`json:"type"`
	External 	bool    `json:"external"`
	User 		string	`json:"user,omitempty"`
	Secret	    string	`json:"secret,omitempty"`
	SecretType  string  `json:"secrettype,omitempty"`
	Semop		string	`json:"semop"`
	Semver  	string	`json:"semver"`
	Required 	bool	`json:"required"`
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

This section of the specification is split in two local and external.

* `local` Environment variables that needs to be set locally. Example: Where to 
get centrally managed config.
* `external` Environment variables that come from some configuration
management solution. This would just reference a unique key and downstream 
processing should produce the desired value or be overridden explicitly in 
a deployment.

#### Local

Environment variables local to the component.

#### External

The external environment variables are split in two sections and introduces 
the concept of which pre-existing variables it **consumes** and what variables
it **introduces**.

* `introduces` Variables introduced by component
* `consumes` Variables consumed by component introduced by other components.

### Config

This section is the spec is for defining traditional configuration variables group by `filename`.

## Status  

WIP

## Full example