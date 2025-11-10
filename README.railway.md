# TURN Server - Railway Deployment

WebRTC TURN (Traversal Using Relays around NAT) server for Sweven Games cloud gaming platform.

## Technology Stack
- **Language**: Go 1.18
- **Library**: Pion TURN/STUN v3
- **Protocol**: UDP

## Important Notes

### UDP Port Requirements
TURN servers require **UDP port ranges** for relay connections (typically 49152-65535). Railway's free tier may have limitations on UDP traffic and port ranges.

### Recommended Deployment
For production TURN servers, consider:
1. **Dedicated Server**: Deploy on a VPS with full UDP support (DigitalOcean, Linode, AWS EC2)
2. **Port Range**: Open ports 49152-65535 for UDP relay traffic
3. **Public IP**: TURN servers need a stable public IP address

## Railway Deployment Steps

### 1. Build the Binary
```bash
cd "d:\cloud gaming\backend-services\turn-server-new"
go build -o turn-server cmd/main.go
```

### 2. Create Railway Project
```bash
railway login
railway init
```

### 3. Set Environment Variables
```bash
railway variables set TURN_USERNAME=swevengames
railway variables set TURN_PASSWORD=your_secure_password
railway variables set PUBLIC_IP=<your-railway-assigned-ip>
railway variables set TURN_PORT=3478
```

### 4. Deploy
```bash
railway up
```

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `TURN_USERNAME` | Username for TURN authentication | - | Yes |
| `TURN_PASSWORD` | Password for TURN authentication | - | Yes |
| `PUBLIC_IP` | Public IP address of server | 0.0.0.0 | No |
| `TURN_PORT` | TURN server port | 3478 | No |
| `WORKER_ID` | Optional Supabase worker ID for monitoring | - | No |
| `TM_PROJECT` | Optional Supabase project URL | - | No |
| `TM_ANONKEY` | Optional Supabase anon key | - | No |

## Port Configuration

### Primary TURN Port
- **Port 3478**: Main TURN server port (UDP)

### Relay Port Range
- **Ports 60000-65535**: Used for WebRTC relay connections
- These ports are configured in the code (`min = 60000`, `max = 65535`)

## Testing TURN Server

### Using turnutils
```bash
# Install turnutils
sudo apt-get install coturn

# Test TURN server
turnutils_uclient -v -u swevengames -w your_password your-turn-server.railway.app 3478
```

### Using WebRTC Test Page
Visit: https://webrtc.github.io/samples/src/content/peerconnection/trickle-ice/

Add your TURN server:
```
turn:your-turn-server.railway.app:3478
Username: swevengames
Password: your_password
```

## Production Deployment Recommendations

### Option 1: Railway (Limited)
Railway deployment may work but has limitations:
- UDP traffic support varies
- Port range limitations
- May not be suitable for high traffic

### Option 2: VPS with Full UDP Support (Recommended)
```bash
# On Ubuntu/Debian server
cd /opt
git clone <your-repo>
cd turn-server-new

# Build
go build -o turn-server cmd/main.go

# Set environment variables
export TURN_USERNAME=swevengames
export TURN_PASSWORD=your_password
export PUBLIC_IP=$(curl -s ifconfig.me)
export TURN_PORT=3478

# Run with systemd
sudo cp service.sh /etc/systemd/system/edge-turn.service
sudo systemctl daemon-reload
sudo systemctl enable edge-turn
sudo systemctl start edge-turn
```

### Firewall Configuration
```bash
# Allow TURN port
sudo ufw allow 3478/udp

# Allow relay port range
sudo ufw allow 60000:65535/udp
```

## Integration with Frontend

Update your frontend WebRTC configuration:

```typescript
const iceServers = [
  {
    urls: 'stun:stun.l.google.com:19302'
  },
  {
    urls: 'turn:your-turn-server.railway.app:3478',
    username: 'swevengames',
    credential: 'your_password'
  }
]
```

## Monitoring

The TURN server includes optional Supabase monitoring:
- Set `WORKER_ID`, `TM_PROJECT`, and `TM_ANONKEY` environment variables
- Server will ping Supabase every 10 seconds with status

## Logging

TURN server logs all STUN/TURN messages:
- Inbound STUN packets
- Outbound STUN packets
- Connection requests with source addresses

## Security Considerations

1. **Strong Password**: Use a strong, unique password for TURN authentication
2. **Rotate Credentials**: Regularly rotate TURN credentials
3. **Rate Limiting**: Consider implementing rate limiting for production
4. **Firewall Rules**: Only open necessary ports
5. **Monitor Traffic**: Watch for abuse and unusual traffic patterns

## Alternative TURN Servers

If Railway doesn't meet your needs, consider:
1. **Coturn**: https://github.com/coturn/coturn
2. **Twilio NAT Traversal**: https://www.twilio.com/stun-turn
3. **Xirsys**: https://xirsys.com/
4. **Managed TURN Services**: Metered, Cloudflare Calls
