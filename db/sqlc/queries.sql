INSERT INTO users (name, dob) 
VALUES ($1, $2) 
RETURNING id, name, dob;

SELECT id, name, dob 
FROM users 
WHERE id = $1;

SELECT id, name, dob 
FROM users 
ORDER BY id;

SELECT id, name, dob 
FROM users 
ORDER BY id
LIMIT $1 OFFSET $2;

-
SELECT COUNT(*) 
FROM users;

UPDATE users 
SET name = $2, dob = $3 
WHERE id = $1 
RETURNING id, name, dob;

DELETE FROM users 
WHERE id = $1;
