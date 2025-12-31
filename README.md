Lumina Core: Spatial AI Neural Visualizer

Lumina Core is a high performance, 3D data visualization engine that transforms live web articles into interactive, AI analyzed neural environments. By combining a Go based microservice with a React Three Fiber (WebGL) frontend, Lumina provides a spatial HUD for digesting information through sentiment driven aesthetics

🚀 How It Works: The Data Pipeline
The Lumina engine operates through a four stage pipeline designed for low latency and high visual fidelity:

1. Ingestion: The user provides a URL. The Go backend utilizes Colly to perform a "deep scrape". Bypassing meta tags to extract the core semantic content of the article.

2. Orchestration: The extracted text is dispatched to Google Gemini 2.5 Flash. The AI is constrained via system prompting to return a strict JSON schema containing a headline, summary, sentiment score(0.0 to 1.0) and a visual mood.

3. Serialization: The Go backend sanitizes the AI response, stripping markdown artifacts and validating the JSON structure before delivering a unified payload to the frontend.

4. Visualization: The React frontend maps these values to GLSL Uniforms. The sentiment score dictates the color interpolation (lerp) from Cyan to Gold, while the visual style modulates the frequency and amplitude of the particle vertex shader

🛠 Challenges & Solutions

1. GPU Overdraw & Battery Optimization
Challenge: Rendering 5 000 active particles with real time physics was consuming high GPU cycles, draining laptop batteries rapidly.

Solution: Vertex shader branching -> particles skip expensive trigonometric calculations if the uSpeed uniform is below a threshold.

  Adaptive Fidelity: Integrated the prefers reduced motion media query. When active, the system slashes particle count by 80% and freezes shader physics, significantly reducing power draw.

  Resolution Capping: Implemented a dynamic Device Pixel Ratio(DPR) cap. On high density displays, the engine caps rendering at 1.5x instead of 3x. saving 50% of GPU fragment shader invocations.

2. Pure Functional 3D Rendering 
Challenge: Generating random particle positions inside the React render cycle caused jumps and memory leaks during state updates.

Solution: Moved particle generation to a Stable Pre-allocated Buffer outside the component. Using useMemo to slice this buffer ensures the 3D scene remains pure, predictable, and lightning fast.

3. Spatial UI Accessibility
Challenge: Standard HTML overlays felt flat,  but placing the UI inside the 3D Canvas made it unselectable and caused it to clip through the particles.

Solution: Holographic occlusion. Utilized the occlude property from @react-three/drei. This allows particles to fly behind the UI card while keeping text interactive.
  
  Perspective Managment: Used distanceFactor={10} to maintain a consistent UI scale relative to the Z-axis camera position

4. Bypassing Scraper Restrictions
Challenge: Major news outlets often return 403 forbidden errors when hit by standard scrapers. 

Solution: Hardened the scraper with Browser Mimicry Headers (User - Agent rotation, Accept - Language, and Referrer spoofing) to ensure content extraction success across diverse domains.

🎮 Navigation & ControlsLeft Click: 
Rotate the Neural Core.Right Click: Pan the Spatial HUD.Scroll: Zoom into the data stream ($Z$-axis depth:  5  to 25 units).

📦 Tech Stack

Frontend React 18 + Vite
3D Engine Three.js (React Three Fiber)
Backend Go (Gin Gonic)
AI Google Gemini 2.5 Flash
Scraper Colly v2
Post Processing Bloom, Vignette (EffectComposer)

🧪 Testing the App
To test the engine try using these URLs in your local environment:
https://science.nasa.gov/mission/mars-2020-perseverance/
https://en.wikipedia.org/wiki/Semiconductor


💡 Local Setup
Backend: Create a .env file in the root, add your GEMINI_API_KEY, and run:

go run ./cmd/api
Frontend: Navigate to your frontend directory and run:

npm install
npm run dev