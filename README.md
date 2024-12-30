# Inventory Platform Collector
High-performance data collection service built in Go for real-time inventory data processing from eBay's API. Part of the Inventory Analytics Platform suite.

## Features
- Real-time eBay API data collection with rate limiting
- gRPC streaming services
- Kafka event distribution
- Concurrent processing using Go routines
- Comprehensive monitoring and metrics
- Docker containerization

## Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose
- Kafka
- Protocol Buffer compiler (protoc)

## Monitoring
- Prometheus metrics exposed at `/metrics` endpoint
- Health check endpoint at `/health`
- Full OpenTelemetry tracing integration for observability

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.





