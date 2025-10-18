# OPENAPI v1
[/api/v1/api.swagger.yaml](./api/v1/api.swagger.yaml)

# TODO

## Database & Architecture
- [x] Design OpenAPI specification
- [ ] Design database schema (members, boards, lists, cards, board_members, starred_boards)
- [ ] Add unique index for `name_board_unique` field in boards table
- [ ] Create database migrations
- [ ] Set up database connection pool and configuration

## Project Structure
- [ ] Set up Go project structure (handlers, services, repositories, models)
- [ ] Configure dependencies
- [ ] Set up configuration management (YAML/ENV)
- [ ] Add logging middleware and error handling
- [ ] Set up CORS and security middleware

## Authentication & Authorization
- [ ] Implement JWT token generation and validation
- [ ] Create access token and refresh token logic
- [ ] Implement user registration endpoint
- [ ] Implement user login endpoint
- [ ] Implement token refresh endpoint
- [ ] Add authentication middleware for protected routes
- [ ] Implement password hashing (bcrypt)

## Members API
- [ ] Implement GET /members/me (get current user info)
- [ ] Implement PUT /members/me (update user profile)
- [ ] Implement POST /members/boards/{nameBoardUnique}/join (join board by unique name)
- [ ] Implement POST /members/boards/{idBoard}/star (star board)
- [ ] Implement DELETE /members/boards/{idBoard}/star (unstar board)
- [ ] Add validation for username/email uniqueness

## Boards API
- [ ] Implement GET /boards (list user's boards with filters)
- [ ] Implement POST /boards (create board with unique name validation)
- [ ] Implement GET /boards/{idBoard} (get board with lists and members)
- [ ] Implement PUT /boards/{idBoard} (update board, including name_board_unique)
- [ ] Implement DELETE /boards/{idBoard} (delete board and cascade delete)
- [ ] Implement DELETE /boards/{idBoard}/members/{idMember} (remove member/leave board)
- [ ] Add board password validation logic
- [ ] Validate name_board_unique format (lowercase, numbers, hyphens only)
- [ ] Ensure name_board_unique uniqueness across all boards

## Lists API
- [ ] Implement POST /lists (create list with fractional indexing)
- [ ] Implement GET /lists/{idList} (get list with cards)
- [ ] Implement PUT /lists/{idList} (update list name, position, archived status)
- [ ] Implement DELETE /lists/{idList} (delete list and cascade delete)
- [ ] Add fractional indexing logic for list positioning

## Cards API
- [ ] Implement POST /cards (create card in list)
- [ ] Implement GET /cards/{idCard} (get card details)
- [ ] Implement PUT /cards/{idCard} (update card, move between lists)
- [ ] Implement DELETE /cards/{idCard} (delete card)
- [ ] Add fractional indexing logic for card positioning

## Business Logic & Validation
- [ ] Implement board membership check (access control)
- [ ] Implement board ownership check (admin operations)
- [ ] Add validation for already joined boards (409 conflict)
- [ ] Add pagination support (limit/offset)
- [ ] Implement starred boards filtering
- [ ] Add archived lists/cards filtering logic

## Testing
- [ ] Write unit tests for services
- [ ] Write integration tests for API endpoints
- [ ] Write tests for authentication middleware
- [ ] Write tests for database repositories
- [ ] Add test coverage reporting

## Deployment & DevOps
- [ ] Create Dockerfile
- [ ] Set up CI/CD pipeline
- [ ] Configure production database
- [ ] Set up environment-specific configs
- [ ] Add health check endpoint