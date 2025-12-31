import { Html } from '@react-three/drei'
import './AnalysisUI.css'

export default function AnalysisUI({ data, url, setUrl, onVisualize, loading }) {
  return (
    <Html
      transform
      distanceFactor={10}
      position={[0, 0, 0]}
      occlude="blending"
      className="ui-html-wrapper"
    >
      <article className="glass-card-3d">
        <header className="card-header">
          <h1 className="title-3d">{data.headline}</h1>
          <div className="title-underline"></div>
        </header>
        
        <div className="summary-scroll-3d">
          <p className="summary-text">{data.summary}</p>
        </div>
        
        <div className="input-row-3d">
          <input 
            type="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="INPUT SOURCE URL..."
            className="input-3d"
          />
          <button 
            onClick={onVisualize}
            disabled={loading}
            className="button-3d"
          >
            {loading ? "SCANNING" : "VISUALIZE"}
          </button>
        </div>
        <div className="holo-border"></div>
      </article>
    </Html>
  )
}