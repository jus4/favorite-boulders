-- name: GetRoutes :many
SELECT * FROM route
ORDER BY title
LIMIT 100;

-- name: GetRoutesByName :many
SELECT 
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
WHERE
  r.title ILIKE $1
ORDER
  by r.title
LIMIT 50;
