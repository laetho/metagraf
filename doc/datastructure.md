# MetaGraf datastructure


The MetaGraf datastructure is inspired by a kubernetes resource (kind).


## Metadata

### Annotations

While the **Spec** structure is rigid, you can add custom information 
about a component using annotations. That your tooling may react to. 

```
    "annotations": {
      "myapp.example.com/application-type": "frontend",
      "myapp.example.com/maintainer": "John Doe <john@example.com>"
    }
```


## Spec


### Resources

The resources section in the file describes a needed attached resource. 

### Environment

This section of the specification is split in two local and external.

* `local` Environment variables that needs to be set locally. Example: Where to get centrally managed config.
* `external` Environment variables that come from some configuration management solution. This would just reference a unique key and downstream processing should produce the desired value or be overridden explicitly in a deployment.

#### Local

#### External

## Status  

WIP

## Full example