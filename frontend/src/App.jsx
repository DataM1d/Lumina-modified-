import { useState, useEffect } from 'react'
import VisualizerScene from './components/VisualizerScene'

function App() {
  const [url, setUrl] = useState('')
  const [loading, setLoading] = useState(false)
  const [reducedMotion, setReducedMotion] = useState(false)
  const [data, setData] = useState({
    headline: 'LUMINA CORE',
    summary: 'SYSTEM READY: Awaiting spatial data stream...',
    sentiment: 0.5,
    visual_style: 'calm'
  })

  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)')
    setReducedMotion(mediaQuery.matches)
    const handler = (event) => setReducedMotion(event.matches)
    mediaQuery.addEventListener('change', handler)
    return () => mediaQuery.removeEventListener('change', handler)
  }, [])

  const handleVisualize = async () => {
    if (!url) return
    setLoading(true)
    try {
      const response = await fetch('http://localhost:8080/process', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url })
      })
      
      if (!response.ok) throw new Error('Network response was not ok')
      
      const result = await response.json()
      const aiData = typeof result.analysis === 'string' 
        ? JSON.parse(result.analysis) 
        : result.analysis
        
      setData(aiData)
    } catch (err) {
      console.error('Visualization Error:', err)
      setData(prev => ({ 
        ...prev, 
        headline: "CORE ERROR", 
        summary: "Data stream interrupted. Please verify the URL or backend status." 
      }))
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="w-full h-screen bg-[#010208]">
      <VisualizerScene 
        data={data} 
        url={url} 
        setUrl={setUrl} 
        onVisualize={handleVisualize} 
        loading={loading}
        reducedMotion={reducedMotion} 
      />
    </main>
  )
}

export default App