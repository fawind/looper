<p align="center">
  <img src="https://user-images.githubusercontent.com/7422050/56345043-1924ca00-61bf-11e9-9832-58a50379851f.png" width="200" alt="Looper"/>
</p>

## Usage

```
>> docker-test --help
Usage:
  app [command]

Available Commands:
  help        Help about any command
  record      Run in record mode
  replay      Run in replay mode

Flags:
      --compose string   docker-compose file for the services (default "./docker-compose.yml")
  -h, --help             help for app
      --out string       File name for the mitm output file (default "out.mitmdump")
      --port int         Port to use for the MITM proxy (default 9999)
      --service string   Docker service name to test (required)
      --sleep int        Time to wait after starting docker services in ms (optional)
      --test string      Test command to execute (required)
```

1. Build the project:
```bash
cd cli && go build
```

2. Run service in record mode:
```bash
./docker-test record \
    --service my-service \
    --compose ./path/to/docker-compose.yml \
    --test '<TESTCMD>'
```

3. Run service in replay mode:
```bash
./docker-test replay \
    --service my-service \
    --compose ./path/to/docker-compose.yml \
    --test '<TESTCMD>'
```

### Example

Browse the [notes-service](https://github.com/fawind/docker-test/tree/master/examples/notes-service) app for an example app with predefined commands for [record](https://github.com/fawind/docker-test/blob/master/examples/notes-service/record.sh) and [replay](https://github.com/fawind/docker-test/blob/master/examples/notes-service/replay.sh) actions.

## Functional Testing for Dockerized Microservices

Nowadays, modern application systems exhibit a numerous amount of different services. All these services have their own unit and integration tests covering the core functionality. Unfortunately, these tests are not enough to ensure the stability and robustness of the whole system. Additional functional tests are required that verify interaction between the dependent services.

Usually when testing a given service, all the dependent service need to be up and running which results in a significant complexity overhead due to the different setups and resource requirements. Specifically, if the dependent services haven’t changed it is unnecessary to set up the complete system again. Another challenge for functional testing is network access which results in slow tests and might lead to flaky behavior due to a lack of determinism or unavailable external services.

### Test Runner
To address these problems we created a test runner for functional testing dockerized services. Our tool provides the ability to simulate the dependent services hence the functional tests can now run without network access and with the setup complexity of the single service. This is done by capturing and replaying the network traffic between the system under test (SUT) and its dependent services.

![System under test](https://user-images.githubusercontent.com/7422050/51401182-74a4d480-1b4a-11e9-80ba-247de6c3859f.png)

Initially, tests have to be run in record mode which brings up the SUT and its dependent services. During record mode, the tests are executed against the full stack of services and requests from the SUT to its dependent services are getting recorded.
Afterwards, tests can be run in replay mode. Here, only the SUT is started. Requests from the SUT to its dependent services are answered using the recorded replies. This leads to faster and more reliable tests during development and as part of continuous integration pipelines which increases developer productivity and system stability.


### Usage
To run a test suite in record mode, the developer has to provide the path to the docker-compose file and the docker service name of the SUT as well as the command for executing the tests suite.

```bash
./docker-test record \
    --service frontend-service \
    --compose ./docker-compose.yml \
    --test 'npm test'
```

This brings up the SUT and its dependent services and executes the test suite. Requests from the SUT are recorded and stored in a file. For further runs, tests can be run in replay mode. Similar to record mode, the developer provides the docker-compose file, docker service name, and command of the test suite to use.

```bash
./docker-test replay \
    --service frontend-service \
    --compose ./docker-compose.yml \
    --test 'npm test'
```

When running in replay mode, only the SUT is started and requests are answered based on the recorded messages.

The test runner is written as a command line app in Go. It registers a man-in-the-middle proxy for the provided docker service which intercepts all http traffic. In record mode, the proxy simply dumps incoming requests and passes them through. In replay mode, all requests are intercepted and answered using the provided dump of previous requests.

### Future Work
The current version of the project allows the execution of tests against a docker stack in record and replay mode. The next steps are to integrate our tool with a production system to evaluate our setup and workflows. We also plan to work on CI integration, making it easy to use our tool as part of a CI pipeline and to automatically detect when recorded request dumps have to be discarded due to changed behavior of dependent services.
