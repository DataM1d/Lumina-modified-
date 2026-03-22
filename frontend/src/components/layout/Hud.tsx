import type { FC, FormEvent, ChangeEvent } from 'react';
import { Html } from '@react-three/drei';
import type { AnalysisResponse } from '@/types/analysis';
import Button from '@/components/shared/Button';
import Input from '@/components/shared/Input';

interface HudProps {
  url: string;
  loading: boolean;
  error: string | null;
  data: AnalysisResponse | null;
  onUrlChange: (e: ChangeEvent<HTMLInputElement>) => void;
  onSubmit: (e: FormEvent<HTMLFormElement>) => void;
}

const Hud: FC<HudProps> = ({ url, loading, error, data, onUrlChange, onSubmit }) => {
  return (
    <Html
      transform
      distanceFactor={10}
      position={[0, 0, 1.5]}
      occlude="blending"
      className="perspective-[1200px]"
    >
      <article className="relative w-[550px] p-12 rounded-[3.5rem] text-white border border-white/10 backdrop-blur-xl shadow-[0_0_100px_rgba(0,0,0,0.6),inset_0_0_20px_rgba(255,255,255,0.03)] bg-[radial-gradient(circle_at_top_left,rgba(255,255,255,0.05),transparent)] preserve-3d">
        
        <header className="mb-2">
          <h1 className="text-5xl font-black italic tracking-tighter uppercase drop-shadow-[0_0_30px_rgba(255,255,255,0.4)]">
            {data?.analysis.headline ?? "LUMINA.AI"}
          </h1>
          <div className="h-0.5 w-24 bg-cyan-400 mb-8 shadow-[0_0_10px_rgba(34,211,238,0.8)]" />
        </header>
        
        <div className="max-h-40 overflow-y-auto mb-10 pr-4 text-slate-300 scrollbar-none">
          <p className="text-lg font-mono leading-relaxed italic">
            {data?.analysis.summary ?? "SYSTEM READY. INPUT SOURCE URL FOR NEURAL MAPPING."}
          </p>
        </div>
        
        <form onSubmit={onSubmit} className="flex items-center gap-4 bg-white/5 p-1 rounded-full border border-white/10">
          <Input 
            type="url"
            value={url}
            onChange={onUrlChange}
            placeholder="INPUT SOURCE URL..."
            required
          />
          <Button 
            type="submit"
            disabled={loading}
            loading={loading}
          >
            VISUALIZE
          </Button>
        </form>

        {error && (
          <p className="absolute bottom-4 left-12 text-[#ff4d4d] font-mono text-[10px] uppercase tracking-widest animate-pulse">
            Error // {error}
          </p>
        )}

        <div className="absolute -inset-0.5 rounded-[3.5rem] -z-10 opacity-20 overflow-hidden pointer-events-none">
           <div className="w-full h-full bg-[linear-gradient(90deg,transparent,rgba(255,255,255,0.5),transparent)] animate-[scan-line_3s_linear_infinite]" />
        </div>
      </article>
    </Html>
  );
};

export default Hud;