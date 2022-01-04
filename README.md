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

## Licenses
Tagwatch is licensed under [GPLv3+](./COPYING).

Tagwatch uses [go-yaml/yaml](https://github.com/go-yaml/yaml/tree/v2).
