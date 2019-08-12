# metaGraf

**metaGraf** is a specification describing the necessary metadata and information about a software component for building it or running it or both. It's intended to be used as a source of truth for CI/CD or GitOps scenarios.
 
It is inspired by the <a href="https://12factor.net">twelve-factor app</a> 
methodology to aid automation tasks or decisions about a component or collection of compoenents.  

<img src="https://github.com/laetho/metagraf/raw/master/doc/component.png" alt="A graph of a software component.">

# mg

**mg** is a tool that operates on *metaGraf* specifications or a collection of specifications.

**mg** can turn a metaGraf specification into Kubernetes resources, inspect the environment to 
determine if the software components requirements are covered, generate documentation, aid configuration
management for CD or GitOps flows and provide graphs of the environment from a collection of
specifications.

TODO: Update this section with actual examples and a video of usage.

You can use the example collection provided to experiment. It produces output like 
this if the resulting file is used with the **dot** utility from graphviz.

<img src="https://github.com/laetho/metagraf/raw/master/example.png" alt="Example graph">

## Status

Current focus is experimentation with the model and trying out some use cases
that might help communicate what this structure and tooling can accomplish.

This repository contains the WIP specification and a tool called **mg**
that consumes the specification or specifications and turns that into
actionable items or kubernets resources.

The repository will split in two in the future. One for the datastructure specification
and one for the **mg** utility.


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

Kudos goes to my current employer <a href="https://www.norsk-tipping.no">Norsk Tipping AS</a>,
for letting me work on this in the open. 
