# metaGraf information inside a container

By default the metaGraf utility will assume that the following directory
strucutre exist inside a metaGraf powered container or buildImage.

```text
/mg
/mg/config
```

It's up to the container image implementation to use the following information
to produce a working image. metaGraf only standardizes how that information will
be presented

## /mg



### /mg/config