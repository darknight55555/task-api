INSERT INTO tasks (title)
VALUES ('learn go')
RETURNING id, title, done, created_at;

SELECT id, title, done, created_at 
FROM tasks
ORDER BY id;

SELECT id, title, done, created_at
FROM tasks
WHERE id = 1;

UPDATE tasks 
SET title = 'learn hard', done = true
WHERE id = 1
RETURNING id, title, done, created_at;

DELETE FROM tasks
WHERE id = 1;