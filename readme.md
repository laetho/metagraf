metaGraf
========

**metaGraf** provides a generic and implementation agnostic
"structure" of metadata about a software component. **metagraf**
is inspired by the <a href="https://12factor.net">twelve-factor app</a> guidelines to 
aid automation tasks or decisions about a component or collection of compoenents.

**metaGraf** operates on an individual or collections of metagraph(s)
 (software components) to produce aggregated metadata to support your
toolchain or pipelines with information that can be acted upon.

One of the goals of **metaGraf** is to indentify missing nodes or edges
(components) when comparing a running enviroment with the graph/branch
of a component not currently deployed or a new version of an existing
component. Desired state vs existing state.

In other words, determining what needs to be present to fulfil the
explicit dependencies of the new component entering the enviroment.

Another goal is to aid in documentation of software components and
their dependencies on a component level.

Background
-
metaGraf is currently a research project and a place to experiment
with a structure for describing software components and how that
information can be used to assist CI and CD pipelines, developers,
architects, operations and a organization as a whole.

I have not found any projects that aims to solve the complexities of
managing software components at scale in this way. If a similar thing
exist or people are working on it, please let me know.

Apparently the paragraph above is no longer true. It seems 
<a href="http://ddd.ward.wiki.org/view/welcome-visitors/view/ward-cunningham">Ward Cunningham</a> 
is dabbling with something in this space at a broader scope: 
http://ddd.ward.wiki.org/view/about-the-el-dorado-project/

Direction
-
Since cloud-native now eats the world, the goal is to enable building 
Kubernetes Operators/Controllers that act on the metadata and 
collections of metadata.


Acknowledgements
-

A thank you to my current employer <a href="https://www.norsk-tipping.no">Norsk Tipping AS</a>, for letting me share this work under
open source license.


