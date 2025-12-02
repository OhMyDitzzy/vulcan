!/bin/bash

# Generate self-signed certificate for HTTPS
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes \
  -subj "/C=US/ST=State/L=City/O=Vulcan/CN=localhost"

echo "Self-signed certificate generated: cert.pem and key.pem"