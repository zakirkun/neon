# Neon - Docker Swarm Management Tool

A CLI tool for automating Docker Swarm deployments with zero-downtime updates and advanced monitoring.

## Features

- Zero-downtime deployments with health checks
- Rolling updates with automatic rollback
- Service monitoring and metrics
- Volume backup and restore
- Network management
- Resource limits and scaling
- Configuration-based deployments

## Installation

```bash
go install github.com/zakirkun/neon@latest
```

## Quick Start

```bash
# Show version
neon --version

# Deploy with zero downtime
neon deploy rolling myapp --image nginx:latest --replicas 3 --update-delay 10s

# Deploy using config file
neon deploy config -f deploy.yaml

# List and manage resources
neon image list
neon network list
neon volume list
```

## Configuration

Default config location: `~/.neon/config.yaml`

```yaml
docker:
  registry: "registry.example.com"
  username: "user"
  password: "pass"

swarm:
  manager_node: "127.0.0.1:2377"
  network_name: "neon-network"

deploy:
  replicas: 3
  update_delay: "10s"
  rollback_delay: "5s"
  failure_action: "rollback"
```

## Deployment Configuration

Example `deploy.yaml`:

```yaml
services:
  webapp:
    image: registry.example.com/webapp:latest
    replicas: 3
    update_config:
      parallelism: 1
      delay: 10s
      order: start-first
      failure_action: rollback
    healthcheck:
      test: ["CMD-SHELL", "curl -f http://localhost/health || exit 1"]
      interval: 5s
      timeout: 3s
      retries: 3
    resources:
      limits:
        cpus: '0.5'
        memory: 512M
```

## Commands

### Deployment
```bash
# Zero-downtime deployment
neon deploy rolling <service> --image <image> [options]
  --replicas      Number of replicas (default: 3)
  --update-delay  Delay between updates (default: 10s)
  --image         New image to deploy

# Config-based deployment
neon deploy config -f deploy.yaml
```

### Resource Management
```bash
# Images
neon image list
neon image rm <image-id>
neon image pull <image-name>

# Networks
neon network list
neon network create <name>
neon network rm <name>

# Volumes
neon volume list
neon volume create <name>
neon volume rm <name>
```

### Monitoring
```bash
# Monitor service metrics
neon monitor service <service-id>
```

## Development

Requirements:
- Go 1.22 or later
- Docker Engine with Swarm mode enabled

Build from source:
```bash
git clone https://github.com/zakirkun/neon.git
cd neon
go build
```

## License

MIT License
