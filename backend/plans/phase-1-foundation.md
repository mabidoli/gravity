# Phase 1: Foundation & Infrastructure

**Sprint:** 1 (Weeks 1-2)  
**Focus:** Funnel logic, session management, and infrastructure setup  
**Status:** Ready for Implementation

---

## Overview

Phase 1 establishes the foundational infrastructure for the Funil Rápido platform. This includes the core funnel logic with URL protection, session management, custom domain support, and database schema design. These components are critical dependencies for all subsequent features.

---

## Objectives

1. Implement secure funnel session management with URL protection
2. Set up custom domain and SSL infrastructure
3. Design and implement database schema
4. Create development and deployment environments
5. Establish CI/CD pipeline

---

## Features

### Feature 2: Lógica de Funil com Proteção de URL

**Priority:** CRITICAL  
**Complexity:** Medium  
**Estimated Time:** 1 week

#### Description

Implement the core funnel flow logic that manages user progression through funnel stages (Landing → Sales → Upsell → Confirmation) with robust URL protection to prevent unauthorized access to protected pages.

#### Key Requirements

**Session Management:**
- Create session on first funnel access
- Store session in Redis with 12-hour TTL
- Track funnel state (viewing_landing, viewing_sales, purchased, upsell_shown, completed)
- Associate session with IP address for security

**URL Protection:**
- Block direct access to Sales, Upsell, and Confirmation pages
- Redirect unauthorized access to Landing Page
- Prevent back button navigation from Upsell page
- Validate session state before serving protected pages

**Security:**
- Session hijacking prevention via IP validation
- Secure session cookies (httpOnly, secure, sameSite)
- CSRF protection for state-changing operations

#### Implementation Tasks

**Backend Tasks:**

1. **Session Service** (`/backend/src/services/session.service.ts`)
   ```typescript
   interface FunnelSession {
     id: string;
     funnelId: string;
     state: 'viewing_landing' | 'viewing_sales' | 'purchased' | 'upsell_shown' | 'completed';
     ipAddress: string;
     userAgent: string;
     createdAt: Date;
     expiresAt: Date;
     metadata: {
       landingPageVisitedAt?: Date;
       salesPageVisitedAt?: Date;
       purchasedAt?: Date;
       upsellShownAt?: Date;
       completedAt?: Date;
     };
   }
   
   class SessionService {
     async createSession(funnelId: string, ipAddress: string, userAgent: string): Promise<FunnelSession>
     async getSession(sessionId: string): Promise<FunnelSession | null>
     async updateSessionState(sessionId: string, newState: string): Promise<void>
     async validateSession(sessionId: string, ipAddress: string): Promise<boolean>
     async invalidateSession(sessionId: string): Promise<void>
   }
   ```

2. **Funnel Guard Middleware** (`/backend/src/middleware/funnel-guard.middleware.ts`)
   ```typescript
   interface FunnelGuardConfig {
     requiredState: string[];
     redirectTo: string;
   }
   
   function funnelGuard(config: FunnelGuardConfig) {
     return async (req, res, next) => {
       const sessionId = req.cookies.funnel_session_id;
       const session = await sessionService.getSession(sessionId);
       
       if (!session) {
         return res.redirect(config.redirectTo);
       }
       
       if (!config.requiredState.includes(session.state)) {
         return res.redirect(config.redirectTo);
       }
       
       if (!await sessionService.validateSession(sessionId, req.ip)) {
         await sessionService.invalidateSession(sessionId);
         return res.redirect(config.redirectTo);
       }
       
       req.funnelSession = session;
       next();
     };
   }
   ```

3. **Funnel Routes** (`/backend/src/routes/funnel.routes.ts`)
   ```typescript
   router.get('/:funnelSlug/landing', funnelController.serveLandingPage);
   
   router.get('/:funnelSlug/sales', 
     funnelGuard({ requiredState: ['viewing_sales', 'purchased'], redirectTo: '/landing' }),
     funnelController.serveSalesPage
   );
   
   router.get('/:funnelSlug/upsell',
     funnelGuard({ requiredState: ['purchased', 'upsell_shown'], redirectTo: '/landing' }),
     funnelController.serveUpsellPage
   );
   
   router.get('/:funnelSlug/confirmation',
     funnelGuard({ requiredState: ['completed'], redirectTo: '/landing' }),
     funnelController.serveConfirmationPage
   );
   ```

