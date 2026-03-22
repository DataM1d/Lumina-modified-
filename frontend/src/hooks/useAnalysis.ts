import { useState } from 'react';
import type { FormEvent, ChangeEvent } from 'react';
import type { AnalysisResponse, ApiError } from '../types/analysis.js';

export const useAnalysis = () => {
  const [url, setUrl] = useState<string>('');
  const [data, setData] = useState<AnalysisResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const handleUrlChange = (e: ChangeEvent<HTMLInputElement>): void => {
    setUrl(e.target.value);
  };

  const executeAnalysis = async (e: FormEvent<HTMLFormElement>): Promise<void> => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const response = await fetch('http://localhost:8080/api/v1/process', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url }),
      });

      if (!response.ok) {
        const errData = (await response.json()) as ApiError;
        throw new Error(errData.error || 'Analysis failed');
      }

      const result = (await response.json()) as AnalysisResponse;
      setData(result);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Unknown network error');
    } finally {
      setLoading(false);
    }
  };

  return {
    url,
    data,
    loading,
    error,
    handleUrlChange,
    executeAnalysis,
  };
};