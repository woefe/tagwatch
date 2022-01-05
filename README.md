# Tagwatch

Watches container registries for new and changed tags and creates an RSS feed for detected changes.

## Configuration
Tagwatch is configured through the [`./tagwatch.yml`](./tagwatch.example.yml) file which is read from the current working directory.
Alternatively, the config can be read from the path configured in the `TAGWATCH_CONF` environment variable.

## Installation
### From Docker Hub
```bash
docker run -v $PWD/tagwatch.yml:/tagwatch.yml woefe/tagwatch:latest
```

### Manually
```bash
git clone https://github.com/woefe/tagwatch
cd tagwatch
go build
```

## Limitations
Every checked tag counts as a pull according to the rate limitation mechanism employed at the Docker Hub.
Hence, be careful which tag patterns you watch with tagwatch.
A simple `.*` pattern will instantly consume all pulls of the free tier in many cases!

As of now, tagwatch is only tested with the Docker Hub registry.

## Licenses
Tagwatch is licensed under [GPLv3+](./COPYING).

Tagwatch uses [go-yaml/yaml](https://github.com/go-yaml/yaml/tree/v2).