4. **Redis Configuration** (`/backend/src/config/redis.config.ts`)
   ```typescript
   import Redis from 'ioredis';
   
   const redis = new Redis({
     host: process.env.REDIS_HOST || 'localhost',
     port: parseInt(process.env.REDIS_PORT || '6379'),
     password: process.env.REDIS_PASSWORD,
     db: 0,
     retryStrategy: (times) => Math.min(times * 50, 2000)
   });
   
   export default redis;
   ```

**Frontend Tasks:**

5. **Back Button Prevention** (`/frontend/src/utils/preventBack.ts`)
   ```typescript
   export function preventBackNavigation() {
     window.history.pushState(null, '', window.location.href);
     
     window.addEventListener('popstate', () => {
       window.history.pushState(null, '', window.location.href);
       showBackBlockedModal();
     });
   }
   
   function showBackBlockedModal() {
     // Show modal: "Você não pode voltar. Aceite ou recuse a oferta."
   }
   ```

6. **Session Initialization** (`/frontend/src/hooks/useFunnelSession.ts`)
   ```typescript
   export function useFunnelSession() {
     useEffect(() => {
       // Initialize session on landing page load
       if (!getCookie('funnel_session_id')) {
         fetch('/api/session/init', { method: 'POST' })
           .then(res => res.json())
           .then(data => {
             // Session cookie set by backend
           });
       }
     }, []);
   }
   ```

**Testing Tasks:**

7. **Session Service Tests** (`/backend/tests/services/session.service.test.ts`)
   - Test session creation
   - Test session validation
   - Test session expiration
   - Test IP validation
   - Test session hijacking prevention

8. **Funnel Guard Tests** (`/backend/tests/middleware/funnel-guard.test.ts`)
   - Test unauthorized access blocking
   - Test valid access allowing
   - Test redirect behavior
   - Test session state validation

9. **E2E Funnel Flow Tests** (`/backend/tests/e2e/funnel-flow.test.ts`)
   - Test complete funnel flow
   - Test direct URL access attempts
   - Test back button blocking
   - Test session expiration handling

#### Acceptance Criteria

- [ ] Session created on landing page access
- [ ] Direct access to upsell page redirects to landing
- [ ] Back button from upsell page is blocked
- [ ] Session expires after 12 hours
- [ ] Session hijacking from different IP is prevented
- [ ] All E2E tests pass

---

### Feature 6: Custom Domain & SSL

**Priority:** HIGH  
**Complexity:** Low  
**Estimated Time:** 3 days

#### Description

Enable users to connect custom domains to their funnels with automatic SSL certificate generation and renewal.

#### Key Requirements

**Domain Configuration:**
- User interface to add custom domain
- DNS verification (CNAME record check)
- Status tracking (pending, active, failed)
- Instructions for DNS configuration

**SSL Management:**
- Automatic certificate generation via Let's Encrypt
- Certificate renewal automation
- HTTPS enforcement
- Certificate status monitoring

#### Implementation Tasks

**Backend Tasks:**

1. **Domain Service** (`/backend/src/services/domain.service.ts`)
   ```typescript
   interface CustomDomain {
     id: string;
     funnelId: string;
     domain: string;
     status: 'pending' | 'verifying' | 'active' | 'failed';
     sslStatus: 'pending' | 'active' | 'expired';
     verifiedAt?: Date;
     sslIssuedAt?: Date;
     sslExpiresAt?: Date;
   }
   
   class DomainService {
     async addDomain(funnelId: string, domain: string): Promise<CustomDomain>
     async verifyDNS(domainId: string): Promise<boolean>
     async generateSSL(domainId: string): Promise<void>
     async renewSSL(domainId: string): Promise<void>
     async getDomainStatus(domainId: string): Promise<CustomDomain>
   }
   ```

2. **DNS Verification** (`/backend/src/utils/dns-verification.ts`)
   ```typescript
   import dns from 'dns/promises';
   
   export async function verifyCNAME(domain: string, expectedTarget: string): Promise<boolean> {
     try {
       const records = await dns.resolveCname(domain);
       return records.some(record => record === expectedTarget);
     } catch (error) {
       return false;
     }
   }
   ```

