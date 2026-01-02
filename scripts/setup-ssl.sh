#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "============================================"
echo "   SSL/TLS Setup with Let's Encrypt"
echo "============================================"
echo ""

# Check if domain is provided
if [ -z "$1" ]; then
    echo -e "${RED}Error: Domain name required${NC}"
    echo "Usage: ./setup-ssl.sh yourdomain.com your@email.com"
    exit 1
fi

if [ -z "$2" ]; then
    echo -e "${RED}Error: Email required for Let's Encrypt${NC}"
    echo "Usage: ./setup-ssl.sh yourdomain.com your@email.com"
    exit 1
fi

DOMAIN=$1
EMAIL=$2

echo -e "${YELLOW}Domain: ${DOMAIN}${NC}"
echo -e "${YELLOW}Email: ${EMAIL}${NC}"
echo ""

# Step 1: Install cert-manager
echo -e "${YELLOW}[1/4] Installing cert-manager...${NC}"
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.3/cert-manager.yaml

# Wait for cert-manager to be ready
echo -e "${YELLOW}Waiting for cert-manager pods...${NC}"
kubectl wait --for=condition=ready pod -l app.kubernetes.io/instance=cert-manager -n cert-manager --timeout=300s

echo -e "${GREEN}cert-manager installed${NC}"

# Step 2: Create ClusterIssuer for Let's Encrypt
echo -e "${YELLOW}[2/4] Creating Let's Encrypt ClusterIssuer...${NC}"

cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: ${EMAIL}
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - http01:
        ingress:
          class: nginx
EOF

echo -e "${GREEN}ClusterIssuer created${NC}"

# Step 3: Update ingress with TLS
echo -e "${YELLOW}[3/4] Creating TLS Ingress...${NC}"

cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: k3s-app
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - ${DOMAIN}
    secretName: k3s-app-tls
  rules:
  - host: ${DOMAIN}
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: k3s-app-frontend
            port:
              number: 80
EOF

echo -e "${GREEN}TLS Ingress created${NC}"

# Step 4: Wait for certificate
echo -e "${YELLOW}[4/4] Waiting for certificate to be issued...${NC}"
echo "This may take 1-2 minutes..."

sleep 30

# Check certificate status
kubectl get certificate -A

echo ""
echo "============================================"
echo -e "${GREEN}   SSL Setup Complete!${NC}"
echo "============================================"
echo ""
echo "Your site should now be available at:"
echo -e "${GREEN}https://${DOMAIN}${NC}"
echo ""
echo "Certificate status:"
kubectl get certificate k3s-app-tls 2>/dev/null || echo "Certificate pending..."
echo ""
echo "If certificate is not ready, check:"
echo "  kubectl describe certificate k3s-app-tls"
echo "  kubectl describe certificaterequest"
echo "  kubectl logs -n cert-manager deployment/cert-manager"
echo ""
