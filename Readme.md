# Suomi topot
Personal project to make all climbing areas accessible to everyone. 

## Stack
- Golang Echo for backend
- Postgres for database
- HTMX + tailwindcss + templ + typescript for frontend

## Running development env
### Setting up the env
1. Add your API key to `envs/app.env`, `MAP_API_KEY="YOUR_API_KEY"`.
2. `topo-repo$ mkdir db/dump`
3. `cp YOUR_DUMP.sql db/dump/dump.sql`

### To Run
```shell
topo-repo$ docker compose build
topo-repo$ docker compose up
```
### To Stop
```shell
topo-repo$ docker compose down
```
### Clean database
Script removes docker volume containing db data. New data will be initialized next time starting the db service.
```shell
topo-repo$ ./scripts/clean_db.sh
```

## Old notes
Make sure you have postgres running and defined in .env

- `make live` for developing
- `make live/live/esbuild`  compile javascript
- `make live/live/tailwind`  compile tailwindcss
