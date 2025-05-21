import React from 'react';
import { Sphere, Line, Text, Html } from '@react-three/drei';
import * as THREE from 'three';

function NetworkGraph({ nodes }) {
  const nodeRadius = 0.5;
  const lineColor = new THREE.Color(0x90caf9); // Light blue color
  const nodeColor = new THREE.Color(0xffffff); // White color
  const textColor = new THREE.Color(0x90caf9); // Light blue color for text
  const glowColor = new THREE.Color(0x90caf9); // Glow color

  const truncateName = (name) => {
    return name.length > 15 ? name.substring(0, 15) + '...' : name;
  };

  return (
    <group>
      {nodes.map((node) => (
        <group key={node.id}>
          {/* Glow effect */}
          <Sphere
            position={[node.position.x, node.position.y, node.position.z]}
            args={[nodeRadius * 1.2, 32, 32]}
          >
            <meshBasicMaterial
              color={glowColor}
              transparent
              opacity={0.2}
              side={THREE.BackSide}
            />
          </Sphere>
          
          {/* Main node sphere */}
          <Sphere
            position={[node.position.x, node.position.y, node.position.z]}
            args={[nodeRadius, 32, 32]}
          >
            <meshStandardMaterial
              color={nodeColor}
              emissive={nodeColor}
              emissiveIntensity={0.2}
              metalness={0.8}
              roughness={0.2}
            />
          </Sphere>

          {/* Node name label with background */}
          <Html
            position={[node.position.x, node.position.y + nodeRadius + 0.5, node.position.z]}
            center
            style={{
              background: 'rgba(18, 18, 18, 0.8)',
              padding: '4px 8px',
              borderRadius: '4px',
              border: '1px solid rgba(144, 202, 249, 0.3)',
              color: '#90caf9',
              fontSize: '12px',
              fontFamily: 'Arial, sans-serif',
              whiteSpace: 'nowrap',
              pointerEvents: 'none',
              transform: 'translate3d(-50%, -50%, 0)',
              backdropFilter: 'blur(4px)',
            }}
          >
            {truncateName(node.name)}
          </Html>

          {/* Connect nodes with lines */}
          {nodes.map((otherNode) => {
            if (node.id !== otherNode.id) {
              return (
                <Line
                  key={`${node.id}-${otherNode.id}`}
                  points={[
                    [node.position.x, node.position.y, node.position.z],
                    [otherNode.position.x, otherNode.position.y, otherNode.position.z],
                  ]}
                  color={lineColor}
                  lineWidth={1}
                  opacity={0.4}
                  transparent
                />
              );
            }
            return null;
          })}
        </group>
      ))}
    </group>
  );
}

export default NetworkGraph; 