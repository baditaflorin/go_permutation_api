-- Initialize database schema for permutation API

-- Create elements table
CREATE TABLE IF NOT EXISTS elements (
    id SERIAL PRIMARY KEY,
    value VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample data
INSERT INTO elements (value) VALUES
    ('apple'),
    ('banana'),
    ('cherry'),
    ('date'),
    ('elderberry');

-- Create index on value column for faster queries
CREATE INDEX IF NOT EXISTS idx_elements_value ON elements(value);

-- Grant permissions
GRANT ALL PRIVILEGES ON TABLE elements TO postgres;
GRANT USAGE, SELECT ON SEQUENCE elements_id_seq TO postgres;
