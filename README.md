# metaGraf

metaGraf is a opinionated specification for describing a software 
component and what its requirements from the runtime environment are. 
The *mg*, the command, turns metaGraf specifications into Kubernetes
resources, supporting CI, CD and GitOps software delivery.

The specification takes inspiration from the <a href="https://12factor.net">twelve-factor app</a> 
methodology.  

## Status

This repository contains the WIP specification and a tool called **[mg](/docs/mg.md)**
that consumes the specification or specifications and turns that into
actionable items or kubernets resources.

The model is maturing but should still be considered a work in progress. It is used heavily
at Norsk Tipping AS for CICD and GitOps based software delivery. I have changed jobs and it's
currently unclear how much effort I will be able to dedicate to this project.
