<p align="center">
  <img src="https://user-images.githubusercontent.com/7422050/56345043-1924ca00-61bf-11e9-9832-58a50379851f.png" width="200" alt="Looper"/>
</p>

Looper is a test runner for functional testing of dockerized microservices. It records and replays interactions between the service under test and its dependant services, drastically reducing the seutp complexity as well as test execution time and eliminating flakiness of functional tests.

## Usage [![Build Status](https://travis-ci.com/fawind/looper.svg?branch=master)](https://travis-ci.com/fawind/looper)

```
>> looper --help

Usage:
  app [command]

Available Commands:
  help        Help about any command
  record      Run in record mode
  replay      Run in replay mode

Flags:
      --compose string     Default docker-compose file for the services (default "./docker-compose.yml")
      --directory string   Directory to store the test dumps in (default "replay-dumps")
  -h, --help               help for app
      --out string         File name for the mitm output file (required when run in stand-alone proxy mode)
      --port int           Port to use for the MITM proxy (default 9999)
      --proxyOnly          Only start proxy in stand-alone mode
      --service string     Docker setService name to test (required when not run in stand-alone proxy mode)
      --sleep int          Time to wait after starting docker services in ms
      --test string        Test command to execute (required when not run in stand-alone proxy mode)
```

### Install or build the project:

```bash
# Build
go build
# Install using go get
go get github.com/fawind/looper
```

### Test docker-compose services

1. Run service in record mode:
```bash
looper record \
    --service my-service \
    --compose ./path/to/docker-compose.yml \
    --test '<TESTCMD>'
```

2. Run service in replay mode:
```bash
looper replay \
    --service my-service \
    --compose ./path/to/docker-compose.yml \
    --test '<TESTCMD>'
```

### Example Application

Browse the [notes-service](https://github.com/fawind/looper/tree/master/examples/notes-service) app for an example app with predefined commands for [record](https://github.com/fawind/looper/blob/master/examples/notes-service/record.sh) and [replay](https://github.com/fawind/looper/blob/master/examples/notes-service/replay.sh) actions.


## Functional Testing for Dockerized Microservices

Nowadays, modern applications exhibit a numerous amount of different services.
All these services have their own unit and integration tests covering the core functionality.
Unfortunately, these tests are not enough to ensure the stability and robustness of the whole system.
Additional functional tests are required to verify the interaction between the dependent services.  
Usually, when functional testing a given service, all the dependant services need to be up and running.
While this often simple for smaller applications, scaling such a test setup to a large system with many intertwined services is a difficult problem.  
Usually, the setup for functional tests comes with a large complexity overhead, which makes it hard to maintain and takes a significant amount of resources and time for starting the systems and executing the tests.
Furthermore, functional tests are often slow because of many network requests getting exchanged between services and flaky due to a lack of determinism and unavailability of external services.  
As a result, many developers struggle with integrating functional tests in their development workflow which reduces their overall confidence in their system and increases the likelihood that bugs are not caught before shipping to production.

To tackle these issues, we introduce a test runner for functional testing of dockerized microservices.
The test runner reduces the overhead for setting up and running functional tests, drastically reduces their execution time, and eliminates flakiness through the recording and replay of requests.
Finally, we make it dead simple for new developers to get up and running by eliminating the need to setup a complex test environment when doing code changes.

### Test Runner

To address these problems we created a test runner for functional testing dockerized services.
Our tool provides the ability to simulate the dependent services.
Hence, the functional tests can now run without network access and with the setup complexity of the single service.
This is done by capturing and replaying the network traffic between the system under test (SUT) and its dependent services.

![System under test](https://user-images.githubusercontent.com/7422050/51401182-74a4d480-1b4a-11e9-80ba-247de6c3859f.png)

Initially, tests have to be run in record mode, which brings up the SUT and its dependent services.
During record mode, the tests are executed against the full stack of services and requests from the SUT to its dependent services are recorded.
Afterwards, tests can be run in replay mode.
Here, only the SUT is started.
Requests from the SUT to its dependent services are answered using the recorded replies.
This leads to faster and reliable tests during development and as part of continuous integration pipelines which increases developer productivity and system stability.

### Usage with Docker Services

To run a test suite in record mode, the developer has to provide the path to the docker-compose file and the docker service name of the SUT, as well as the command for executing the tests suite.

```bash
looper record \
    --service frontend-service \
    --compose ./docker-compose.yml \
    --test 'npm test'
```

This brings up the SUT and its dependent services and executes the test suite.
Requests from the SUT are recorded and stored in a file.

For further runs, tests can be run in replay mode.
Similar to record mode, the developer provides the docker-compose file, docker service name, and command of the test suite.

```bash
looper replay \
    --service frontend-service \
    --compose ./docker-compose.yml \
    --test 'npm test'
```

When running in replay mode, only the SUT is started and requests are answered based on the recorded messages.
The file containing the request dumps is automatically selected based on the used service and test command.
This allows to seamlessly work with multiple test suites.

The test runner is written as a command-line app in Go.
It registers a man-in-the-middle proxy for the provided docker service which intercepts all http traffic.
In record mode, the proxy records incoming requests and passes them through.
In replay mode, all requests are intercepted and answered using the provided dump of previous requests.

### Usage in Proxy-Only Mode

Our test runner works best with dockerized microservices.
However, it can also be used to test systems that can not easily containerized such as iOS apps.
For this, the test runner can be run in a proxy-only mode.

```bash
looper [record|replay] \
    --proxyOnly \
    --out request-dump.mitmdump
```

In proxy-only mode, the test runner still allows to record and replay the requests of the SUT.
However, users are responsible for registering the proxy with the SUT and starting the required services.

### Developer Workflow

Initially, tests have to be run in record mode.
After the tests finish, a file will be automatically created containing the recorded request dumps.

```
./replay-dumps/my-service-302ffec00f89dcb36d50bd7cb59f1c77d70d5ce3.mitmdump
```

The dump files should be checked into source control in order to make them accessible to other developers.
When onboarding new developers that want to work on the service, they just have to clone the repository and can get started working immediately.
After making their code changes, the developer can verify the integrity of the service by running the tests against the requests dumps without having to set up the dependant services.
As long as the developer does not introduce breaking changes to APIs or adds new tests, the test runner can guarantee test integrity without needing to start any dependant services.
When introducing breaking changes to APIs or adding new tests, the test runner will fail.
In this case, tests need to be run again in record mode and the updated request dumps have to be committed together with the changed code.

Furthermore, the test runner can be used as part of your continuous integration pipeline. Here, it has the potential to drastically speed up CI checks and reduce the time that developers spend waiting for CI to finish.

To sum it up, our test runner reduces the setup complexity of tests, speeds up their execution time, and gives an end to flaky functional tests.
We hope that it enables you to ship faster and with more confidence!
