<p align="center">
  <img src="./.github/images/axolotl-cloud.png" alt="Logo" width="200"/>
</p>

<h1 align="center">Axolotl Cloud</h1>

**Axolotl Cloud** is a lightweight self-hosted platform to manage and run your Docker-based projects.
Easily create, import, and configure containers per project — all from a clean web UI.


```sh
mkdir -p ./volumes
mkdir -p ./data
```

```yml
services:
  app:
    image: ghcr.io/bastien2203/axolotl-cloud:latest
    ports:
      - "8080:8080"
    environment:
      HTTP_PORT: 8080
      ENV: production
      VOLUMES_PATH: /app/volumes
      GIN_MODE: release
      DATABASE_PATH: /app/data/data.db
    volumes:
      - ./volumes:/app/volumes
      - ./data:/app/data
      - /var/run/docker.sock:/var/run/docker.sock

```
