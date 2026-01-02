# Tagwatch

Watches container registries for new and changed tags and creates an RSS feed for detected changes.

## Configuration
Tagwatch is configured through the [`./tagwatch.yml`](./tagwatch.example.yml) file which is read from the current working directory.
Alternatively, the config can be read from the path configured in the `TAGWATCH_CONF` environment variable.

## Installation
### From Docker Hub
Checks configured tags every 4 hours and serves the generated feed on http://container.ip:8080/feed.xml
```bash
docker run -v $PWD/tagwatch.yml:/tagwatch.yml woefe/tagwatch:latest
```

### Manually
```bash
# Clone and build tagwatch
git clone https://github.com/woefe/tagwatch
cd tagwatch
go build

# Show help
./tagwatch help

# Print the generated feed for the example config once
TAGWATCH_CONF=tagwatch.example.yml ./tagwatch run
```

## Limitations
Every checked tag counts as a pull according to the rate limitation mechanism employed at the Docker Hub.
Hence, be careful which tag patterns you watch with tagwatch.
A simple `.*` pattern will instantly consume all pulls of the free tier in many cases!

As of now, tagwatch is only tested against Docker Hub and the Google Container Registry.

## Licenses
Tagwatch is licensed under [GPLv3+](./COPYING).

Tagwatch uses [yaml/go-yaml](https://github.com/yaml/go-yaml).
