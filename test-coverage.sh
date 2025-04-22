go test $(go list ./internal/... | grep -v /mock) -coverprofile=profile.cov
go tool cover -func=profile.cov