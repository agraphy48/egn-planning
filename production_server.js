const http = require('http');
const fs = require('fs');
const path = require('path');
const zlib = require('zlib');
const { spawn } = require('child_process');

const PORT = process.env.PORT || 5080;
const BACKEND_PORT = process.env.BACKEND_PORT || 5081;
const DIST_DIR = path.join(__dirname, 'frontend', 'dist');

// MIME types dictionary
const MIME_TYPES = {
  '.html': 'text/html; charset=utf-8',
  '.js': 'application/javascript; charset=utf-8',
  '.css': 'text/css; charset=utf-8',
  '.json': 'application/json; charset=utf-8',
  '.svg': 'image/svg+xml',
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.gif': 'image/gif',
  '.ico': 'image/x-icon',
  '.woff': 'font/woff',
  '.woff2': 'font/woff2',
  '.ttf': 'font/ttf',
  '.eot': 'application/vnd.ms-fontobject'
};

// 1. Rate Limiter Map (IP-based)
const rateLimits = new Map();
const RATE_LIMIT_WINDOW = 60 * 1000; // 1 minute
const MAX_GENERAL_REQUESTS = 180;
const MAX_LOGIN_REQUESTS = 15;

function isRateLimited(ip, isLogin = false) {
  const now = Date.now();
  let record = rateLimits.get(ip);
  if (!record || now - record.resetTime > RATE_LIMIT_WINDOW) {
    record = { count: 0, loginCount: 0, resetTime: now };
    rateLimits.set(ip, record);
  }

  record.count++;
  if (isLogin) {
    record.loginCount++;
    if (record.loginCount > MAX_LOGIN_REQUESTS) {
      return true;
    }
  }

  return record.count > MAX_GENERAL_REQUESTS;
}

// Clean up stale rate limit entries every 5 minutes
setInterval(() => {
  const now = Date.now();
  for (const [ip, record] of rateLimits.entries()) {
    if (now - record.resetTime > RATE_LIMIT_WINDOW) {
      rateLimits.delete(ip);
    }
  }
}, 5 * 60 * 1000);

// 2. Start Golang Backend Process as Managed Child Process
let backendProcess = null;

function startBackend() {
  const exePath = path.join(__dirname, 'backend', 'server.exe');
  if (!fs.existsSync(exePath)) {
    console.error(`[Server] Backend executable not found at ${exePath}`);
    return;
  }

  console.log(`[Server] Launching backend server on port ${BACKEND_PORT}...`);
  backendProcess = spawn(exePath, [], {
    cwd: path.join(__dirname, 'backend'),
    env: {
      ...process.env,
      PORT: String(BACKEND_PORT),
      DB_PATH: 'pulse.db'
    },
    stdio: ['ignore', 'pipe', 'pipe']
  });

  backendProcess.stdout.on('data', (data) => {
    const text = data.toString().trim();
    if (text) console.log(`[Go-Backend] ${text}`);
  });

  backendProcess.stderr.on('data', (data) => {
    const text = data.toString().trim();
    if (text) console.error(`[Go-Backend-Err] ${text}`);
  });

  backendProcess.on('exit', (code, signal) => {
    console.warn(`[Server] Go backend exited (code=${code}, signal=${signal}). Restarting in 3s...`);
    setTimeout(startBackend, 3000);
  });
}

startBackend();

// 3. Security Headers Setter
function setSecurityHeaders(res) {
  res.setHeader('X-Content-Type-Options', 'nosniff');
  res.setHeader('X-Frame-Options', 'SAMEORIGIN');
  res.setHeader('Referrer-Policy', 'strict-origin-when-cross-origin');
  res.setHeader('Permissions-Policy', 'camera=(), microphone=(), geolocation=()');
  res.setHeader('X-XSS-Protection', '1; mode=block');
  res.setHeader('Strict-Transport-Security', 'max-age=31536000; includeSubDomains');
  res.setHeader(
    'Content-Security-Policy',
    "default-src 'self' 'unsafe-inline' https: data:; font-src 'self' https: data:; img-src 'self' https: data:; connect-src 'self' https: http:;"
  );
  res.removeHeader('X-Powered-By');
}

// 4. HTTP Proxy to Go Backend
function proxyToBackend(req, res) {
  const options = {
    hostname: '127.0.0.1',
    port: BACKEND_PORT,
    path: req.url,
    method: req.method,
    headers: {
      ...req.headers,
      host: `127.0.0.1:${BACKEND_PORT}`,
      'x-forwarded-for': req.socket.remoteAddress,
      'x-forwarded-proto': 'https'
    }
  };

  const proxyReq = http.request(options, (proxyRes) => {
    setSecurityHeaders(res);
    res.writeHead(proxyRes.statusCode, proxyRes.headers);
    proxyRes.pipe(res);
  });

  proxyReq.on('error', (err) => {
    console.error(`[Proxy] Error communicating with backend: ${err.message}`);
    setSecurityHeaders(res);
    res.writeHead(502, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ error: 'Backend service unavailable', status: 502 }));
  });

  req.pipe(proxyReq);
}

