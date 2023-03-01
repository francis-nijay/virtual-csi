
# Builds the binary
build:
	go build .


# Generates the podman image
docker: build
	docker build -t vcsi-driver .

# Generates the podman image
podman: build
	podman build -t vcsi-driver .