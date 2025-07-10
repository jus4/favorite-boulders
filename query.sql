-- name: GetRoutes :many
SELECT * FROM route
ORDER BY title
LIMIT 100;

-- name: GetRoutesByName :many
SELECT 
  r.route_id,
  r.title,
  r.grade,
  c.name as grag_name,
  s.name as sector_name
FROM route as r
JOIN 
  sector as s 
ON 
  s.sector_id = r.sector_id
JOIN
  crag as c 
ON
  c.crag_id = r.crag_id
JOIN areas AS a 
ON 
  c.area_id = a.area_id 
WHERE
  r.title ILIKE $1
AND 
  a.country = 'Finland'
ORDER
  by r.title
LIMIT 50;

-- name: CreateUser :one
INSERT INTO users (email) VALUES ($1)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 
LIMIT 1;

-- name: CreateFavoriteClimbList :one
INSERT INTO user_favourite_lists (name, description, owner)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetFavoriteClimbsList :many
SELECT * from user_favourite_lists
WHERE owner = $1;

-- name: InsertFavuriteClimbsItems :many
INSERT INTO user_favourite_climb_items (route_id, favourite_climb_list)
SELECT unnest($1::int[]), $2::uuid
RETURNING *;

-- name: GetSectors :many
SELECT 
  s.latitude,
  s.longitude,
  s.name,
  s.sector_id,
  c.name as crag_name
FROM sector AS s 
INNER JOIN crag AS c 
  on s.crag_id = c.crag_id 
JOIN areas AS a 
  on c.area_id = a.area_id 
WHERE a.country = 'Finland';

-- name: GetRoutesBySector :many
SELECT 
  r.route_id,
  r.title,
  r.route_type,
  r.grade,
  s.name as sector_name
FROM route as r
JOIN sector as s
  ON s.sector_id = r.sector_id
WHERE
  r.sector_id = $1;

