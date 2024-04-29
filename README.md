# go_ws_game
This project aims to build a WebSocket-based game service using Golang, designed to run on Google Cloud Run. The goal is to create a fully functional game accessible through a simple API.

# Features
- Multiplayer Gameplay: Players can compete against each other in a fun and engaging game.
- Real-time Interaction: The WebSocket connection enables real-time updates and interactions between players.
- Leaderboards: Track player progress and rankings.
- API-driven: All game logic and interactions are managed through a well-defined API.
# Project Structure
## Backend
- Game Logic: Implements the core game mechanics, including question generation, answer validation, and score calculation.
- WebSocket Server: Handles WebSocket connections, manages player interactions, and broadcasts game updates.
- API Endpoints: Provides endpoints for player login, question retrieval, answer submission, leaderboard access, and other - game-related operations.
- Data Storage: Utilizes a database (e.g., Cloud SQL) to store game data, player information, and leaderboard rankings.
## Frontend
- User Interface: Provides a user-friendly interface for players to interact with the game.
- WebSocket Client: Connects to the WebSocket server and handles real-time communication.
- API Integration: Communicates with the backend API to retrieve game data and submit player actions.
## Testing
- Unit Tests: Thoroughly test individual components and functions.
- Integration Tests: Verify the interaction between different components and the overall functionality of the game.
- End-to-End Tests: Simulate real-world scenarios to ensure the game works as expected.
## Deployment
- Google Cloud Run: Deploy the game service to Google Cloud Run for scalability and ease of management.
- GitHub Actions: Automate the build and deployment process using GitHub Actions.
# Solution Overview
This project utilizes a WebSocket-based architecture to provide real-time gameplay. The backend API handles all game logic and data management, while the frontend provides a user interface for players to interact with the game.

## API Endpoints
- Login: `POST /login?username=apple` - Allows players to log in with a username.
- Get Quiz: `GET /quiz` - Retrieves a new quiz question.
- Answer Quiz: `POST /quiz` - Submits an answer to a quiz question.
- Add Quiz: `PUT /quiz` - Adds a new quiz question.
- Reset Quiz: `GET /quiz/init` - Resets the quiz to its initial state.
- List Quiz: `GET /quiz/list` - Lists all available quiz questions.
- Easter Egg: `GET /easter_egg` - Checks for easter eggs.
- Health Check: `GET /health_check` - Checks the health of the service.
- Report: `POST /report` - Allows players to report issues or concerns.
## Data Storage
- Quiz Questions: Store quiz questions, answers, and answer ratios.
- Leaderboards: Store player scores and rankings.
## Workflow
- Player Login: Players enter their username and log in using the `/login` endpoint.
- Quiz Retrieval: Players request a new quiz question using the `/quiz` endpoint.
- Answer Submission: Players submit their answers using the `/quiz` endpoint.
- Result Display: The backend validates the answer and displays the result to the player.
- Leaderboard Access: Players can view the leaderboard using the `/leaderboard` endpoint.
## Deployment Considerations
- Google Cloud Run: Deploy the game service to Google Cloud Run for automatic scaling and high availability.
- Docker: Containerize the game service using Docker for portability and consistent execution.
- Persistent Storage: Use Cloud SQL or other persistent storage solutions to store game data.
- SSL/TLS: Secure the game service using SSL/TLS certificates.
## Testing
- Unit Tests: Test individual functions and components.
- Integration Tests: Test the interaction between different components.
- End-to-End Tests: Test the entire game flow from start to finish.
## Next Steps
- Implement Game Logic: Develop the core game mechanics and logic.
- Build Frontend: Create a user interface for players to interact with the game.
- Deploy to Cloud Run: Deploy the game service to Google Cloud Run.
- Test Thoroughly: Perform comprehensive testing to ensure the game works as expected.

This README provides a comprehensive overview of the `go_ws_game` project. It outlines the project's goals, features, architecture, and deployment considerations. By following the steps outlined in this document, you can successfully build and deploy a fun and engaging WebSocket-based game service on Google Cloud Run.