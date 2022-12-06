# Random number generator collector

Collect random number and write to files under `logs` directory.

## Cross platform solution (Unix/MacOS/Windows)

- Install [Docker desktop](https://www.docker.com/products/docker-desktop) on Windows

- Install [Docker desktop](https://www.docker.com/products/docker-desktop) on MacOS

- Install [Docker engine](https://docs.docker.com/engine/install/centos/) on CentOS

- Install [Docker engine](https://docs.docker.com/engine/install/debian/) on Debain

- Install [Docker engine](https://docs.docker.com/engine/install/fedora/) on Fedora

- Install [Docker engine](https://docs.docker.com/engine/install/rhel/) on RHEL

- Install [Docker engine](https://docs.docker.com/engine/install/sles/) on SLES

- Install [Docker engine](https://docs.docker.com/engine/install/ubuntu/) on Ubuntu

### Quick start

- Launch service

    ```bash
    docker-compose up -d 
    ```

- Stop service

    ```bash
    docker-compose up -d 
    ```

### RNG Tool

- Build image

    ```bash
    docker build -f Dockerfile.rng-tool -t rng-tool .
    ```

- Launch tests

    ```bash
    docker-compose -f ./assets/docker-compose.yml up -d 
    ```

- Launch a few tests

    ```bash
    docker-compose -f ./assets/docker-compose.yml up -d set14 set15
    ```

#### Cross platform build

- Build image for compiler.

    ```bash
    docker build -f Dockerfile -t rng-tool:build .
    ```

- Compile 

    ```bash
    docker run -it --rm --entrypoint make -v ${PWD}/release:/workspace/release rng-tool:build cross-compile
    ```
