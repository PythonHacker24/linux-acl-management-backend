# Linux ACL Management Interface - Backend Component

A robust web-based management interface for Linux Access Control Lists (ACLs), designed to enhance data protection and simplify ACL administration. This project provides a modern, user-friendly solution for managing file system permissions in Linux environments.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[View Documentation](https://pythonhacker24.github.io/linux-acl-management/)

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
   git clone https://github.com/yourusername/linux-acl-management.git
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
   git clone https://github.com/yourusername/linux-acl-management.git
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

1. Start the server:
   ```bash
   ./acl-manager
   ```

2. Access the api at `http://<ip-address>:<port>`

3. Configure your settings in `config.yaml`

For detailed usage instructions, please refer to our [documentation](https://pythonhacker24.github.io/linux-acl-management/).

## Project Structure

```
.
├── cmd/          # Application entry points
├── internal/     # Private application code
├── pkg/          # Public library code
├── api/          # API definitions and handlers
├── docs/         # Documentation
└── deployments/  # Deployment configurations
```

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
