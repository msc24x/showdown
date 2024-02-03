FROM ubuntu

ENV COMPILERS_PATH ""

# C & C++
RUN apt-get update && \
    apt-get install -y build-essential

RUN apt-get -y install curl sudo zlib1g-dev

# GOLANG
ENV GOLANG_VERSION 1.20.7
RUN curl -OLf https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go$GOLANG_VERSION.linux-amd64.tar.gz

ENV COMPILERS_PATH=$COMPILERS_PATH:/usr/local/go/bin

# Python
ENV PYTHON_VERSION 3.12.0
RUN curl -Of https://www.python.org/ftp/python/${PYTHON_VERSION}/Python-${PYTHON_VERSION}.tgz && \
    tar -xvzf Python-${PYTHON_VERSION}.tgz
RUN cd Python-${PYTHON_VERSION} && \
    ./configure \
        --prefix=/opt/python/${PYTHON_VERSION} \
        --enable-shared \
        --enable-optimizations \
        --enable-ipv6 \
        LDFLAGS=-Wl,-rpath=/opt/python/${PYTHON_VERSION}/lib,--disable-new-dtags
RUN cd Python-${PYTHON_VERSION} && \
    make && \
    sudo make install

ENV COMPILERS_PATH=$COMPILERS_PATH:/opt/python/$PYTHON_VERSION/bin/

# ISOLATE
RUN apt-get -y install git libcap-dev

RUN git clone https://github.com/ioi/isolate && \
    cd isolate && \
    sudo make isolate && \
    sudo make install

ENV ISOLATE_PATH=${pwd}/isolate/isolate

RUN rm -rf /tmp/*