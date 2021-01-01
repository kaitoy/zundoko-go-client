[![License: CC0-1.0](https://licensebuttons.net/l/zero/1.0/80x15.png)](http://creativecommons.org/publicdomain/zero/1.0/)

zundoko-go-client
=================

zundoko-go-client is Zundoko Kiyoshi client written in Go.

# Build

1. Install Git and Go 1.15+.
2. Clone the project.

    ```console
    $ git clone --recursive https://github.com/kaitoy/zundoko-go-client.git
    $ cd zundoko-go-client
    ```

3. Run Go Build.

    ```console
    $ make build
    ```

# Start Zundoko Server
Node 8+ is required.

In the project root directory, execute the following command to start Zundoko Server.

```console
$ make start-server
```

You can stop Zundoko Server by `make stop-server`.

# Start Zundoko Kiyoshi
Execute the built zundoko-client binary to start a Zundoko Kiyoshi.

```console
$ ./bin/zundoko-client
```

zundoko-client interacts with Zundoko Server and exits after making a Kiyoshi.

# Development

## Generate JSON Decoders
Java 8+ is required.

This project uses [Swagger Codegen](https://github.com/swagger-api/swagger-codegen) to generate JSON decoders.

Write a swagger spec in `swagger/swagger.yaml` and run `make model` to generate decoders.

## Unit Tests
This project uses [Ginkgo](https://onsi.github.io/ginkgo/) and [gomock](https://godoc.org/github.com/golang/mock/gomock) for unit tests.

Steps to write unite tests are as follows:

1. Generate a test suite by `make test-suite path/to/pkg`.

    Just one test suite is needed for each package, and usually there is no need to modify generated ones.

2. Generate a test template by `make test-template path/to/file.go` for each Go file you want to test.
3. Generate mocks by `make mock`.

    This command finds interfaces declared in files under pkg directory,
    and generates a mock for each interface into mock directory.

4. Write tests in the generated templates using the generated mocks.

Execute `make test` to run unit tests.

# License
This project is licensed under the Creative Commons license (CC0 1.0).
