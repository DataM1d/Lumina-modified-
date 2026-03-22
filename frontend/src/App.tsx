import type { FC } from 'react';
import { Canvas } from '@react-three/fiber';
import { OrbitControls, Stars, Float } from '@react-three/drei';
import { EffectComposer, Bloom, Vignette } from '@react-three/postprocessing';

import { useAnalysis } from './hooks/useAnalysis.js';
import Hud from './components/layout/Hud.js';
import Visualizer from './components/scene/Visualizer.js';

const App: FC = () => {
  const { url, data, loading, error, handleUrlChange, executeAnalysis } = useAnalysis();

  return (
    <main className="w-screen h-screen bg-[#010208] overflow-hidden">
      <Canvas 
        camera={{ position: [0, 0, 20], fov: 40 }} 
        dpr={[1, 2]}
        gl={{ antialias: true }}
      >
        <color attach="background" args={['#010208']} />
        
        <Stars 
          radius={100} 
          depth={50} 
          count={5000} 
          factor={4} 
          saturation={0} 
          fade 
          speed={1} 
        />
        
        <Float speed={2} rotationIntensity={0.4} floatIntensity={0.4}>
          <Visualizer analysis={data?.analysis ?? null} />
          
          <Hud 
            url={url} 
            loading={loading} 
            error={error} 
            data={data} 
            onUrlChange={handleUrlChange} 
            onSubmit={executeAnalysis} 
          />
        </Float>

        <EffectComposer>
          <Bloom intensity={1.5} luminanceThreshold={0.1} mipmapBlur />
          <Vignette offset={0.1} darkness={1.2} />
        </EffectComposer>

        <OrbitControls 
          enableZoom={true} 
          autoRotate 
          autoRotateSpeed={0.5} 
          maxDistance={30}
          minDistance={10}
        />
      </Canvas>
    </main>
  );
};

export default App;