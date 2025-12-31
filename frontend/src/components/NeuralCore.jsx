import { useRef, useMemo } from 'react'
import { useFrame } from '@react-three/fiber'
import * as THREE from 'three'

const MAX_PARTICLES = 5000
const PRE_GEN = new Float32Array(MAX_PARTICLES * 3)
for (let i = 0; i < MAX_PARTICLES * 3; i++) {
  PRE_GEN[i] = (Math.random() - 0.5) * 10
}

export default function NeuralCore({ sentiment = 0.5, visualStyle = 'calm', reducedMotion = false }) {
  const materialRef = useRef()
  const currentCount = reducedMotion ? 1000 : MAX_PARTICLES

  const positions = useMemo(() => PRE_GEN.slice(0, currentCount * 3), [currentCount])

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
  }), [])

  useFrame((state) => {
    if (!materialRef.current) return
    materialRef.current.uniforms.uTime.value = state.clock.getElapsedTime()
    
    const targetColor = new THREE.Color(sentiment > 0.5 ? "#ffcc00" : "#00ccff")
    materialRef.current.uniforms.uColor.value.lerp(targetColor, 0.05)
    
    const targetSpeed = reducedMotion ? 0.0 : (visualStyle === 'energetic' ? 2.5 : 0.6)
    materialRef.current.uniforms.uSpeed.value = THREE.MathUtils.lerp(
      materialRef.current.uniforms.uSpeed.value, 
      targetSpeed, 
      0.05
    )
  })

  return (
    <points>
      <bufferGeometry>
        <bufferAttribute
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
  )
}