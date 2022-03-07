# internalizer

### What is this?

`internalizer` is a tool for figuring out which folders on a Go application can be made [internal][1]. 

It works by building an import graph of your project's packages and suggesting how to move things around to keep those imports working.

It only makes sense for Go *applications* and not on Go *libraries*, since `internalizer` can't figure out which symbols are used by external dependencies.

This is a work in progress, which I am trying to use for the [datadog-agent][2].

[1]: https://docs.google.com/document/d/1e8kOo3r51b2BWtTs_1uADIA5djfXhPT36s6eHVRIvaU/edit
[2]: https://github.com/DataDog/datadog-agent
