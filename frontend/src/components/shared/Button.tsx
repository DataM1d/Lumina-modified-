import type { FC, ButtonHTMLAttributes } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  loading?: boolean;
}

const Button: FC<ButtonProps> = ({ children, loading, className = '', ...props }) => {
  return (
    <button
      {...props}
      className={`bg-white text-slate-950 font-black px-10 py-3 rounded-full transition-all 
      hover:bg-cyan-400 hover:scale-105 active:scale-95 text-[10px] tracking-widest 
      shadow-[0_0_20px_rgba(255,255,255,0.1)] uppercase disabled:opacity-50 ${className}`}
    >
      {loading ? "SCANNING" : children}
    </button>
  );
};

export default Button;