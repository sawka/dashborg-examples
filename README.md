# Dashborg Examples

This repository contains example code for Dashborg in both Golang and Python.  The repository is structured with 3 directories: golang, python, and panels.  The panels directory is shared between both the Go code and the Python code as the HTML templates are identical between the different backend SDKs.  Because of this, all the examples should be run from the *root* of the repository as they access the panels via a relative path (e.g. ```panels/todo.html```).  So here's how to run the "todo" example for both python and go:

```
go run ./golang/todo
python python/todo/todo.py
```

## Python Setup

Clone this repository
```
git clone https://github.com/sawka/dashborg-examples.git
cd dashborg-examples
```

Set up virtual environment
```
python -m venv env
source env/bin/activate
```

Download requirements
```
python -m pip install -r requirements.txt
```

Run an example!
```
python python/todo/todo.py
```

## Golang Setup

Clone this repository
```
git clone https://github.com/sawka/dashborg-examples.git
cd dashborg-examples
```

Run an example:
```
go run ./golang/todo
```

Note that this format, with the initial ```./```, allows for our *main* package to be spread over multiple files in the example directory.

