# API blueprint
A template is for building a HTTPServer purpose codebase. It is influenced Clean Architecture and featuring Simplicity, Flexibility in mind

## Getting started
- Install [make](https://www.geeksforgeeks.org/how-to-install-make-on-ubuntu/) utility command
- Install [go](https://go.dev/doc/install) compiler with the latest release
- Install [air](https://github.com/air-verse/air) for hot reloading
- Install [python](https://www.python.org/downloads/) compiler with the latest release
- Replace your project name with the default name
- Install virtual env (assume python installed, skip if already installed)
```shell
python3 -m venv venv
```
- Activate virtual env 

For Windows:
```shell
./venv/bin/Activate
```
For Linux/MacOs:
```shell
source ./venv/bin/activate
```
- Install pre-commit
```shell
pip install pre-commit
pre-commit install
```
- Change ```config.yaml``` to an appropriate configuration.
- Find **TODO** comments in the source code and follow guidance

## What's next 
For further customization, please learn the source code as it's designed with Simple mindset at first.
- Install [protoc](https://grpc.io/docs/protoc-installation/) compiler if you need grpc
- Install [swag](https://github.com/swaggo/swag) command if you need Restful API Documentation served by Swagger Spec
- Install [docker](https://docs.docker.com/desktop/install/mac-install/) command if you need to build an image ([Dockerfile](https://github.com/hyuti/api-blueprint/blob/main/Dockerfile) also provided out of the box)

## Contributions
Please fire pull requests if you get any bright ideas.