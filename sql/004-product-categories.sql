CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(256) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO categories (code, name) VALUES 
    ('CLOTHING', 'Clothing'),
    ('SHOES', 'Shoes'),
    ('ACCESSORIES', 'Accessories')
ON CONFLICT (code) DO NOTHING;

ALTER TABLE products ADD COLUMN category_id INTEGER
REFERENCES categories(id) ON DELETE SET NULL;

UPDATE products 
SET category_id = (SELECT id FROM categories WHERE code = 'CLOTHING')
WHERE code IN ('PROD001', 'PROD004', 'PROD007');

UPDATE products 
SET category_id = (SELECT id FROM categories WHERE code = 'SHOES')
WHERE code IN ('PROD002', 'PROD006');

UPDATE products 
SET category_id = (SELECT id FROM categories WHERE code = 'ACCESSORIES')
WHERE code IN ('PROD003', 'PROD005', 'PROD008');
