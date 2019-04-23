FROM golang:latest
RUN echo "deb http://deb.debian.org/debian stretch-backports main" >> /etc/apt/sources.list && \
    apt-get update && \
    apt-get -y install libz-dev libbz2-dev libsnappy-dev && \
    apt-get -y -t stretch-backports install librocksdb-dev && \
    rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/facebook/rocksdb.git /tmp/rocksdb

WORKDIR /tmp/rocksdb
# RUN make static_lib
RUN make shared_lib
RUN make install
RUN rm -rf /tmp/rocksdb

ENV CGO_CFLAGS  "-I/usr/local/include"
ENV CGO_LDFLAGS "-L/usr/local -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy"

WORKDIR /mas
COPY . /mas

RUN ls

ARG username=${username:-""}
ENV username=${username}

ARG password=${password:-""}
ENV password=${password}

RUN echo "machine gitlab.zalopay.vn login ${username} password ${password}" > ~/.netrc
# RUN ping localhost
RUN go mod download
RUN rm -rf ~/.netrc 
# Build the binary.
RUN GOOS=linux GOARCH=amd64 GO111MODULE=on go build -ldflags="-w -s" -o /go/bin/mas
RUN go get github.com/mattn/goreman
# ENTRYPOINT ["./run.sh"]
############################
# STEP 2 build a small image
############################
# FROM scratch
# # Import the user and group files from the builder.
# COPY --from=builder /etc/passwd /etc/passwd

# WORKDIR /mas
# # Copy our static executable.
# COPY --from=builder /go/bin/mas ./

# # Use an unprivileged user.
# USER appuser
# # Run the hello binary.



