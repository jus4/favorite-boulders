# Suomi topot
Personal project to make all climbing areas accessible to everyone. 

## Stack
- Golang Echo for backend
- Postgres for database
- HTMX + tailwindcss + templ + typescript for frontend

## Run
Make sure you have postgres running and defined in .env

- `make live` for developing
- `make live/live/esbuild`  compile javascript
- `make live/live/tailwind`  compile tailwindcss
