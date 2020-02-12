-- name:Table-Existence
SELECT EXISTS
           (
               SELECT 1
               FROM information_schema.tables
               WHERE table_name = 'item'
           );

-- name:Table-Create
CREATE TABLE item
(
    id   uuid primary key,
    name varchar(300)
);

-- name:Table-Set-Owner
ALTER TABLE item OWNER TO $1;

-- name:Select-All
SELECT *
FROM item
ORDER BY name
LIMIT $1 OFFSET $2;

--name:Insert
INSERT INTO item(id, name)
VALUES (gen_random_uuid(), $1)