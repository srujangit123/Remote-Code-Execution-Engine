# Remote Code Execution Engine

![Go](https://github.com/srujangit123/Remote-Code-Execution-Engine/actions/workflows/go.yml/badge.svg)


This repository contains a remote code execution engine that allows users to submit code in various programming languages, execute it in isolated Docker containers, and retrieve the output. The engine supports both `x86_64` and `arm64` architectures.

## Features

- Supports C++ and Go programming languages.
- Executes code in isolated Docker containers.
- Supports both `x86_64` and `arm64` architecture machines.
- Cleans up zombie containers to avoid memory leaks.
- Provides a REST API for code submission and execution.
- Restricts the usage of system resources (Memory, CPU, max processes, max files, max file size)
- Kills a container if it is taking more than a minute to complete the execution.
- Supports custom docker images and compilation commands for each programming language.

## Prerequisites

- Docker
- Go

## Setup

1. Clone the repository:

    ```sh
    git clone https://github.com/srujangit123/Remote-Code-Execution-Engine.git
    cd Remote-Code-Execution-Engine
    ```

2. Build the Docker images:
This builds the default docker images and uses the default config (`config.yml`)
Skip this step if you want to use custom docker images and commands. Make sure you modify the `config.yml` if you are using custom docker images and commands

    ```sh
    bash scripts/build_docker.sh
    ```

3. Build the server:

    ```sh
    go build -o server cmd/server.go
    ```

## Configuration

The configuration file [config.yml](http://_vscodecontentref_/0) specifies the settings for each supported language. Here is an example configuration:

```yaml
cpp:
  extension: ".cpp"
  image: "cpp_arm64:latest"
  command: "/usr/bin/run-code.sh {{LANGUAGE}} {{FILE}} {{INPUT}}"
golang:
  extension: ".go"
  image: "golang_arm64:latest"
  command: "/usr/bin/run-code.sh {{LANGUAGE}} {{FILE}} {{INPUT}}"
```

### Variables available in the command config
- {{LANGUAGE}} - programming language
- {{FILE}} - Code file with the extension as specified in the config
- {{INPUT}} - Input file if the input is provided by the user

These variables are replaced with appropriate values before creating the code container.

## Running the server
To start the server, run the following command:
```sh
./server
```

### Flags
- `--code-dir`
    The default directory where the code files will be stored is `/tmp/`
    A separate directory is created for every language.
    To change the directory where the code files will be stored(eventually removed by the garbage collector), use
```sh
./server --code-dir /path/to/code/dir
```

- `--resource-constraints`
    By default, resource constraints are turned off to improve the performance, if you want to enable it, use
```sh
./server --resource-constraints true
```

## API

### Submit Code
- URL: `/api/v1/submit`
- Method: `POST`
- Content-Type: `application/json`
- Request Body:

```json
{
    "code": "base64_encoded_code",
    "input": "base64_encoded_input",
    "language": "cpp"
}
```
- Response:
```json
{
    "output": "execution_output"
}
```

### Example
To submit a code execution request, you can use the following `curl` command:
```sh
curl -X POST http://localhost:9000/api/v1/submit \
     -H "Content-Type: application/json" \
     -d '{"code": "base64_encoded_code", "input": "base64_encoded_input", "language": "cpp"}'
```