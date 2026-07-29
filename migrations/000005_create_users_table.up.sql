INSERT INTO roles (name, description) VALUES
    ('admin', 'Administrator with full access'),
    ('user', 'Regular user with limited access'),
    ('viewer', 'Read-only access')
ON CONFLICT (name) DO NOTHING;

INSERT INTO users (name, email, age, role_id) VALUES
    ('Admin', 'admin@example.com', 30, (SELECT id FROM roles WHERE name = 'admin')),
    ('User1', 'user1@example.com', 25, (SELECT id FROM roles WHERE name = 'user')),
    ('User2', 'user2@example.com', 28, (SELECT id FROM roles WHERE name = 'user'))
ON CONFLICT (email) DO NOTHING;
