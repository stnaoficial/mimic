.SILENT:

COUNT=$$(find ./migrations -maxdepth 1 -type f | wc -l)

example:
	../../bin/mimic dump \
	-v "count=$(COUNT)" \
	-v "timestamp=$$(date +%s)"