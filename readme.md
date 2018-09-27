# metaGraf

**metaGraf** provides a generic structure of metadata about a software component.
 
**metaGraf** is inspired by the <a href="https://12factor.net">twelve-factor app</a> 
methodology to aid automation tasks or decisions about a component or collection of compoenents.  

**metaGraf** operates on an individual or collections of metagraph(s)
 (software components) to produce metadata or aggregated metadata to support 
 your toolchain or pipelines with information that can be acted upon.

A aggregated service graph can be created and used to automate entire
environments. A collection of metaGraf's should form the declarative state
necessary for a reconciliation loop on a complete solution or environment. 

This will also aid in documentation of software components and their dependencies.

## Status

Current focus is experimentation with the model and trying out some use cases
that might help communicate what this structure and tooling might accomplish.

This repository contains the WIP specification and a tool called **mg** 
that consumes the specification or specifications and turns that into 
actionable items or kubernets resources.

## Background

metaGraf tries to experiment with a structure for describing software components
and how that information can be used to assist CI and CD pipelines, developers,
architects, operations and a organization as a whole.

I have not found many projects that solve the complexities of
managing software components in an environment similar to the scope
of metaGraf.

The <a href="https://getnelson.github.io/nelson/">Nelson</a> project is
the closest thing I have found (after some pointers). It has a concept of
topology map (graph) of deployments in an environment.

<a href="http://ddd.ward.wiki.org/view/welcome-visitors/view/ward-cunningham">Ward Cunningham</a> 
is dabbling with something in this space at a broader scope: 
http://ddd.ward.wiki.org/view/about-the-el-dorado-project/

If anyone is  interested in this subject, please reach out and hopefully
get a discussion going. Input and suggestions are always welcome.
  

Direction
-
Since cloud-native now eats the world, the goal is to enable building 
Kubernetes Operators/Controllers that act on the metadata and 
collections of metadata. The structure so far is also inspired by a 
Kubernetes resource so a metaGraf could be a CRD. 

mg
-
A tool that understands metaGraf specifications.

TODO: Update this section with actual examples and a video of usage.

You can use the example collection provided to experiment. It produces output like 
this if the resulting file is used with the **dot** utility from graphviz.

<img src="https://github.com/laetho/metagraf/raw/master/example.png" alt="Example graph">

Acknowledgement
-
A shout out to my current employer <a href="https://www.norsk-tipping.no">Norsk Tipping AS</a>,
for letting me share this work under an open source license.


