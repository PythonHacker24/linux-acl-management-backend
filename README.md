<div align="center">

# Linux ACL Management Interface - Backend Component

<img width="600" hegith="600" src="https://github.com/user-attachments/assets/a1625f58-0cd8-4df9-babc-31547b18d55a">

A robust web-based management interface for Linux Access Control Lists (ACLs), designed to enhance data protection and simplify ACL administration. This project provides a modern, user-friendly solution for managing file system permissions in Linux environments.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[View Documentation](https://pythonhacker24.github.io/linux-acl-management-docs/)

</div>

## Project Summary 

Institutional departments, such as the Biomedical Informatics (BMI) Department of Emory University School of Medicine, manage vast amounts of data, often reaching petabyte scales across multiple Linux-based storage servers. Researchers storing data in these systems need a streamlined way to modify ACLs to grant or revoke access for collaborators. Currently, the IT team at BMI is responsible for manually handling these ACL modifications, which is time-consuming, error-prone, and inefficient, especially as data volume and user demands grow. To address this challenge at BMI and similar institutions worldwide, a Web Management Interface is needed to allow users to modify ACLs securely. This solution would eliminate the burden on IT teams by enabling on-demand permission management while ensuring security and reliability. The proposed system will feature a robust and highly configurable backend, high-speed databases, orchestration daemons for file storage servers, and an intuitive frontend. The proposal includes an in-depth analysis of required components, high-level and low-level design considerations, technology selection, and the demonstration of a functional prototype as proof of concept. The goal is to deliver a production-ready, secure, scalable, and reliable system for managing ACLs across multiple servers hosting filesystems such as NFS, BeeGFS, and others. This solution will streamline access control management and prepare it for deployment at BMI and other institutions worldwide, significantly reducing the manual workload for IT teams.

## Features

- Intuitive web interface for ACL management
- High-performance backend written in Go
- Real-time ACL updates
- Comprehensive ACL reporting and visualization
- Integration with OpenLDAP for authentication

## Quick Start

### Prerequisites

- Go 1.20 or higher
- Docker (optional)
- Redis
- OpenLDAP server

### Local Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/PythonHacker24/linux-acl-management-backend.git
   cd linux-acl-management
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the application:
   ```bash
   go build -o acl-manager
   ```

### Production Build 

For production build, it is recommended to use the Makefile. This allows you to build the complete binary on locally for security purposes. Since the project is in development mode, complete local build is not possible since dependencies are managed via GitHub and external vendors. Tarball based complete local builds will be developed in later stages.

1. Clone the repository:
   ```bash
   git clone https://github.com/PythonHacker24/linux-acl-management.git
   cd linux-acl-management
   ```

2. Use make:
    ```bash
    make build
    ```

3. Execute the binary
    ```bash
    ./bin/laclm --config config.yaml
    ```

### Docker Testbench Deployment

A simulated environment has been developed on docker-compose for testing and experimenting purposes. It's not a production level build but a training ground for testing your config.yaml file for specific scenario. 

```bash
docker-compose up -d
```

A complete optional Docker based deployment option will be developed in later stages of development 

## Usage

1. Configure your settings in `config.yaml`

2. Start the server:
   ```bash
   ./laclm --config <config.yaml>
   ```

3. Access the api at `http://<ip-address>:<port>`

For detailed usage instructions, please refer to our [documentation](https://pythonhacker24.github.io/linux-acl-management/).

## Development

### Branches

- `main`: Production-ready code
- `development-v<version>`: Development branches for specific versions

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and development process.

## About

This project is developed as part of Google Summer of Code 2025, in collaboration with the Department of Biomedical Informatics at Emory University.

### Team

- **Contributor:** Aditya Patil
- **Mentors:** 
  - Robert Tweedy
  - Mahmoud Zeydabadinezhad, PhD

### Technologies

- **Backend:** Golang, net/http
- **API:** gRPC, REST
- **Infrastructure:** Docker, Redis, OpenLDAP
- **Packaging:** Tarball

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Department of Biomedical Informatics, Emory University
- Google Summer of Code Program
- Open Source Community
