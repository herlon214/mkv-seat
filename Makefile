release:
	gox -output="./releases/{{.Dir}}_{{.OS}}_{{.Arch}}" github.com/herlon214/mkv-seat/cmd/mkv-seat