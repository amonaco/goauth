-- Seed development/test data (DEV ONLY - DO NOT USE IN PRODUCTION)

-- Insert test roles
INSERT INTO roles (name, description) VALUES
  ('user', 'Standard user role'),
  ('admin', 'Administrator role'),
  ('superadmin', 'Super administrator role')
ON CONFLICT (name) DO NOTHING;

-- Insert test companies
INSERT INTO companies (name, slug, description) VALUES
  ('Acme Corp', 'acme-corp', 'Test company Acme'),
  ('Tech Startup', 'tech-startup', 'Test company Tech Startup')
ON CONFLICT (slug) DO NOTHING;

-- Insert test users
INSERT INTO users (email, password_hash, first_name, last_name) VALUES
  ('admin@example.com', '$2a$10$placeholder_hash_admin', 'Admin', 'User'),
  ('user@example.com', '$2a$10$placeholder_hash_user', 'Test', 'User'),
  ('dev@example.com', '$2a$10$placeholder_hash_dev', 'Dev', 'User')
ON CONFLICT (email) DO NOTHING;

-- Assign roles to users
INSERT INTO user_roles (user_id, company_id, role_id) VALUES
  (1, 1, 3), -- admin@example.com -> Acme Corp -> superadmin
  (2, 1, 1), -- user@example.com -> Acme Corp -> user
  (3, 2, 2), -- dev@example.com -> Tech Startup -> admin
  (3, 1, 1)  -- dev@example.com -> Acme Corp -> user
ON CONFLICT (user_id, company_id, role_id) DO NOTHING;
