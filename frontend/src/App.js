import React, { useState, useEffect } from 'react';
import { Canvas } from '@react-three/fiber';
import { OrbitControls } from '@react-three/drei';
import { Box, Container, Typography, ThemeProvider, createTheme } from '@mui/material';
import NetworkGraph from './components/NetworkGraph';
import NodeCreation from './components/NodeCreation';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
const API_BASE = `${API_URL}/v1`;

// Create a dark theme
const darkTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#90caf9',
    },
    background: {
      default: '#121212',
      paper: 'rgba(30, 30, 30, 0.8)',
    },
  },
});

function App() {
  const [nodes, setNodes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [publicIP, setPublicIP] = useState(null);

  useEffect(() => {
    fetchNodes();
    fetchPublicIP();
  }, []);

  const fetchPublicIP = async () => {
    try {
      const response = await fetch(`${API_BASE}/ip`);
      if (!response.ok) {
        throw new Error('Failed to fetch IP address');
      }
      const data = await response.json();
      setPublicIP(data.ip);
    } catch (err) {
      console.error('Failed to fetch public IP:', err);
    }
  };

  const fetchNodes = async () => {
    try {
      const response = await fetch(`${API_BASE}/nodes`);
      if (!response.ok) {
        throw new Error('Failed to fetch nodes');
      }
      const data = await response.json();
      setNodes(data);
      setError(null);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateNode = async (name) => {
    if (!publicIP) {
      throw new Error('Unable to detect public IP address');
    }

    try {
      const response = await fetch(`${API_BASE}/nodes`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
          name,
          ip: publicIP 
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Failed to create node');
      }

      setNodes([...nodes, data]);
      setError(null);
    } catch (err) {
      throw err; // Propagate error to NodeCreation component
    }
  };

  if (loading) return <Typography>Loading...</Typography>;
  if (error) return <Typography>Error: {error}</Typography>;

  return (
    <ThemeProvider theme={darkTheme}>
      <Box sx={{ position: 'relative', width: '100vw', height: '100vh', overflow: 'hidden' }}>
        {/* Full screen network visualization */}
        <Canvas camera={{ position: [0, 0, 10], fov: 75 }}>
          <color attach="background" args={['#121212']} />
          <ambientLight intensity={0.3} />
          <pointLight position={[10, 10, 10]} intensity={0.5} />
          <NetworkGraph nodes={nodes} />
          <OrbitControls enableZoom={true} enablePan={true} enableRotate={true} />
        </Canvas>

        {/* Overlay controls */}
        <Box
          sx={{
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
            p: 2,
            backgroundColor: 'background.paper',
            backdropFilter: 'blur(8px)',
            zIndex: 1,
            width: '380px',
          }}
        >
          <Typography variant="h4" gutterBottom color="primary">
            IP Proximity Network Visualizer
          </Typography>
          {publicIP && (
            <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
              Your IP: {publicIP}
            </Typography>
          )}
          <NodeCreation onCreateNode={handleCreateNode} />
        </Box>
      </Box>
    </ThemeProvider>
  );
}

export default App; 