3. **SSL Certificate Manager** (`/backend/src/services/ssl.service.ts`)
   ```typescript
   import { exec } from 'child_process';
   import { promisify } from 'util';
   
   const execAsync = promisify(exec);
   
   class SSLService {
     async generateCertificate(domain: string): Promise<void> {
       // Use certbot to generate Let's Encrypt certificate
       await execAsync(`certbot certonly --webroot -w /var/www/certbot -d ${domain} --non-interactive --agree-tos --email admin@funilrapido.com`);
     }
     
     async renewCertificate(domain: string): Promise<void> {
       await execAsync(`certbot renew --cert-name ${domain}`);
     }
     
     async getCertificateExpiry(domain: string): Promise<Date> {
       const { stdout } = await execAsync(`certbot certificates -d ${domain}`);
       // Parse expiry date from output
       return new Date(/* parsed date */);
     }
   }
   ```

4. **Domain Routes** (`/backend/src/routes/domain.routes.ts`)
   ```typescript
   router.post('/funnels/:funnelId/domain', domainController.addDomain);
   router.get('/funnels/:funnelId/domain', domainController.getDomain);
   router.post('/funnels/:funnelId/domain/verify', domainController.verifyDomain);
   router.delete('/funnels/:funnelId/domain', domainController.removeDomain);
   ```

5. **Nginx Configuration Generator** (`/backend/src/utils/nginx-config.ts`)
   ```typescript
   export function generateNginxConfig(domain: string, funnelId: string): string {
     return `
   server {
     listen 80;
     listen 443 ssl;
     server_name ${domain};
     
     ssl_certificate /etc/letsencrypt/live/${domain}/fullchain.pem;
     ssl_certificate_key /etc/letsencrypt/live/${domain}/privkey.pem;
     
     location / {
       proxy_pass http://localhost:3000;
       proxy_set_header Host $host;
       proxy_set_header X-Real-IP $remote_addr;
       proxy_set_header X-Funnel-ID ${funnelId};
     }
   }
     `.trim();
   }
   ```

**Frontend Tasks:**

6. **Domain Settings UI** (`/frontend/src/components/DomainSettings.tsx`)
   ```typescript
   export function DomainSettings({ funnelId }: { funnelId: string }) {
     const [domain, setDomain] = useState('');
     const [status, setStatus] = useState<'idle' | 'pending' | 'active'>('idle');
     
     const handleAddDomain = async () => {
       await fetch(`/api/funnels/${funnelId}/domain`, {
         method: 'POST',
         body: JSON.stringify({ domain }),
         headers: { 'Content-Type': 'application/json' }
       });
       
       setStatus('pending');
       pollDomainStatus();
     };
     
     // UI shows:
     // - Input for domain
     // - DNS instructions (CNAME record)
     // - Status indicator
     // - Verify button
   }
   ```

**Infrastructure Tasks:**

7. **Nginx Setup** (`/infrastructure/nginx/`)
   - Configure Nginx as reverse proxy
   - Set up certbot for Let's Encrypt
   - Create domain configuration directory
   - Set up automatic certificate renewal cron job

8. **DNS Configuration Documentation** (`/docs/custom-domain-setup.md`)
   - Step-by-step guide for users
   - Screenshots of common DNS providers
   - Troubleshooting guide

**Testing Tasks:**

9. **Domain Service Tests** (`/backend/tests/services/domain.service.test.ts`)
   - Test domain addition
   - Test DNS verification
   - Test SSL generation
   - Test domain status tracking

10. **Integration Tests** (`/backend/tests/integration/custom-domain.test.ts`)
    - Test complete domain setup flow
    - Test SSL certificate generation
    - Test domain routing

#### Acceptance Criteria

- [ ] User can add custom domain via UI
- [ ] DNS verification detects CNAME correctly
- [ ] SSL certificate generated automatically after DNS verification
- [ ] Funnel accessible via custom domain with HTTPS
- [ ] Clear instructions shown to user
- [ ] Status updates in real-time

---

## Database Schema

### Tables

**1. funnels**
```sql
CREATE TABLE funnels (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL UNIQUE,
  status VARCHAR(50) DEFAULT 'draft', -- draft, published, archived
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_funnels_user_id ON funnels(user_id);
CREATE INDEX idx_funnels_slug ON funnels(slug);
```

**2. funnel_pages**
```sql
CREATE TABLE funnel_pages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  funnel_id UUID NOT NULL REFERENCES funnels(id) ON DELETE CASCADE,
  type VARCHAR(50) NOT NULL, -- landing, sales, upsell, confirmation
  content JSONB NOT NULL, -- Page structure and components
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_funnel_pages_funnel_id ON funnel_pages(funnel_id);
CREATE INDEX idx_funnel_pages_type ON funnel_pages(type);
```

