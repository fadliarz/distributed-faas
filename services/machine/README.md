# Machine Service

A gRPC service that handles function execution requests from the dispatcher service.

## Features

- Receives `ExecuteFunctionRequest` via gRPC
- Logs the request payload (InvocationID, FunctionID, SourceCodeURL)
- Returns a success response

## Usage

### Build and Run

```bash
# Build the service
make build

# Run the service
make run
```

The service will start listening on port `50051`.

### Manual Commands

```bash
# Build manually
go build -o bin/machine-service .

# Run manually
./bin/machine-service
```

## API

The service implements the `MachineService` gRPC interface:

```protobuf
service MachineService {
  rpc ExecuteFunction(ExecuteFunctionRequest) returns (ExecuteFunctionResponse);
}
```

### Request Format

```protobuf
message ExecuteFunctionRequest {
  string invocation_id = 1;
  string function_id = 2;
  string source_code_url = 3;
}
```

### Response Format

```protobuf
message ExecuteFunctionResponse {
  string invocation_id = 1;
  string status = 2;
  string message = 3;
  string result = 4;
}
```
