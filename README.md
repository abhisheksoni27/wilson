# wilson
**wilson** - quickly write tests for APIs using JSON configs

If you want to quickly write tests for your APIs and, then `wilson` is the tool for you.

Defining what your APIs do on different inputs shouldn't be that hard. All it takes is some info about the parameters, or request body, and what to expect in the response (along with data types).

Here's an example `wilson` config for the `https://www.thecolorapi.com/id`

```
[
  {
    "url": "https://www.thecolorapi.com/id",
    "type_of_request": "get",
    "expected_status_code": 200,
    "request_params": {
      "hex": "00ff00"
    },
    "expected_response": {
      "hex": {
        "value": "String",
        "clean": "String"
      }
    }
  }
]

```

# Installation

```
go get github.com/abhisheksoni27/wilson
```

This will install the `wilson` command in your `GOBIN`

# Usage
  wilson [command]

### Available Commands
  `help`        Help about any command
  
  
  `run`         run all tests in config directory

### The run command
run all tests in config directory

**Usage:**

```
wilson run [flags]
```

**Flags:**

  `-c, --config` **string**               _config directory (default is $HOME/.wilson/)_
  
  
  `-m, --max-tests-at-a-time` **int**   _max-tests-at-a-time number of tests to run in parallel (default 4)_
