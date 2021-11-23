# GRPC Bi Drectional Stream

Simple project to validate bi directional streams are working across AWS LB's (NLB/ALB)
## Building
Edit the Makefile to match your specific requirements or pass environment variables

**Example:**
`REGISTRY=foobar PKG=github.com/foobar/grpc-stream make build`

Run `make` or `make build` to compile your app.  This will use a Docker image
to build your app, with the current directory volume-mounted into place.  This
will store incremental state for the fastest possible build.  Run `make
all-build` to build for all architectures.

Run `make container` to build the container image.  It will calculate the image
tag based on the most recent git tag, and whether the repo is "dirty" since
that tag (see `make version`).  Run `make all-container` to build containers
for all supported architectures.

Run `make push` to push the container image to `REGISTRY`.  Run `make all-push`
to push the container images for all architectures.

Run `make clean` to clean up.

Run `make help` to get a list of available targets.

## Docker
A docker compose file is provided that can be tweaked to get this working in docker 



