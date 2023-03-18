import { defineConfig } from 'astro/config'
import svelte from '@astrojs/svelte'
import yaml from '@rollup/plugin-yaml'
import compress from 'astro-compress'
import AstroPWA from '@vite-pwa/astro'
import rehypeExternalLinks from 'rehype-external-links'

const SITE = 'https://vmail.leopard.in.ua/'

// https://astro.build/config
export default defineConfig({
  site: SITE,
  base: '/',
  integrations: [svelte(), AstroPWA({
    injectRegister: null,
    strategies: 'injectManifest',
    registerType: 'prompt',
    srcDir: 'src',
    filename: 'sw.js',
    base: '/',
    scope: '/',
    includeAssets: ['favicon.svg', 'favicon.ico', 'icon-192x192.png', 'icon-512x512.png', 'maskable_icon.png'],
    injectManifest: {
      globPatterns: ['**/*.{css,js}']
    },
    devOptions: {
      enabled: true,
      type: 'module'
    },
    manifest: {
      name: 'VMail',
      short_name: 'VMail',
      description: 'VMail - check the markup (HTML, CSS) of HTML email template compatibility with email clients',
      theme_color: '#f9fafb',
      icons: [{
        'src': '/icon-192x192.png',
        'type': 'image/png',
        'sizes': '192x192'
      }, {
        'src': '/icon-512x512.png',
        'type': 'image/png',
        'sizes': '512x512'
      }, {
        'src': '/maskable_icon.png',
        'type': 'image/png',
        'sizes': '1024x1024',
        'purpose': 'maskable'
      }]
    }
  }), compress({
    css: true,
    html: true,
    js: true,
    img: false,
    svg: false,
    exclude: ['email_sample.html']
  })],
  markdown: {
    extendDefaultPlugins: true,
    rehypePlugins: [[rehypeExternalLinks, {
      target: '_blank',
      rel: 'noopener noreferrer'
    }]]
  },
  build: {
    format: 'file'
  },
  vite: {
    plugins: [yaml()],
    build: {
      cssCodeSplit: false,
      minify: 'esbuild'
    }
  }
})
