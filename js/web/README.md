# Lumo Frontend

This is a simple React frontend for the Lumo application. It connects to the Lumo backend using Connect RPC (a modern gRPC-web compatible protocol) to create and display a sample link between Lumes.

## Getting Started

### Prerequisites

- Node.js (v14 or later)
- Yarn package manager
- Lumo backend running (default: http://localhost:8080)

### Installation

1. Install dependencies:

```bash
yarn install
```

### Running the Application

Start the development server:

```bash
yarn start
```

This will run the app in development mode. Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

## Connecting to the Backend

The frontend is configured to connect to the backend at `http://localhost:8080`. If your backend is running on a different URL, you'll need to update the `baseUrl` in `src/App.tsx`:

```typescript
const transport = createConnectTransport({
  baseUrl: "http://localhost:8080", // Change this to your backend URL
});
```

## Features

- Creates and displays a sample link between Lumes
- Shows detailed information about the link, including:
  - Source and destination Lume IDs
  - Link type
  - Travel details (when applicable)
  - Notes

## Implementation Details

- The client is initialized inside the component to ensure it's properly created
- Uses the CreateLink endpoint to create a sample link instead of ListLinks
- Displays a single link instead of a list of links

## Troubleshooting

If you see an error message saying "Failed to create link", check that:

1. The backend server is running
2. The backend URL is correctly configured in `src/App.tsx`
3. CORS is properly configured on the backend to allow requests from the frontend

## Building for Production

To build the app for production:

```bash
yarn build
```

This creates optimized production files in the `build` folder that can be deployed to a web server.

## Cleanup

See `files-to-delete.md` for a list of files that can be safely deleted from the project.
