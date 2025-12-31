import * as THREE from 'three'
import { Canvas } from '@react-three/fiber'
import { OrbitControls, Float, Stars } from '@react-three/drei'
import { EffectComposer, Bloom, Vignette } from '@react-three/postprocessing'
import NeuralCore from './NeuralCore'
import AnalysisUI from './AnalysisUI'
import './VisualizerScene.css'

export default function VisualizerScene({ data, url, setUrl, onVisualize, loading, reducedMotion }) {
  return (
    <div className="canvas-container w-full h-full">
      <Canvas camera={{ position: [0, 0, 20], fov: 40 }} dpr={[1, 2]}>
        <color attach="background" args={['#010208']} />
        
        <Stars radius={100} depth={50} count={5000} factor={4} saturation={0} fade speed={1} />
        
        <Float speed={2} rotationIntensity={0.4} floatIntensity={0.4}>
          <NeuralCore 
            sentiment={data.sentiment} 
            visualStyle={data.visual_style} 
            reducedMotion={reducedMotion} 
          />
          
          <group position={[0, 0, 1.5]}>
            <AnalysisUI 
              data={data} 
              url={url} 
              setUrl={setUrl} 
              onVisualize={onVisualize} 
              loading={loading} 
            />
          </group>
        </Float>

        <EffectComposer>
          <Bloom intensity={1.5} luminanceThreshold={0.1} mipmapBlur />
          <Vignette eskil={false} offset={0.1} darkness={1.2} />
        </EffectComposer>

        <OrbitControls 
          enableZoom={true} 
          enablePan={true}
          minDistance={5}
          maxDistance={25}
          autoRotate={!reducedMotion}
          autoRotateSpeed={0.5}
          mouseButtons={{
            LEFT: THREE.MOUSE.ROTATE,
            MIDDLE: THREE.MOUSE.DOLLY,
            RIGHT: THREE.MOUSE.PAN
          }}
        />
      </Canvas>
    </div>
  )
}