// 5. Static File Server with Gzip & In-Memory / Stream Caching
function serveStaticFile(req, res, filePath) {
  fs.stat(filePath, (err, stats) => {
    if (err || !stats.isFile()) {
      // Fallback to index.html for Single Page Application client-side routing
      const indexHtmlPath = path.join(DIST_DIR, 'index.html');
      return serveIndexHtml(req, res, indexHtmlPath);
    }

    const ext = path.extname(filePath).toLowerCase();
    const contentType = MIME_TYPES[ext] || 'application/octet-stream';
    setSecurityHeaders(res);

    // Cache-Control headers
    if (filePath.includes(path.sep + 'assets' + path.sep)) {
      // Hashed assets: cache for 1 year immutable
      res.setHeader('Cache-Control', 'public, max-age=31536000, immutable');
    } else {
      res.setHeader('Cache-Control', 'public, max-age=3600');
    }

    // Compression check
    const acceptEncoding = req.headers['accept-encoding'] || '';
    res.setHeader('Content-Type', contentType);

    if (/\b(html|text|javascript|json|xml|svg)\b/.test(contentType) && acceptEncoding.includes('gzip')) {
      res.setHeader('Content-Encoding', 'gzip');
      res.writeHead(200);
      const rawStream = fs.createReadStream(filePath);
      rawStream.pipe(zlib.createGzip()).pipe(res);
    } else {
      res.setHeader('Content-Length', stats.size);
      res.writeHead(200);
      fs.createReadStream(filePath).pipe(res);
    }
  });
}

function serveIndexHtml(req, res, indexPath) {
  fs.readFile(indexPath, (err, data) => {
    if (err) {
      res.writeHead(404, { 'Content-Type': 'text/plain' });
      return res.end('Index file not found. Ensure frontend is built.');
    }

    setSecurityHeaders(res);
    // index.html should not be aggressively cached to allow instant updates
    res.setHeader('Cache-Control', 'no-cache, no-store, must-revalidate');
    res.setHeader('Content-Type', 'text/html; charset=utf-8');

    const acceptEncoding = req.headers['accept-encoding'] || '';
    if (acceptEncoding.includes('gzip')) {
      res.setHeader('Content-Encoding', 'gzip');
      zlib.gzip(data, (_, compressed) => {
        res.writeHead(200);
        res.end(compressed);
      });
    } else {
      res.writeHead(200);
      res.end(data);
    }
  });
}

// 6. Main HTTP Server
const server = http.createServer((req, res) => {
  const clientIP = req.headers['x-forwarded-for'] || req.socket.remoteAddress || 'unknown';
  const isLoginRoute = req.url.startsWith('/api/v1/auth/login');

  if (isRateLimited(clientIP, isLoginRoute)) {
    setSecurityHeaders(res);
    res.writeHead(429, { 'Content-Type': 'application/json' });
    return res.end(JSON.stringify({ error: 'Too many requests. Please slow down.', status: 429 }));
  }

  // Route: Health check endpoint
  if (req.url === '/health' || req.url === '/ping') {
    setSecurityHeaders(res);
    res.writeHead(200, { 'Content-Type': 'application/json' });
    return res.end(JSON.stringify({
      status: 'healthy',
      subdomain: 'project.kedukanaksara.com',
      service: 'Pulse Workspace Production Host',
      time: new Date().toISOString()
    }));
  }

  // Route: Golang Backend REST API
  if (req.url.startsWith('/api/') || req.url.startsWith('/health-backend')) {
    return proxyToBackend(req, res);
  }

  // Route: Static Frontend SPA
  let safePath = path.normalize(decodeURIComponent(req.url.split('?')[0]));
  if (safePath === '/' || safePath === '') {
    safePath = '/index.html';
  }

  const resolvedFilePath = path.join(DIST_DIR, safePath);

  // Security: Prevent directory traversal
  if (!resolvedFilePath.startsWith(DIST_DIR)) {
    setSecurityHeaders(res);
    res.writeHead(403, { 'Content-Type': 'text/plain' });
    return res.end('Access denied.');
  }

  serveStaticFile(req, res, resolvedFilePath);
});

// Listen on all local loopback & container interfaces (supports both IPv4 and IPv6)
server.listen(PORT, () => {
  console.log(`[Pulse] Production Server listening on port ${PORT} (IPv4/IPv6)`);
  console.log(`[Pulse] Subdomain Target: project.kedukanaksara.com`);
});

// Graceful Termination
function shutdown() {
  console.log('[Pulse] Shutting down production server...');
  if (backendProcess) {
    backendProcess.kill();
  }
  server.close(() => {
    process.exit(0);
  });
}

process.on('SIGINT', shutdown);
process.on('SIGTERM', shutdown);

process.on('uncaughtException', (err) => {
  console.error(`[Server UncaughtException] ${err.message}\n${err.stack}`);
});

process.on('unhandledRejection', (reason) => {
  console.error(`[Server UnhandledRejection] ${reason}`);
});
