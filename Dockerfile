# FIRST STAGE:  build the app.
FROM registry.met.no/baseimg/ubuntu:22.04 AS build-app
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get -y upgrade && apt-get -y dist-upgrade
RUN apt-get -y install apt-utils pkg-config curl git

WORKDIR /build/app

RUN apt-get -y install cmake g++ wget unzip build-essential
RUN wget -O opencv.zip https://github.com/opencv/opencv/archive/4.8.1.zip
RUN wget -O opencv_contrib.zip https://github.com/opencv/opencv_contrib/archive/4.8.1.zip
RUN unzip opencv.zip
RUN unzip opencv_contrib.zip

RUN mkdir -p build
WORKDIR /build/app/build
RUN cmake -D OPENCV_GENERATE_PKGCONFIG=YES ../opencv-4.8.1/
RUN cmake --build .
RUN make install
RUN ldconfig -v

ENV GOPATH=/go
RUN curl -L https://go.dev/dl/go1.21.1.linux-amd64.tar.gz | tar xz --directory /usr/local
ENV PATH="/go/bin:/usr/local/go/bin:${PATH}"

RUN ldconfig -v
COPY . /tmp/roadlabels
WORKDIR /tmp/roadlabels
ARG S3SecretKey
ARG S3AccessKey
ENV S3SecretKey $S3SecretKey
ENV S3AccessKey $S3AccessKey

RUN go mod tidy && make build && make install

# Second stage. Install 
FROM registry.met.no/baseimg/ubuntu:22.04
ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get -y upgrade && apt-get -y dist-upgrade
RUN apt-get -y install apt-utils pkg-config python3-netcdf4 python3-pyproj curl

COPY --from=build-app /usr/local /usr/local
RUN ldconfig -v
COPY --from=build-app /tmp/roadlabels/exttools /app/exttools
ENV PYTHONPATH=/app/extools
ENV HDF5_USE_FILE_LOCKING=FALSE
RUN mkdir -p /lustre
ADD entry.sh /usr/local/bin/entry.sh
ADD var/lib/roadlabels/roadcams.db /roadlabels/roadcams.db
ADD var/lib/roadlabels/userdb-empty.db /roadlabels/userdb-empty.db
RUN mkdir -p /var/lib/roadlabels
RUN chown -R nobody.nogroup /var/lib/roadlabels
USER nobody:nogroup

EXPOSE 25260

ENTRYPOINT /usr/local/bin/entry.sh

