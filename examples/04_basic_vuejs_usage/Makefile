.SILENT:
example: ../../bin/mimic
	yes | ../../bin/mimic \
	--debug \
	--expr-open="<%" \
	--expr-close="%>" \
	-v name="Foo Bar Baz" \
	-v date="$$(date +%Y-%m-%d)" \
	.mimic/ .