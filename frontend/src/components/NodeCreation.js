import React, { useState } from 'react';
import { Box, TextField, Button, Typography } from '@mui/material';

function NodeCreation({ onCreateNode }) {
  const [name, setName] = useState('');
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (name.trim()) {
      setIsSubmitting(true);
      setError('');
      try {
        await onCreateNode(name.trim());
        setName('');
      } catch (err) {
        // Check for specific error messages
        if (err.message.includes('Node name already taken')) {
          setError('This name is already taken. Please choose a different name.');
        } else if (err.message.includes('IP node already exists')) {
          setError('A node with your IP address already exists.');
        } else {
          setError(err.message || 'Failed to create node. Please try again.');
        }
      } finally {
        setIsSubmitting(false);
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
          disabled={isSubmitting}
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
          disabled={!name.trim() || isSubmitting}
        >
          {isSubmitting ? 'Creating...' : 'Create Node'}
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