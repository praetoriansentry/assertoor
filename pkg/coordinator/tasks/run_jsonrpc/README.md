## `run_jsonrpc` Task

### Description
The `run_jsonrpc` task makes a single JSON RPC request against an endpoint.

### Configuration Parameters

- **`clientPattern`**:\
  A regular expression pattern used to specify which clients to check. This allows for targeted health checks of specific clients or groups of clients within the network. A blank pattern targets all clients.
- **`method`**:\
  The JSON RPC method that will be called during the test

### Defaults

Default settings for the `run_jsonrpc` task:

```yaml
- name: run_jsonrpc
  config:
    clientPattern: ""
    method: ""
    params: []
```
