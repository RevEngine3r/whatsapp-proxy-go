# Step 2: Configuration Management

## Objective
Implement flexible configuration system supporting CLI arguments, YAML files, and environment variables with proper validation.

## Tasks

### 1. Configuration Structure
Define complete config schema:
- Server settings (port, bind address)
- SOCKS5 upstream (host, port, auth)
- SSL settings (cert paths, generation)
- Logging (level, format, output)
- Metrics (enabled, port)

### 2. CLI Implementation with Cobra
- Root command setup
- Flag definitions for all settings
- Help text and examples
- Version command

### 3. YAML Config Support with Viper
- Config file parsing
- Default config path search
- Custom config file flag
- Config file generation command

### 4. Configuration Priority
1. CLI flags (highest)
2. Environment variables
3. Config file
4. Defaults (lowest)

### 5. Validation
- Port range validation
- Required field checks
- SOCKS5 URL parsing
- Path existence verification

## Implementation Details

### Config Structure Example
```go
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    SOCKS5   SOCKS5Config   `yaml:"socks5"`
    SSL      SSLConfig      `yaml:"ssl"`
    Logging  LoggingConfig  `yaml:"logging"`
    Metrics  MetricsConfig  `yaml:"metrics"`
}

type ServerConfig struct {
    Port       int    `yaml:"port"`
    BindAddr   string `yaml:"bind_addr"`
    IdleTimeout int   `yaml:"idle_timeout"`
}

type SOCKS5Config struct {
    Enabled  bool   `yaml:"enabled"`
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
}
```

### CLI Flags
```
--port, -p              Server port (default: 8443)
--bind                  Bind address (default: 0.0.0.0)
--config, -c            Config file path
--socks5-proxy          SOCKS5 proxy (format: socks5://[user:pass@]host:port)
--log-level             Logging level (debug, info, warn, error)
--metrics-port          Metrics endpoint port (default: 8199)
```

### Example YAML Config
```yaml
server:
  port: 8443
  bind_addr: 0.0.0.0
  idle_timeout: 300

socks5:
  enabled: true
  host: 127.0.0.1
  port: 1080
  username: ""
  password: ""

ssl:
  auto_generate: true
  cert_file: ""
  key_file: ""
  dns_names: []
  ip_addresses: []

logging:
  level: info
  format: text
  output: stdout

metrics:
  enabled: true
  port: 8199
```

## Testing
- Parse CLI arguments correctly
- Load YAML config file
- Environment variable override works
- Validation catches errors
- Priority order respected

## Deliverables
- internal/config package with complete implementation
- configs/config.example.yaml
- CLI command structure
- Unit tests for config validation
