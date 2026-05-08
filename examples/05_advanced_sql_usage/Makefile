.SILENT:
example: ../../bin/mimic
	yes | ../../bin/mimic copy \
	-p tableName="Table name: " \
	-p tableColumns="Columns definition (separated by comma): " \
	-p tableFields="Fields definition (separated by comma): " \
	-p tableValues="Values definition (separated by comma): " \
	-v tableName="users" \
	-v tableColumns="id UUID PRIMARY KEY DEFAULT gen_random_uuid(), name TEXT NOT NULL, created_at TIMESTAMPTZ DEFAULT NOW()" \
	-v tableFields="name" \
	-v tableValues="\"Foo Bar Baz\"" \
	-s .mimic/MIGRATION \
	-s .mimic/SEED