# Tasks: Google login and multi-user data isolation

- [x] 1. `domain.User` + validation + unit tests.
- [x] 2. `UserID` added to `domain.Asset`/`domain.Operation` + validation + unit tests.
- [x] 3. `repository.UserRepository` interface; `userID` param added to `AssetRepository`/`OperationRepository`.
- [x] 4. `sqlite.UserRepository` (upsert-by-GoogleID) + integration test proving it updates, not duplicates, on re-login.
- [x] 5. Existing sqlite repos scoped by `user_id` in every query.
- [x] 6. `internal/auth` package: OAuth2 flow, ID token verification, stateless session JWT + unit tests.
- [x] 7. `middleware.RequireAuth`; CORS updated for credentials.
- [x] 8. `service.AuthService`; `userID` threaded through `AssetService`/`OperationService`/`PositionService`.
- [x] 9. `/api/v1/auth/{google/login,google/callback,logout,me}` routes; every data route wrapped in `RequireAuth`.
- [x] 10. Frontend: `AuthProvider`/`useAuth`, `AuthGate`, `UserMenu`, `credentials: "include"` on the API client.
- [x] 11. Update docs: README (Google Cloud setup steps), `docs/catalog.md`, spec 001's non-goal note.
- [x] 12. ADR 0009 for the session/OAuth design.
- [ ] 13. Manual smoke test against a real Google OAuth client (needs credentials only the project owner can create — not exercised in CI).
