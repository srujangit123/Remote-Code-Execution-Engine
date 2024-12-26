# Remote Code Execution Engine

This repository contains a remote code execution engine that allows users to submit code in various programming languages, execute it in isolated Docker containers, and retrieve the output. The engine supports both `x86_64` and `arm64` architectures.

## Features

- Supports C++ and Go programming languages.
- Executes code in isolated Docker containers.
- Supports both `x86_64` and `arm64` architectures.
- Cleans up zombie containers to avoid memory leaks.
- Provides a REST API for code submission and execution.

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

## Running the server
To start the server, run the following command:
```sh
./server
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

## TODO
<input disabled="" type="checkbox"> Rate limiter and tokens<br>
<input disabled="" type="checkbox"> Add frontend and OAuth