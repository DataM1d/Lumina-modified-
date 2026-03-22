import type { FC, InputHTMLAttributes } from 'react';

const Input: FC<InputHTMLAttributes<HTMLInputElement>> = ({ className = '', ...props }) => {
  return (
    <input
      {...props}
      className={`flex-1 bg-transparent px-6 py-4 text-white placeholder-slate-600 
      focus:outline-none text-[10px] font-mono ${className}`}
    />
  );
};

export default Input;