import * as THREE from 'three';

export interface Analysis {
    headline: string;
    summary: string;
    sentiment: number;
    visual_style: 'energetic' | 'calm' | 'dramatic';
}

export interface AnalysisResponse {
    source: string;
    analysis: Analysis;
    meta: {
        timestamp: number;
        version: string;
    };
}

export interface VisualizerState {
    isActive: boolean;
    intensity: number;
    colorTheme: string;
}

export interface ApiError {
    error: string;
}

export interface NeuralShaderMaterial extends THREE.ShaderMaterial {
  uniforms: {
    uTime: { value: number };
    uColor: { value: THREE.Color };
    uSpeed: { value: number };
  };
}