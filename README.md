# Dashborg Examples

See all the examples running here: https://demo.dashborg.net

This repository contains example code for Dashborg.  Go code is contained in the "golang" directory and HTML for the apps is located in the "panels" directory.  All the examples should be run from the *root* of the repository as they access the panels via a relative path (e.g. ```panels/todo.html```).  So here's how to run the "todo" example:

```
go run ./golang/todo/
```

## Golang Setup

Clone this repository
```
git clone https://github.com/sawka/dashborg-examples.git
cd dashborg-examples
```

Run an example:
```
go run ./golang/todo/
```

Note that this format, with the initial ```./```, allows for our *main* package to be spread over multiple files in the example directory.

## Python Setup

The python examples are currently being updated to work with the new Dashborg service.
