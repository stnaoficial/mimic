-- Migration version: 69701285

-- UP
CREATE TABLE IF NOT EXISTS carts (
    id UUID PRIMARY KEY,
	user_id UUID NOT NULL,
	product_id UUID NOT NULL,
	quantity INTEGER NOT NULL DEFAULT 1
);

-- DOWN
DROP TABLE IF EXISTS carts;