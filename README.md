<div align="center">
  <h1>
    zep
  </h1>
  <p>
    <strong>Z</strong> generation enviroment variable gotemplate processor
  </p>
  <p>
    <a href="https://github.com/aasaam/zep/actions/workflows/build.yml" target="_blank"><img src="https://github.com/aasaam/zep/actions/workflows/build.yml/badge.svg" alt="build" /></a>
    <a href="https://codecov.io/gh/aasaam/zep" target="_blank"><img src="https://codecov.io/gh/aasaam/zep/branch/main/graph/badge.svg" alt="Coverage" /></a>
    <a href="https://hub.docker.com/r/aasaam/zep" target="_blank"><img src="https://img.shields.io/docker/image-size/aasaam/zep?label=docker%20image" alt="docker" /></a>
    <a href="https://github.com/aasaam/zep/blob/master/LICENSE"><img alt="License" src="https://img.shields.io/github/license/aasaam/zep"></a>
  </p>
</div>

## Usage

Include zep as a small binary in your Docker image using multi-stage builds:

```Dockerfile
# build layer
FROM ghcr.io/aasaam/zep:latest as zep

# final layer
FROM yourimage

COPY --from=zep /usr/bin/zep /usr/local/bin/zep
```

Then in your entrypoint script:

```sh
#!/bin/sh
set -e
# Use zep to process your template file
/usr/local/bin/zep /path/to/your/template > final-config.ext
# Run your application
/run/my/awesome-process
```

<div>
  <p align="center">
    <a href="https://aasaam.com" title="aasaam software development group">
      <img alt="aasaam software development group" width="64" src="https://raw.githubusercontent.com/aasaam/information/master/logo/aasaam.svg">
    </a>
    <br />
    aasaam software development group
  </p>
</div>
