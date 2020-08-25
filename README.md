# data-receiver

[![Build Status](https://travis-ci.org/efark/data-receiver.svg?branch=master)](https://travis-ci.org/efark/data-receiver)
[![Go Report Card](https://goreportcard.com/badge/github.com/efark/data-receiver)](https://goreportcard.com/report/github.com/efark/data-receiver)

This is a small webserver project to receive http requests and write the data somewhere.
I've made interfaces for the most relevant actions in the data flow:
- parsing and loading configuration.
- extracting parameters from the request.
- authenticating a request.
- writing the data.

I tried to keep it small and simple, so it works on a one service-one client basis.
It doesn't cover every use case, but with the examples in the code it should be easier to implement new classes and functions to cover other use cases.

Things that should be modified to make it production ready. (Besides adapting everything to meet your needs)
- Metrics.
- Configuration.
- Writer.

In a similar fashion, you can add an interface to validate the content of the request you're going to write.

Also, you can add another interface to create some more complex messages, in which case you would have to modify the writers to accept this new format.

Some thoughts and words on decisions made while coding this project:

Authenticator was made for signature authentication with some shared key.
Other kinds of authentication can be also made and applied, but they probably require some extra work and ended up being out of scope. For example, some things that could be applied here LDAP authentication, token auth, basic auth (Gin has it already out of the box).

You may notice that some interfaces are implemented by pointers and others by structs. In few words, most times using a pointer is the way to go and having methods receiving a struct is the exception.
One such case is when no method modify anything in the struct. Signer implements Authenticator and has only fixed values in the structs inner fields (key, and the functions to generate the hash and encode it), and it only calculates a hash and returns an error message. This kind of calculation can be implemented by a struct, and it's not a big struct so passing it by value shouldn't generate much overhead.

I isolated configuration in one package to have the logic together and write some unit tests.
I would recommend Json only for simple structures, if the structure gets bigger and more complex then human-readable formats are better (like yaml).
Both json and yaml are implemented here. JsonFile expects a JsonLines file and reads it line by line, but yaml expects a file containing several service configurations nested.

Writer interface is really basic.
This is something to work on, you'll probably want to have some metadata, like saving the timestamp when the request was received or some data regarding the service.
Also, if you're writing to a db or Kafka, or something else, then the struct you need will be a bit bigger.

I tried to make as many tests as possible.
Again, it doesn't cover everything. But I hope the examples help to build the things that are missing. 

main.go seems a bit bulky, and I wanted to make it simpler, but didn't have the time yet to improve it.
Some things that may be important in main.go (were to me at some point): having an easy way to increase timeout as the default values were too short and having a graceful shutdown for the app.

I also tried not to import too many packages, trying to make this really small. You can check some other things below that may be useful to work on this, or for some other projects.

I included the dockerfile to build the app.

## Tech and docs
This project was developed with Go v1.14.2.

[The Go Programming Language Specification](https://golang.org/ref/spec): This one is basic, and I recommend to go back to it every once in a while. Programming is all about learning and practicing, and many things take time to understand. If you get stuck, maybe you need to read it again after some rest.

[Effective Go](https://golang.org/doc/effective_go.html): This one is really important to write code in a Go-ish style.

Gin-Gonic is the framework I've chosen to write this app. You can find Gin's official documentation [here](https://github.com/gin-gonic/gin).

For metrics and monitoring, Prometheus may be the way to go. [Docs](https://github.com/prometheus/client_golang)

Json. [Official Documentation](https://www.json.org/json-en.html). This is for Go: ["encoding/json"](https://golang.org/pkg/encoding/json/). There are also other packages to work with Json, but as it's only for configuration, the standard library is fast and simple enough.

Yaml. [Official Documentation](https://yaml.org/). For Go: [go-yaml](https://github.com/go-yaml/yaml)

To manage configuration, you may want to check this project: [Viper](https://github.com/spf13/viper).

To develop a more powerful Cli, this may be useful: [Cobra](https://github.com/spf13/cobra).

## How to build it and run it.
Clone the repo as usual.

Create the executable:
`go build .`

Run:
`./data-receiver`

Configuration can be added with `-inline-config` (and passing the whole json) or `-config` + filepath flags.

Run the tests, from the app's root:

`go test ./...`

And the docker image can be built using the Dockerfile: `docker build -t my-data-receiver:latest .`

Then you can start the docker image with:
```
docker run \
-p 8080:8080 \
my-data-receiver:latest \
/app/data-receiver \
-inline-config '{"services": {"test_service": {"extractor": {"type": ""}, "authenticator": {"type": ""}, "writer": {"type": "ConsoleWriter"}}}}'
```

This would work to test that the webserver receives messages:
`curl -X POST localhost:8080/data/test_service -d '{"hello": "world"}' -i`

## Or just take what you need and be on your way

Just import the packages you need and use them in your application.
You know, like this:

```
package myPackage

import (
	"github.com/efark/data-receiver/authenticator"
	"github.com/efark/data-receiver/configuration"
	"github.com/efark/data-receiver/extractor"
	"github.com/efark/data-receiver/logger"
	"github.com/efark/data-receiver/writer"
)
```

[MIT license](LICENSE.md).