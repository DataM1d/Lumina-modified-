import React, { useRef, useMemo } from 'react';
import type { FC } from 'react';
import { useFrame } from '@react-three/fiber';
import * as THREE from 'three';
import type { Analysis, NeuralShaderMaterial } from '../../types/analysis.js';

interface VisualizerProps {
  analysis: Analysis | null;
}

const MAX_PARTICLES = 5000;
const PRE_GEN = new Float32Array(MAX_PARTICLES * 3);
for (let i = 0; i < MAX_PARTICLES * 3; i++) {
  PRE_GEN[i] = (Math.random() - 0.5) * 10;
}

const Visualizer: FC<VisualizerProps> = ({ analysis }) => {
  const materialRef = useRef<NeuralShaderMaterial>(null!);
  const positions = useMemo(() => PRE_GEN.slice(0, MAX_PARTICLES * 3), []);

  const shaderArgs = useMemo(() => ({
    uniforms: {
      uTime: { value: 0 },
      uColor: { value: new THREE.Color("#00ccff") },
      uSpeed: { value: 1.0 }
    },
    vertexShader: `
      uniform float uTime;
      uniform float uSpeed;
      void main() {
        vec3 pos = position;
        if (uSpeed > 0.1) {
          pos.x += sin(uTime * uSpeed + position.y) * 0.2;
          pos.y += cos(uTime * uSpeed + position.x) * 0.2;
        }
        vec4 mvPosition = modelViewMatrix * vec4(pos, 1.0);
        gl_PointSize = 4.0 * (10.0 / -mvPosition.z);
        gl_Position = projectionMatrix * mvPosition;
      }
    `,
    fragmentShader: `
      uniform vec3 uColor;
      void main() {
        float strength = distance(gl_PointCoord, vec2(0.5));
        strength = 1.0 - strength;
        strength = pow(strength, 3.0);
        gl_FragColor = vec4(uColor, strength);
      }
    `
  }), []);

  useFrame((state) => {
    if (!materialRef.current) return;
    
    materialRef.current.uniforms.uTime.value = state.clock.getElapsedTime();
    
    const sentiment = analysis?.sentiment ?? 0.5;
    const visualStyle = analysis?.visual_style ?? 'calm';

    const targetColor = new THREE.Color(sentiment > 0.5 ? "#ffcc00" : "#00ccff");
    materialRef.current.uniforms.uColor.value.lerp(targetColor, 0.05);
    
    const targetSpeed = visualStyle === 'energetic' ? 2.5 : 0.6;
    materialRef.current.uniforms.uSpeed.value = THREE.MathUtils.lerp(
      materialRef.current.uniforms.uSpeed.value, 
      targetSpeed, 
      0.05
    );
  });

  return (
    <points>
      <bufferGeometry>
        <bufferAttribute
          args={[positions, 3]}
          attach="attributes-position"
          count={positions.length / 3}
          array={positions}
          itemSize={3}
        />
      </bufferGeometry>
      <shaderMaterial
        ref={materialRef}
        args={[shaderArgs]}
        transparent
        depthWrite={false}
        blending={THREE.AdditiveBlending}
      />
    </points>
  );
};

export default Visualizer;