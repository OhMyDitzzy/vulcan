# Vulcan Explorer

Vulcan Explorer is the official web-based dashboard and block explorer for the Vulcan blockchain. Built with **React**, **TypeScript**, **Vite**, and **Tailwind CSS**, it provides a clean and responsive interface for interacting with Vulcan nodes.

## Features
- 🔍 Real‑time blockchain explorer
- 📦 Block list and block detail views
- 💸 Transaction explorer with decoding
- 👛 Wallet creation & management (demo mode)
- ⚙️ Miner control panel
- 📊 Dashboard with node statistics
- 🌐 Fully powered by Vulcan REST API

## Getting Started

### Install dependencies
```bash
npm install
````

### Start development server

```bash
npm run dev
```

Access the UI at: [http://localhost:5173](http://localhost:5173)

### Build for production

```bash
npm run build
```

### Preview production build

```bash
npm run preview
```

## Environment Variables

The frontend automatically connects to:

* **API**: [http://localhost:8080](http://localhost:8080)

To change API target, modify `/src/api/client.ts`.

## License

MIT — part of the Vulcan project.