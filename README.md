# nrunc: Experimental Testing Playground

Welcome to the **nrunc** repository! This project serves as an **experimental testing playground** for exploring and validating features related to [nanos](https://github.com/nanovms/nanos) and [urunc](https://github.com/nubificus/urunc) integration. Our goal is to experiment with various functionalities, test APIs, and gather insights.

## Purpose
The findings and results from our experiments will be pushed to the upstream repository [urunc](https://github.com/nubificus/urunc) once the APIs are finalized. This ensures that our contributions are integrated into the main project, benefiting the broader community.

## Usage

### Requirement

This project is intended for Linux environments only. To run the experiments, you need to have one of the following installed:

- **Docker**: [https://docs.docker.com/engine/install/debian/](https://docs.docker.com/engine/install/debian/)
- **Ops**: [https://ops.city](https://ops.city) Needed to build Nanos unikernels.
- **Go**: You need to have Go installed to build nrunc. Please follow the [official Go installation guide](https://golang.org/doc/install) for instructions.
- **VMM**: Being tested right now:
  - **Firecracker**: [https://github.com/firecracker-microvm/firecracker](https://github.com/firecracker-microvm/firecracker)
  - **Cloud Hypervisor**: [https://github.com/cloud-hypervisor/cloud-hypervisor](https://github.com/cloud-hypervisor/cloud-hypervisor)
  - **Qemu**: [https://www.qemu.org/download/#linux](https://www.qemu.org/download/#linux)

Make sure to follow the installation instructions for your chosen hypervisor before proceeding with the usage steps.

To use the **nrunc** project, follow these steps:

1. **Clone the Repository**:
   Clone the repository to your local machine using the following command:
   ```sh
   git clone https://github.com/rinor/nrunc.git
   cd nrunc
   ```

2. **Build and Install**:
   Make sure to install any required dependencies. You can typically do this with:
   ```sh
   make
   make install # or sudo make install
   ```

3. **Run Experiments**:
   Check some container files in [examples](examples) folder. Assuming you already have some nanos unikernels ready,
   ```sh
   ops image list
   +------+----------------------------+---------+-------------+
   | NAME |            PATH            |  SIZE   |  CREATEDAT  |
   +------+----------------------------+---------+-------------+
   | app  | /home/user/.ops/images/app | 14.9 MB | 4 hours ago |
   +------+----------------------------+---------+-------------+
   ```

   copy that `app` image to `examples` folder

   ```sh
   cd examples
   cp /home/user/.ops/images/app app
   ```

   to build a "container" for `firecracker` from [Containerfile.firecracker](examples/Containerfile.firecracker)

   ```dockerfile
   #syntax=harbor.nbfc.io/nubificus/pun:latest
   FROM scratch

   COPY app /nanos/image

   LABEL "com.urunc.unikernel.binary"="/nanos/kernel"
   LABEL "com.urunc.unikernel.cmdline"=""
   LABEL "com.urunc.unikernel.block"="/nanos/image"
   LABEL "com.urunc.unikernel.useDMBlock"="false"
   LABEL "com.urunc.unikernel.unikernelType"="nanos"
   LABEL "com.urunc.unikernel.hypervisor"="firecracker"
   ```

   let's build it

   ```sh
   docker buildx build --builder=default --output "type=image,oci-mediatypes=true" -f Containerfile.firecracker  -t "nanos/app:firecracker" .
   ```

   let's try to run it

   ```sh
   docker run --rm --name app_firecracker --runtime io.containerd.nrunc.v2 nanos/app:firecracker
   ```

4. **Known issues**:
   - `qemu` and `cloud-hypervisor` **do not** support **networking**
   - `firecracker` **does** support **networking**

5. **Contribute**:
   If you have insights or improvements, feel free to contribute by creating a pull request or an issue.
