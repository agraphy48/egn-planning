/**
 * Pulse Project Workspace - Auto-Recovery Daemon
 * Memastikan Pulse Workspace (Port 5080 & Backend Port 5081) terus berjalan 24/7 tanpa henti.
 * Jika proses terhenti atau crash, daemon otomatis me-restart dalam 3 detik.
 */

const { spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

const ROOT_DIR = __dirname;
const LOGS_DIR = path.join(ROOT_DIR, 'logs');
const DAEMON_LOG = path.join(LOGS_DIR, 'daemon.log');
const OUT_LOG = path.join(LOGS_DIR, 'pulse-server.log');
const ERR_LOG = path.join(LOGS_DIR, 'pulse-error.log');

function ensureLogDir() {
  if (!fs.existsSync(LOGS_DIR)) {
    fs.mkdirSync(LOGS_DIR, { recursive: true });
  }
}

function log(msg) {
  const timestamp = new Date().toLocaleString('id-ID', { timeZone: 'Asia/Jakarta' });
  const line = `[${timestamp} WIB] [PULSE-DAEMON] ${msg}\n`;
  console.log(line.trim());
  try {
    ensureLogDir();
    fs.appendFileSync(DAEMON_LOG, line, 'utf8');
  } catch (e) {}
}

let childProcess = null;
let isStopping = false;
let restartCount = 0;

function startServer() {
  if (isStopping) return;

  ensureLogDir();
  const outStream = fs.openSync(OUT_LOG, 'a');
  const errStream = fs.openSync(ERR_LOG, 'a');

  log(`Memulai proses Pulse Workspace Production Server (Run #${restartCount + 1})...`);

  childProcess = spawn(process.execPath, [path.join(ROOT_DIR, 'production_server.js')], {
    cwd: ROOT_DIR,
    stdio: ['ignore', outStream, errStream],
    shell: false,
    env: {
      ...process.env,
      PORT: '5080',
      BACKEND_PORT: '5081'
    },
    windowsHide: true,
  });

  childProcess.on('spawn', () => {
    log(`Pulse Workspace Server aktif dengan PID: ${childProcess.pid} (Port 5080)`);
  });

  childProcess.on('exit', (code, signal) => {
    try { fs.closeSync(outStream); } catch {}
    try { fs.closeSync(errStream); } catch {}

    log(`Pulse Workspace Server terhenti (code: ${code}, signal: ${signal}).`);
    childProcess = null;

    if (!isStopping) {
      restartCount++;
      log('Merestart Pulse Workspace Server dalam 3 detik...');
      setTimeout(startServer, 3000);
    }
  });

  childProcess.on('error', (err) => {
    log(`Gagal menjalankan Pulse Workspace Server: ${err.message}`);
  });
}

function stopServer() {
  log('Menghentikan Pulse Workspace Server & Daemon...');
  isStopping = true;
  if (childProcess) {
    childProcess.kill();
  }
  process.exit(0);
}

process.on('SIGINT', stopServer);
process.on('SIGTERM', stopServer);

process.on('uncaughtException', (err) => {
  log(`Uncaught Exception in Daemon: ${err.message}\n${err.stack}`);
});

process.on('unhandledRejection', (reason) => {
  log(`Unhandled Rejection in Daemon: ${reason}`);
});

// Perpetual Heartbeat to guarantee Node.js event loop never idles out
setInterval(() => {
  if (!childProcess && !isStopping) {
    log('Heartbeat: childProcess null, starting server...');
    startServer();
  }
}, 10000);

log('====================================================');
log('Pulse Workspace Auto-Recovery Daemon Dimulai');
log(`Working Directory: ${ROOT_DIR}`);
log('Target Domain: project.kedukanaksara.com');
log('====================================================');

startServer();
