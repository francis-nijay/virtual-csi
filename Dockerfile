# Use rocky Linux as the base image
FROM rockylinux:8

# Copy the files from the host to the container
COPY "vcsi-driver" .

RUN yum install -y \
    e2fsprogs \
    which \
    xfsprogs \
    device-mapper-multipath \
    libaio \
    numactl \
    libuuid \
    e4fsprogs \
    nfs-utils \
    && \
    yum clean all \
    && \
    rm -rf /var/cache/run

# validate some cli utilities are found
RUN which mkfs.ext4
RUN which mkfs.xfs

# Set the command to run when the container starts
ENTRYPOINT ["/vcsi-driver"]