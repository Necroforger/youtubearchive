db="cmd/youtube-archive/archive.db"
tpl="templates/*.html"
static="static"
addr="127.0.0.1:8080"
pass="wriggles-lantern"


go run cmd/youtube-archive-server/main.go \
	-db "$db" \
	-templates "$tpl" \
	-static "$static" \
	-addr "$addr" \
	-pass "$pass"