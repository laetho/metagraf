# mg

**mg** is a tool that operates on *metaGraf* specifications or a collection of specifications.

**mg** can turn a metaGraf specification into Kubernetes resources, inspect the environment to 
determine if the software components requirements are covered, generate documentation, aid configuration
management for CD or GitOps flows and provide graphs of the environment from a collection of
specifications.

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

