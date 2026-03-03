/**
 * esbuild configuration for TowCommand Lambda functions
 * Bundles all Lambda handlers for deployment
 * Pattern adapted from gutguard-ai esbuild config
 */

import * as esbuild from 'esbuild';
import { readdirSync, statSync, writeFileSync, mkdirSync, existsSync } from 'fs';
import { join, dirname } from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

/**
 * Find all Lambda handler entry points in services/
 * Looks for index.ts or handler.ts files in src/ directories
 */
function findEntryPoints() {
  const entries = [];
  const servicesDir = join(__dirname, 'services');

  if (!existsSync(servicesDir)) {
    console.warn('No services/ directory found');
    return entries;
  }

  const services = readdirSync(servicesDir);

  for (const service of services) {
    const serviceDir = join(servicesDir, service);
    if (!statSync(serviceDir).isDirectory()) continue;

    const srcDir = join(serviceDir, 'src');
    if (!existsSync(srcDir)) continue;

    // Check for handler.ts or index.ts
    for (const entryFile of ['handler.ts', 'index.ts']) {
      const entryPath = join(srcDir, entryFile);
      if (existsSync(entryPath)) {
        entries.push({ name: service, path: entryPath });
        break;
      }
    }

    // Also check for subdirectory handlers (e.g., services/api-gateway/src/handlers/)
    const handlersDir = join(srcDir, 'handlers');
    if (existsSync(handlersDir) && statSync(handlersDir).isDirectory()) {
      const handlers = readdirSync(handlersDir);
      for (const handler of handlers) {
        const handlerPath = join(handlersDir, handler);
        if (statSync(handlerPath).isDirectory()) {
          const indexPath = join(handlerPath, 'index.ts');
          if (existsSync(indexPath)) {
            entries.push({
              name: `${service}-${handler}`,
              path: indexPath,
            });
          }
        }
      }
    }
  }

  return entries;
}

// Common build options
const commonOptions = {
  bundle: true,
  platform: 'node',
  target: 'node20',
  format: 'esm',
  sourcemap: process.env.NODE_ENV !== 'production',
  minify: process.env.NODE_ENV === 'production',
  treeShaking: true,

  // External packages (AWS SDK v3 is included in Lambda runtime)
  external: [
    '@aws-sdk/*',
    'aws-sdk',
  ],

  // Banner for ESM compatibility
  banner: {
    js: `
import { createRequire } from 'module';
import { fileURLToPath } from 'url';
import { dirname } from 'path';
const require = createRequire(import.meta.url);
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
    `.trim(),
  },

  // Resolve extensions
  resolveExtensions: ['.ts', '.js', '.mjs', '.json'],

  // Log level
  logLevel: 'info',
};

// Watch mode for development
const watchMode = process.argv.includes('--watch');

async function build() {
  try {
    const entries = findEntryPoints();

    if (entries.length === 0) {
      console.warn('No entry points found. Ensure services have handler.ts or index.ts in src/');
      return;
    }

    console.log(`Found ${entries.length} Lambda functions to build:`);
    entries.forEach(e => console.log(`  - ${e.name}`));

    for (const entry of entries) {
      const outdir = join(__dirname, 'dist', entry.name);

      // Ensure output directory exists
      mkdirSync(outdir, { recursive: true });

      // Write package.json for ESM support
      writeFileSync(
        join(outdir, 'package.json'),
        JSON.stringify({ type: 'module' }, null, 2)
      );

      const buildOptions = {
        ...commonOptions,
        entryPoints: [entry.path],
        outfile: join(outdir, 'index.js'),
      };

      if (watchMode) {
        const context = await esbuild.context(buildOptions);
        await context.watch();
        console.log(`Watching ${entry.name}...`);
      } else {
        await esbuild.build(buildOptions);
        console.log(`Built ${entry.name}`);
      }
    }

    console.log('Build complete!');
  } catch (error) {
    console.error('Build failed:', error);
    process.exit(1);
  }
}

build();
