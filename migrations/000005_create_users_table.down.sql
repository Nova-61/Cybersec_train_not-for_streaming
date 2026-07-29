DELETE FROM users WHERE email IN ('admin@example.com', 'user1@example.com', 'user2@example.com');
DELETE FROM roles WHERE name IN ('admin', 'user', 'viewer');