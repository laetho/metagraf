# metaGraf

**metaGraf** provides a generic datastructure specification of information and metadata about a software component.
 
It is inspired by the <a href="https://12factor.net">twelve-factor app</a> 
methodology to aid automation tasks or decisions about a component or collection of compoenents.  

**metaGraf** operates on a individual or collections of metagraf specifications (software components) 
to produce metadata or aggregated metadata to support  your toolchain or pipelines with information that can be acted upon.

## Status

Current focus is experimentation with the model and trying out some use cases
that might help communicate what this structure and tooling can accomplish.

This repository contains the WIP specification and a tool called **mg** 
that consumes the specification or specifications and turns that into 
actionable items or kubernets resources.

The repository will split in two in the future. One for the datastructure specification
and one for the **mg** utility.


# Direction

Since cloud-native now eats the world, the goal is to enable building 
Kubernetes Operators/Controllers that act on the metadata and 
collections of metadata. The structure so far is also inspired by a 
Kubernetes resource so a metaGraf could be a CRD. 

# mg

A tool that understands metaGraf specifications.

TODO: Update this section with actual examples and a video of usage.

You can use the example collection provided to experiment. It produces output like 
this if the resulting file is used with the **dot** utility from graphviz.

<img src="https://github.com/laetho/metagraf/raw/master/example.png" alt="Example graph">

## Building

You'll need a working go lang environment. Consult with the internet on how to do that.
The project also use "dep" for managing go lang dependencies.

Get the source:

    git clone git@github.com:laetho/metagraf.git
   
Download vendored code with dep:

    dep ensure
    
To build the mg utility go into the mg catalog and build it:

    cd mg/
    go build

# Acknowledgement

A thank you to my current employer <a href="https://www.norsk-tipping.no">Norsk Tipping AS</a>,
for letting me share this work under an open source license.
