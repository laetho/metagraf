metagraf
========

**metagraf** provides a generic and implementation agnostic "structure" of metadata about a software component. '''metagraf''' is inspired by the 12 factor app manifesto to aid automation tasks or decisions about a component or collection of compoenents.

**metagraf** operates on an individual or collections of metagraph(s) (software components) to produce aggregated metadata to support your toolchain or pipelines with information that can be acted upon.

One of the goals of metagraf is to indentify missing nodes or edges (components) when comparing a running enviroment with the graph/branch of a component not currently deployed. Desired state vs existing state.

In other words determining what needs to be present to fulfil the explicit dependencies of the new component entering the enviroment.

Another goal is to aid in documentation of software components and their dependencies on a component level.





