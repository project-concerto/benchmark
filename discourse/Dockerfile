FROM ubuntu:latest

WORKDIR /benchmark/discourse/
COPY . /benchmark/discourse/

# Setup JAVA_HOME -- useful for docker commandline
ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64/
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
    && apt-get install -y git \
    && apt-get install -y wget \
    && apt-get install -y curl \
    && apt-get install -y brotli \
    && apt-get install -y libpq-dev \
    && apt-get install -y maven \
    # Install OpenJDK-8
    && apt-get install -y openjdk-8-jdk \
    && apt-get install -y ant \
    # Fix certificate issues
    && apt-get install -y ca-certificates-java \
    && update-ca-certificates -f \
    && export JAVA_HOME
