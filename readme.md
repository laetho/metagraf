# metaGraf

**metaGraf** is a specification describing the necessary metadata and information about a software component for building it, running it or both. It's intended to be used as a source of truth for CI/CD or GitOps scenarios.
 
It is inspired by the <a href="https://12factor.net">twelve-factor app</a> 
methodology to aid automation tasks or decisions about a component or collection of compoenents.  

<img src="https://github.com/laetho/metagraf/raw/master/docs/component.png" alt="A graph of a software component.">

# Status

The model is maturing but should still be considered a work in progress. It is getting quite heavy
usage internally in our CI processes. It is also getting internal usage in a GitOps setting for the 
CD parts. Better public examples are forthcoming.

This repository contains the WIP specification and a tool called **[mg](/docs/mg.md)**
that consumes the specification or specifications and turns that into
actionable items or kubernets resources.

The repository will split in two in the future. One for the datastructure specification
and one for the **[mg](/docs/mg.md)** utility.


# Acknowledgement

Kudos goes to my current employer <a href="https://www.norsk-tipping.no">Norsk Tipping AS</a>,
for letting me work on this in the open. 
