# wilson
**wilson** - quickly write tests for APIs using JSON configs

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
