# Image optimizer Traefik middleware plugin

Image optimizer middleware is a [Traefik](https://traefik.io) plugin designed to optimize image responses on the fly.
It can remove unwanted metadata, convert to the desired format (like webp) and resize images.

Easy to extend, you can implement your own image processing logic and caching systems, please refer to the [documentation](https://doc.traefik.io/traefik/plugins/).

## Usage

To be active for a given Traefik instance, it must be declared in the static configuration.

You can request your images as usual, and just add `w=<width>` in query param.
```bash 
curl "http://demo.localhost/very_big.jpg" # return converted to webp and without metadata 
curl "http://demo.localhost/very_big.jpg?w=1725" # return resized with 1725px width, converted to webp and without metadata
```

### Configuration

For each plugin, the Traefik static configuration must define the module name (as is usual for Go packages).

The following declaration (given here in YAML) defines an plugin:

```yaml
# Static configuration
pilot:
  token: xxxxx

experimental:
  plugins:
    image_optimizer:
      moduleName: github.com/agravelot/image_optimizer
      version: v0.1.0
```

Here is an example of a file provider dynamic configuration (given here in YAML), where the interesting part is the `http.middlewares` section:

```yaml
# Dynamic configuration

http:
  routers:
    my-router:
      rule: host(`demo.localhost`)
      service: service-foo
      entryPoints:
        - web
      middlewares:
        - image_optimizer

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
  
  middlewares:
    image_optimizer:
      plugin:
        config:
          processor: <processor>
          imaginary:
            url: http://imaginary:9000
          cache: <cache>
          file:
            path: /tmp
          redis:
            url: redis://<user>:<pass>@localhost:6379/<db>
```

List of available processors:

| Name         | Note                         |
| -------------|:---------------------------:|
| imaginary    | Use [imaginary](https://github.com/h2non/imaginary) as processor to manipulate images, can be easily scaled. (recommended)     |
| local        | Process images in Traefik itself, ⚠️ currently **not implemented** cause of interpreter limitations. |
| none         | Keep images untouched (default)    |

List of available caches:

| Name         | Note                         |
| -------------|:---------------------------:|
| file         | Save images in given directory. (recommended)     |
| redis        | Save images in redis, work best in HA environments.  ⚠️ currently **not implemented** cause of interpreter limitations. |
| memory       | Keep images directly in memory, only recommended in development. ⚠️ Cache invalidity not implemented yet.    |
| none         | Do not cache images (default)    |

### Dev Mode

An easy to bootstrap development environment using docker is available in the `demo/` folder.

Make sure to copy `demo/.env.example` to `demo/.env` and provide the required Traefik pilot token.
To quickly start the demo, please run at the root of this repository:

```bash
make demo
```

Services are now accessible with these endpoints.

| Service           | URL                         |
| ----------------- |:---------------------------:|
| Demo frontend     | http://demo.localhost       |
| Traefik dashboard | http://localhost:8080       |
| Grafana           | http://grafana.localhost    |
| Prometheus        | http://prometheus.localhost |

You can now implement our own image processing or caching systems by implementing `processor.Processor` and `cache.Cache` interfaces.
After that, make sure to add your new configuration and add it to the corresponding factory.

Note that only one plugin can be tested in dev mode at a time, and when using dev mode, Traefik will shut down after 30 minutes.

