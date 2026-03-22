Lumina AI Analysis Engine:
Lumina is a high performance fullstack application that transforms raw web content into structured emotional and thematic data. It leverages a Go based concurrent backend, PostgreSQL for persistence, and a Three.js React frontend to visualize sentiment through real time mathematical shaders.

Technical Architecture:
The system is built on a decoupled architecture to ensure scalability and maintainability.

Backend Strategy:
The Go server uses a repository pattern to manage data flow. Incoming requests trigger a scraper that sanitizes HTML content, which is then processed by the Gemini 2.5 Flash model. To optimize costs and performance, the system implements a caching layer that checks PostgreSQL for existing analyses before calling the AI API.

Frontend Engineering:
The React frontend uses custom hooks to manage asynchronous state. Sentiment data from the API is injected into a WebGL context, where custom GLSL shaders respond to variables like emotional intensity and visual style.

Core Features
1. Automated content extraction from any valid URL.

2. Sentiment analysis and thematic summary generation.

3. PostgreSQL caching with a 24 hour TTL logic to minimize API latency.

4. Dynamic 3D visualization using Three.js and custom shader materials.

5. RESTful API with rate limiting and graceful shutdown procedures.

Infrastructure and Tools
1. Language: Go 1.22

2. Database: PostgreSQL 16

3. AI Model: Google Gemini 2.5 Flash

4. Frontend: React, TypeScript, Three.js

5. Routing: Gin Gonic

6. Environment Management: Dotenv

Technical Challenges and Solutions:
API Rate Limiting and Resiliency ->
To handle the constraints of LLM usage, the backend implements an exponential backoff strategy. This ensures that transient network issues or rate limit triggers do not result in failed user requests.

Data Consistency ->
The PostgreSQL schema uses unique constraints and B-tree indexing on URL columns. This ensures $O(1)$ lookup times for the caching layer, allowing the system to serve repeat requests in milliseconds.

Setup and Installation
1. Clone the repository and install Go dependencies.

2. Configure the .env file with your GEMINI_API_KEY and DATABASE_URL.

3. Run the PostgreSQL migration script located in the scripts directory.

4. Start the backend using go run main.go.

5. Navigate to the frontend directory and run npm install followed by npm run dev.

Future Roadmap
1. Implementation of a history dashboard for authenticated users.

2. Support for PDF and local file uploads.

3. Enhanced shader physics for more complex emotional mapping.