example: ../../bin/mimic
	yes | ../../bin/mimic \
	-v table="users" \
	-v fields="id UUID PRIMARY KEY, name TEXT NOT NULL, created_at TIMESTAMPTZ DEFAULT NOW()" \
	.mimic/ .

	yes | ../../bin/mimic \
	-v table="products" \
	-v fields="id UUID PRIMARY KEY, name TEXT NOT NULL, price INTEGER NOT NULL, currency SMALLINT NOT NULL, created_at TIMESTAMPTZ DEFAULT NOW()" \
	.mimic/ .

	yes | ../../bin/mimic \
	-v table="carts" \
	-v fields="id UUID PRIMARY KEY, user_id UUID NOT NULL, product_id UUID NOT NULL, quantity INTEGER NOT NULL DEFAULT 1" \
	.mimic/ .