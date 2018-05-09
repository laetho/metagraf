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

In other words, determining what needs to be present or changed to 
fulfil the explicit dependencies of the new component entering an
environment.

Another goal is to aid in documentation of software components and
their dependencies on a component level.

Background
-
metaGraf is currently a research project and a place to experiment
with a structure for describing software components and how that
information can be used to assist CI and CD pipelines, developers,
architects, operations and a organization as a whole.

I have not found many projects that solve the complexities of
managing software components in an enviroment similar to the goals
of metaGraf.

The <a href="https://getnelson.github.io/nelson/">Nelson</a> project is
the closest thing I have found (after some pointers). It has a concept of
topology map (graph) of deployments in an environment.

<a href="http://ddd.ward.wiki.org/view/welcome-visitors/view/ward-cunningham">Ward Cunningham</a> 
is dabbling with something in this space at a broader scope: 
http://ddd.ward.wiki.org/view/about-the-el-dorado-project/

If anyone is  interested in this subject, please reach out and hopefully
 we can get a discussion going. Input and suggestions are always welcome.
  

Direction
-
Since cloud-native now eats the world, the goal is to enable building 
Kubernetes Operators/Controllers that act on the metadata and 
collections of metadata. The structure so far is also inspired by a 
Kubernetes resource so a metaGraf could be a CRD. 

mgraf
-
A little tool to help communicate what metagraf attempt to solve and basis for
further discussion.

Usage:

> mgraf -c /path/to/collection/of/metagraphs 

Usage straight from source:

> go run mgraf -c /path/to/collection/of/metagraphs

You can use the example collection provided to experiment. It produces output like 
this if the resulting file is used with the **dot** utility from graphviz.

<img src="https://github.com/laetho/metagraf/raw/master/example.png" alt="Example graph">

Acknowledgement
-

A thank you to my current employer <a href="https://www.norsk-tipping.no">Norsk Tipping AS</a>, for letting me share this work under an
open source license.