**3. custom_domains**
```sql
CREATE TABLE custom_domains (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  funnel_id UUID NOT NULL REFERENCES funnels(id) ON DELETE CASCADE,
  domain VARCHAR(255) NOT NULL UNIQUE,
  status VARCHAR(50) DEFAULT 'pending', -- pending, verifying, active, failed
  ssl_status VARCHAR(50) DEFAULT 'pending', -- pending, active, expired
  verified_at TIMESTAMP,
  ssl_issued_at TIMESTAMP,
  ssl_expires_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_custom_domains_funnel_id ON custom_domains(funnel_id);
CREATE INDEX idx_custom_domains_domain ON custom_domains(domain);
```

**4. users**
```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
```

### Redis Schema

**Session Storage:**
```
Key: session:{sessionId}
Value: {
  id: string,
  funnelId: string,
  state: string,
  ipAddress: string,
  userAgent: string,
  createdAt: string,
  metadata: object
}
TTL: 43200 seconds (12 hours)
```

---

## Development Environment Setup

### Prerequisites

- Node.js 18+
- PostgreSQL 14+
- Redis 7+
- Docker & Docker Compose (optional)

### Setup Steps

1. **Clone Repository**
   ```bash
   git clone https://github.com/mabidoli/gravity.git
   cd gravity
   ```

2. **Install Dependencies**
   ```bash
   # Backend
   cd backend
   npm install
   
   # Frontend
   cd ../frontend
   npm install
   ```

3. **Environment Configuration**
   ```bash
   # Backend .env
   cp backend/.env.example backend/.env
   
   # Configure:
   DATABASE_URL=postgresql://user:password@localhost:5432/funilrapido
   REDIS_URL=redis://localhost:6379
   SESSION_SECRET=your-secret-key
   NODE_ENV=development
   ```

4. **Database Setup**
   ```bash
   cd backend
   npm run db:migrate
   npm run db:seed
   ```

5. **Start Development Servers**
   ```bash
   # Terminal 1: Backend
   cd backend
   npm run dev
   
   # Terminal 2: Frontend
   cd frontend
   npm run dev
   
   # Terminal 3: Redis
   redis-server
   ```

### Docker Setup (Alternative)

```bash
docker-compose up -d
```

---

## CI/CD Pipeline

### GitHub Actions Workflow

**File:** `.github/workflows/ci.yml`

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:14
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      
      redis:
        image: redis:7
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install dependencies
        run: |
          cd backend && npm ci
          cd ../frontend && npm ci
      
      - name: Run backend tests
        run: cd backend && npm test
        env:
          DATABASE_URL: postgresql://postgres:postgres@localhost:5432/test
          REDIS_URL: redis://localhost:6379
      
      - name: Run frontend tests
        run: cd frontend && npm test
      
      - name: Build
        run: |
          cd backend && npm run build
          cd ../frontend && npm run build
```

---

## Deliverables

### Week 1
- [ ] Database schema implemented and migrated
- [ ] Session service with Redis integration
- [ ] Funnel guard middleware
- [ ] Basic funnel routes
- [ ] Unit tests for session service

### Week 2
- [ ] Custom domain service
- [ ] DNS verification utility
- [ ] SSL certificate management
- [ ] Domain settings UI
- [ ] E2E tests for funnel flow
- [ ] CI/CD pipeline configured

---

## Testing Checklist

### Unit Tests
- [ ] Session creation
- [ ] Session validation
- [ ] Session expiration
- [ ] IP validation
- [ ] Domain DNS verification
- [ ] SSL certificate generation

### Integration Tests
- [ ] Funnel flow with session management
- [ ] Custom domain setup
- [ ] SSL certificate issuance

### E2E Tests
- [ ] Complete funnel flow (Landing → Sales → Upsell → Confirmation)
- [ ] Direct URL access blocking
- [ ] Back button prevention
- [ ] Session expiration handling
- [ ] Custom domain access

---

## Next Phase

Proceed to **[Phase 2: No-Code Editor](./phase-2-editor.md)** after completing all deliverables and tests.

---

**Phase Owner:** Backend Team  
**Last Updated:** January 3, 2026
