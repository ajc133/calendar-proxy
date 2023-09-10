.PHONY: server client

server:
	exec cargo watch -q -c -w src/ -x "run"

client:
	exec cargo watch -q -c -w examples/ -x "run --example quick_dev"
