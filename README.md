<p align="center">
  <img src="./.github/images/axolotl-cloud.png" alt="Logo" width="200"/>
</p>

<h1 align="center">Axolotl Cloud</h1>

**Axolotl Cloud** is a lightweight self-hosted platform to manage and run your Docker-based projects.
Easily create, import, and configure containers per project — all from a clean web UI.


Create a `.env` file in the backend directory with the following content:

```sh
HTTP_PORT=8080
VOLUMES_PATH="/path/to/volumes" # Absolute path to the volumes directory
DATABASE_PATH=./data.db
```

```sh
mkdir -p /path/to/volumes
```

