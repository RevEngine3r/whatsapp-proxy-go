# Configuration Reference

Complete reference for all configuration options in WhatsApp Proxy Go.

## Table of Contents

- [Configuration Methods](#configuration-methods)
- [Configuration Priority](#configuration-priority)
- [Server Configuration](#server-configuration)
- [SOCKS5 Configuration](#socks5-configuration)
- [SSL/TLS Configuration](#ssltls-configuration)
- [Logging Configuration](#logging-configuration)
- [Metrics Configuration](#metrics-configuration)
- [Environment Variables](#environment-variables)
- [CLI Flags](#cli-flags)
- [Configuration Examples](#configuration-examples)
- [Best Practices](#best-practices)

## Configuration Methods

The proxy supports three configuration methods:

1. **YAML Configuration File** (recommended)
2. **Command-line Flags**
3. **Environment Variables**

### Configuration Priority

When the same setting is defined in multiple places, the following priority applies (highest to lowest):

1. Command-line flags
2. Environment variables
3. Configuration file
4. Default values

## Server Configuration

### `server.port`

**Type:** `int`  
**Default:** `8443`  
**Description:** Port number for the proxy server to listen on. This single port handles all protocols (HTTP, HTTPS, Jabber/XMPP).

```yaml
server:
  port: 8443
```

**CLI Flag:** `--port`  
**Environment Variable:** `PROXY_PORT`

**Notes:**
- Ports below 1024 require root/administrator privileges on Unix systems
- Port 8443 is commonly used for HTTPS proxies
- Ensure the port is not already in use

### `server.bind_addr`

**Type:** `string`  
**Default:** `0.0.0.0`  
**Description:** IP address to bind the server to.

```yaml
server:
  bind_addr: 0.0.0.0
```

**CLI Flag:** `--bind-addr`  
**Environment Variable:** `PROXY_BIND_ADDR`

**Common Values:**
- `0.0.0.0` - Listen on all network interfaces (public access)
- `127.0.0.1` - Listen only on localhost (local access only)
- Specific IP - Listen on a specific network interface

### `server.idle_timeout`

**Type:** `int`  
**Default:** `300` (5 minutes)  
**Description:** Connection idle timeout in seconds. Connections with no activity for this duration will be closed.

```yaml
server:
  idle_timeout: 300
```

**Environment Variable:** `PROXY_IDLE_TIMEOUT`

**Recommendations:**
- Standard use: 300-600 seconds (5-10 minutes)
- High-traffic: 60-120 seconds (1-2 minutes)
- Long-polling apps: 900-1800 seconds (15-30 minutes)

### `server.max_connections`

**Type:** `int`  
**Default:** `1000`  
**Description:** Maximum number of concurrent connections. Set to 0 for unlimited (not recommended in production).

```yaml
server:
  max_connections: 1000
```

**Environment Variable:** `PROXY_MAX_CONNECTIONS`

**Sizing Guide:**
- Small deployment: 100-500
- Medium deployment: 500-2000
- Large deployment: 2000-10000
- Enterprise: 10000+

**System Requirements:**
- Each connection uses ~1-2MB of memory
- Ensure your system's file descriptor limit (ulimit) is set appropriately
- Linux: Typically needs `ulimit -n` set to at least 2x max_connections

## SOCKS5 Configuration

### `socks5.enabled`

**Type:** `bool`  
**Default:** `false`  
**Description:** Enable upstream SOCKS5 proxy for all outbound connections.

```yaml
socks5:
  enabled: true
```

**CLI Flag:** `--socks5-enabled`  
**Environment Variable:** `SOCKS5_ENABLED`

### `socks5.host`

**Type:** `string`  
**Default:** `""`  
**Description:** SOCKS5 proxy server hostname or IP address.

```yaml
socks5:
  host: 127.0.0.1
```

**CLI Flag:** Part of `--socks5-proxy` URL  
**Environment Variable:** `SOCKS5_HOST`

### `socks5.port`

**Type:** `int`  
**Default:** `1080`  
**Description:** SOCKS5 proxy server port.

```yaml
socks5:
  port: 1080
```

**CLI Flag:** Part of `--socks5-proxy` URL  
**Environment Variable:** `SOCKS5_PORT`

### `socks5.username`

**Type:** `string`  
**Default:** `""`  
**Description:** SOCKS5 authentication username. Leave empty if authentication is not required.

```yaml
socks5:
  username: myuser
```

**CLI Flag:** Part of `--socks5-proxy` URL  
**Environment Variable:** `SOCKS5_USERNAME`

**Security Note:** Consider using environment variables for credentials instead of storing in config files.

### `socks5.password`

**Type:** `string`  
**Default:** `""`  
**Description:** SOCKS5 authentication password. Leave empty if authentication is not required.

```yaml
socks5:
  password: mypassword
```

**CLI Flag:** Part of `--socks5-proxy` URL  
**Environment Variable:** `SOCKS5_PASSWORD`

**Security Note:** Always use environment variables or secrets management for passwords.

### `socks5.timeout`

**Type:** `int`  
**Default:** `30`  
**Description:** Connection timeout in seconds when connecting to the SOCKS5 proxy.

```yaml
socks5:
  timeout: 30
```

**Environment Variable:** `SOCKS5_TIMEOUT`

## SSL/TLS Configuration

### `ssl.auto_generate`

**Type:** `bool`  
**Default:** `true`  
**Description:** Automatically generate self-signed SSL certificates on startup.

```yaml
ssl:
  auto_generate: true
```

**CLI Flag:** `--ssl-auto-generate`  
**Environment Variable:** `SSL_AUTO_GENERATE`

**When to use:**
- Development and testing: `true` (convenient)
- Production with custom certs: `false`
- Internal use: `true` (acceptable)
- Public-facing: `false` (use proper certificates)

### `ssl.cert_file`

**Type:** `string`  
**Default:** `""`  
**Description:** Path to custom SSL certificate file (PEM format). Only used when `auto_generate` is `false`.

```yaml
ssl:
  auto_generate: false
  cert_file: /etc/ssl/certs/proxy.crt
```

**CLI Flag:** `--ssl-cert`  
**Environment Variable:** `SSL_CERT_FILE`

### `ssl.key_file`

**Type:** `string`  
**Default:** `""`  
**Description:** Path to custom SSL private key file (PEM format). Only used when `auto_generate` is `false`.

```yaml
ssl:
  auto_generate: false
  key_file: /etc/ssl/private/proxy.key
```

**CLI Flag:** `--ssl-key`  
**Environment Variable:** `SSL_KEY_FILE`

**Security:** Ensure private key file has restrictive permissions (chmod 600).

### `ssl.dns_names`

**Type:** `[]string`  
**Default:** `["localhost"]`  
**Description:** DNS names to include in Subject Alternative Names (SANs) for auto-generated certificates.

```yaml
ssl:
  dns_names:
    - localhost
    - proxy.example.com
    - whatsapp-proxy.local
```

**Environment Variable:** `SSL_DNS_NAMES` (comma-separated)

**Important:** Include all hostnames that clients will use to connect.

### `ssl.ip_addresses`

**Type:** `[]string`  
**Default:** `["127.0.0.1"]`  
**Description:** IP addresses to include in Subject Alternative Names (SANs) for auto-generated certificates.

```yaml
ssl:
  ip_addresses:
    - 127.0.0.1
    - 192.168.1.100
```

**Environment Variable:** `SSL_IP_ADDRESSES` (comma-separated)

### `ssl.validity_days`

**Type:** `int`  
**Default:** `365`  
**Description:** Validity period in days for auto-generated certificates.

```yaml
ssl:
  validity_days: 365
```

**Environment Variable:** `SSL_VALIDITY_DAYS`

**Recommendations:**
- Development: 365 days (1 year)
- Testing: 90 days
- Production: Use proper certificates with certificate authority

## Logging Configuration

### `logging.level`

**Type:** `string`  
**Default:** `info`  
**Description:** Minimum log level to output.

```yaml
logging:
  level: info
```

**CLI Flag:** `--log-level`  
**Environment Variable:** `LOG_LEVEL`

**Valid Values:**
- `debug` - Verbose output, includes protocol detection details
- `info` - Standard operational messages (recommended)
- `warn` - Warning messages only
- `error` - Error messages only

### `logging.format`

**Type:** `string`  
**Default:** `text`  
**Description:** Log output format.

```yaml
logging:
  format: text
```

**CLI Flag:** `--log-format`  
**Environment Variable:** `LOG_FORMAT`

**Valid Values:**
- `text` - Human-readable format (good for console)
- `json` - Structured JSON format (good for log aggregation)

**When to use JSON:**
- Production environments
- Log aggregation systems (ELK, Splunk, etc.)
- Automated log parsing
- Cloud deployments

### `logging.output`

**Type:** `string`  
**Default:** `stdout`  
**Description:** Log output destination.

```yaml
logging:
  output: stdout
```

**CLI Flag:** `--log-output`  
**Environment Variable:** `LOG_OUTPUT`

**Valid Values:**
- `stdout` - Standard output (console)
- `stderr` - Standard error
- `/path/to/file.log` - Write to file

**File Logging Tips:**
- Use absolute paths
- Ensure directory exists and is writable
- Set up log rotation (logrotate on Linux)
- Monitor disk space

## Metrics Configuration

### `metrics.enabled`

**Type:** `bool`  
**Default:** `true`  
**Description:** Enable Prometheus-compatible metrics endpoint.

```yaml
metrics:
  enabled: true
```

**CLI Flag:** `--metrics-enabled`  
**Environment Variable:** `METRICS_ENABLED`

### `metrics.port`

**Type:** `int`  
**Default:** `8199`  
**Description:** Port for the metrics HTTP server.

```yaml
metrics:
  port: 8199
```

**CLI Flag:** `--metrics-port`  
**Environment Variable:** `METRICS_PORT`

### `metrics.bind_addr`

**Type:** `string`  
**Default:** `127.0.0.1`  
**Description:** Bind address for the metrics server.

```yaml
metrics:
  bind_addr: 127.0.0.1
```

**CLI Flag:** `--metrics-bind-addr`  
**Environment Variable:** `METRICS_BIND_ADDR`

**Security Recommendation:** Use `127.0.0.1` (localhost only) unless you need external access to metrics.

### Available Metrics

The `/metrics` endpoint exposes the following metrics in OpenMetrics format:

- `whatsapp_proxy_connections_total` - Total connection count (counter)
- `whatsapp_proxy_connections_active` - Active connections (gauge)
- `whatsapp_proxy_connections_failed` - Failed connections (counter)
- `whatsapp_proxy_protocol_connections{protocol}` - Connections by protocol (counter)
- `whatsapp_proxy_bytes_sent_total` - Total bytes sent (counter)
- `whatsapp_proxy_bytes_received_total` - Total bytes received (counter)
- `whatsapp_proxy_errors_total` - Total errors (counter)
- `whatsapp_proxy_uptime_seconds` - Server uptime (gauge)

## Environment Variables

All configuration options can be set via environment variables:

```bash
# Server
export PROXY_PORT=8443
export PROXY_BIND_ADDR=0.0.0.0
export PROXY_IDLE_TIMEOUT=300
export PROXY_MAX_CONNECTIONS=1000

# SOCKS5
export SOCKS5_ENABLED=true
export SOCKS5_HOST=127.0.0.1
export SOCKS5_PORT=1080
export SOCKS5_USERNAME=myuser
export SOCKS5_PASSWORD=mypass
export SOCKS5_TIMEOUT=30

# SSL
export SSL_AUTO_GENERATE=true
export SSL_CERT_FILE=/path/to/cert.pem
export SSL_KEY_FILE=/path/to/key.pem
export SSL_DNS_NAMES=localhost,proxy.local
export SSL_IP_ADDRESSES=127.0.0.1,192.168.1.1
export SSL_VALIDITY_DAYS=365

# Logging
export LOG_LEVEL=info
export LOG_FORMAT=text
export LOG_OUTPUT=stdout

# Metrics
export METRICS_ENABLED=true
export METRICS_PORT=8199
export METRICS_BIND_ADDR=127.0.0.1
```

## CLI Flags

All configuration options are available as command-line flags:

```bash
./whatsapp-proxy \
  --config config.yaml \
  --port 8443 \
  --bind-addr 0.0.0.0 \
  --socks5-proxy socks5://user:pass@127.0.0.1:1080 \
  --ssl-auto-generate \
  --log-level info \
  --log-format text \
  --log-output stdout \
  --metrics-enabled \
  --metrics-port 8199 \
  --metrics-bind-addr 127.0.0.1
```

### Special Flags

- `--config` - Path to YAML configuration file
- `--version` - Show version information
- `--help` - Show help message

### SOCKS5 Proxy URL Format

The `--socks5-proxy` flag accepts a URL in the following formats:

```bash
# Without authentication
--socks5-proxy socks5://127.0.0.1:1080

# With authentication
--socks5-proxy socks5://username:password@127.0.0.1:1080

# Using hostname
--socks5-proxy socks5://proxy.example.com:1080
```

## Configuration Examples

### Minimal Configuration

```yaml
server:
  port: 8443
```

### Development Setup

```yaml
server:
  port: 8443
  bind_addr: 127.0.0.1

socks5:
  enabled: false

ssl:
  auto_generate: true

logging:
  level: debug
  format: text

metrics:
  enabled: true
```

### Production with SOCKS5

```yaml
server:
  port: 8443
  bind_addr: 0.0.0.0
  idle_timeout: 600
  max_connections: 5000

socks5:
  enabled: true
  host: proxy.example.com
  port: 1080
  username: proxyuser
  password: ${SOCKS5_PASSWORD}  # From environment
  timeout: 30

ssl:
  auto_generate: false
  cert_file: /etc/ssl/certs/proxy.crt
  key_file: /etc/ssl/private/proxy.key

logging:
  level: warn
  format: json
  output: /var/log/whatsapp-proxy/proxy.log

metrics:
  enabled: true
  port: 8199
  bind_addr: 127.0.0.1
```

### High-Performance Setup

```yaml
server:
  port: 443
  bind_addr: 0.0.0.0
  idle_timeout: 300
  max_connections: 10000

socks5:
  enabled: true
  host: 127.0.0.1
  port: 1080
  timeout: 15

ssl:
  auto_generate: false
  cert_file: /etc/ssl/certs/proxy.crt
  key_file: /etc/ssl/private/proxy.key

logging:
  level: error
  format: json
  output: /var/log/whatsapp-proxy/proxy.log

metrics:
  enabled: true
  port: 8199
  bind_addr: 127.0.0.1
```

## Best Practices

### Security

1. **Store credentials securely**
   - Use environment variables for passwords
   - Never commit credentials to version control
   - Use secrets management systems in production

2. **Restrict metrics endpoint**
   - Bind to `127.0.0.1` or use firewall rules
   - Use reverse proxy with authentication if external access needed

3. **Use proper certificates**
   - Self-signed certificates OK for internal/development use
   - Use proper CA-signed certificates for production
   - Regularly rotate certificates

4. **File permissions**
   ```bash
   chmod 600 config.yaml          # Config file
   chmod 600 /etc/ssl/private/*.key  # Private keys
   chmod 644 /etc/ssl/certs/*.crt    # Certificates
   ```

### Performance

1. **Tune connection limits**
   - Set `max_connections` based on expected load
   - Monitor system resources (CPU, memory, file descriptors)
   - Adjust `idle_timeout` based on usage patterns

2. **System limits**
   ```bash
   # Increase file descriptor limit
   ulimit -n 65535
   
   # Make permanent (add to /etc/security/limits.conf)
   * soft nofile 65535
   * hard nofile 65535
   ```

3. **Logging in production**
   - Use `warn` or `error` level to reduce I/O
   - Use JSON format for better parsing
   - Set up log rotation
   - Monitor log file sizes

### Monitoring

1. **Enable metrics**
   - Always enable metrics endpoint
   - Integrate with Prometheus or similar
   - Set up alerts for:
     - High connection count
     - High error rate
     - High memory usage

2. **Health checks**
   - Use `/health` endpoint for load balancer health checks
   - Monitor uptime and restarts
   - Set up automated recovery

### Reliability

1. **Use systemd or equivalent**
   - Auto-restart on failure
   - Start on system boot
   - Resource limits

2. **Graceful shutdown**
   - Proxy handles SIGTERM/SIGINT gracefully
   - Allows 30 seconds for active connections to finish
   - Clean resource cleanup

3. **Backup configuration**
   - Regularly backup config files
   - Version control your configurations
   - Document custom settings
