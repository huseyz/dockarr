# Dockarr

> Dockarr is still pre-release. Use at your own risk.

Dockarr is a companion app for people running the [\*arr](https://wiki.servarr.com/) stack using docker-compose and using [Prowlarr](https://wiki.servarr.com/prowlarr) to keep Indexers in sync.

Dockarr automatically discovers your "\*arr" services and adds them to the discovered Prowlarr instance as applications. It runs in an infite loop to apply the future changes.

It also adds the discovered download clients to discovered \*arr services.

By default, Dockarr uses your docker container name as the hostname, extracts the service apiKey, urlBase, protocol and the port from the `config.xml` file.

It currently supports:

Arr* services:

- Radarr
- Sonarr
- Lidarr
- Readarr

Download clients:

- Transmission


Coming later:

- Other \*arr services: LazyLibrarian, Mylar, Whisparr etc.
- Download clients: SABnzbd etc.

## Usage

1. Add dockarr to your `docker-compose.yaml` and mount docker socket to it.
2. Label your \*arr services with discovery label: `dockarr.discover`
3. Profit?

```
dockarr:
  container_name: dockarr
  image: ghcr.io/huseyz/dockarr:latest
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
```

```
prowlarr:
  container_name: prowlarr
  image: ghcr.io/hotio/prowlarr:latest
  restart: unless-stopped
  ports:
    - 9696:9696
  labels:
    - dockarr.discover

sonarr:
  container_name: sonarr
  image: ghcr.io/hotio/sonarr:latest
  restart: unless-stopped
  ports:
    - 8989:8989
  labels:
    - dockarr.discover
```

Check the example docker-compose file [here](https://github.com/huseyz/dockarr/blob/main/deploy/docker-compose.yaml) for a more complete example.

## Configuration

Dockarr has two types of configuration support. For general config it uses environment variables. For service specific configuration, use labels.

### Environment variables

| Environment Variable | Description                                                                              | Default | Possible Values          |
| -------------------- | ---------------------------------------------------------------------------------------- | ------- | ------------------------ |
| LOG_LEVEL            | Log level for the Dockarr service.                                                       | Info    | Debug, Info, Warn, Error |
| DELETE_BEHAVIOUR     | When Dockarr does not discover an existing *arr application, this tells it what to do | Ignore  | Delete, Disable, Ignore  |
| SYNC_INTERVAL        | The interval Dockarr discovers and syncs your services, in seconds.                      | 60      | Any integer              |

### Labels

You can override certain settings per service using docker labels.

#### *arr Services

You can override the SyncLevel or the complete address of a service.

| Label                      | Description                                                                                                 |
| -------------------------- | ----------------------------------------------------------------------------------------------------------- |
| dockarr.override.address   | Full address of the application, including protocol(http(s)), hostname, port and base url.       |
| dockarr.override.syncLevel | By default Dockarr sets application-indexer sync level to fullSync. Use `addOnly` or `disabled` to override |

#### Transmission

| Label                      | Description                                                                                                 |
| -------------------------- | ----------------------------------------------------------------------------------------------------------- |
| dockarr.override.host | Hostname of the transmision |
| dockarr.override.port | Port of the transmision |
| dockarr.override.ssl | Whether to use SSL for the transmision connection |
| dockarr.override.urlbase | RPC Url base of transmision |