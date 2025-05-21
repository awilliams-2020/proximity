import React, { useState } from 'react';
import { Box, TextField, Button, Typography } from '@mui/material';

function NodeCreation({ onCreateNode }) {
  const [name, setName] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (name.trim()) {
      try {
        await onCreateNode(name.trim());
        setName('');
        setError('');
      } catch (err) {
        setError(err.message);
      }
    }
  };

  return (
    <Box>
      <Box
        component="form"
        onSubmit={handleSubmit}
        sx={{
          display: 'flex',
          gap: 1,
          alignItems: 'center',
        }}
      >
        <TextField
          size="small"
          value={name}
          onChange={(e) => {
            setName(e.target.value);
            setError(''); // Clear error when user types
          }}
          placeholder="Enter node name"
          variant="outlined"
          error={!!error}
          sx={{
            flex: 1,
            '& .MuiOutlinedInput-root': {
              backgroundColor: 'rgba(255, 255, 255, 0.05)',
            },
          }}
        />
        <Button
          type="submit"
          variant="contained"
          color="primary"
          size="small"
          disabled={!name.trim()}
        >
          Create Node
        </Button>
      </Box>
      {error && (
        <Typography
          variant="caption"
          color="error"
          sx={{
            mt: 1,
            display: 'block',
            color: '#f44336',
            fontSize: '0.75rem',
          }}
        >
          {error}
        </Typography>
      )}
    </Box>
  );
}

export default NodeCreation; 