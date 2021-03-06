# somux

somux is simple port forwarder via stdio.

## Usage

```
Usage: somux [-L value]... [-R value]... [command args...]

  -L value
        local to remote forwarding (e.g. 8080:127.0.0.1:80)
  -R value
        remote to local forwarding (e.g. 9000:127.0.0.1:9000)
  -v    verbose move
```

## Install

```sh
go get github.com/ngyuki/somux
```

And

```sh
docker pull ngyuki/somux
```

## Example

```sh
# Specify remote docker host
export DOCKER_HOST=example.com

# Run nginx container
docker run --name=nginx --rm -d nginx:alpine

# Run local somux and somux container
somux -L 8080:nginx:80 docker run --name=somux --rm -i --link=nginx ngyuki/somux &

# Forwarded from local 8080 port to nginx container 80 port
curl http://localhost:8080/
```

See [example/](./example/) directory for more example.

## Licence

[MIT License](https://opensource.org/licenses/mit-license.php)
