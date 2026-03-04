# ADR 0001: Clean Architecture Dependency Rule

## Status

Accepted

## Context

The project is structured in layers:

- **`business/`** — domain models (`business/models`) and application use cases (`business/usecases`). This is the core of the system.
- **`api/`** — HTTP handlers and API-specific logic (`api/file`, `api/income`, `api/login`, `api/reminder`, `api/site`, `api/user`).
- **`pkg/`** — infrastructure concerns: database sessions (`pkg/mongo`), authentication (`pkg/auth`), configuration (`pkg/config`), file export (`pkg/file`), messaging (`pkg/slack`), and shared utilities (`pkg/utils`).
- **`repositories/`** — data access implementations.
- **`worker/`** — background job processing.

Without an explicit rule about dependency direction, business logic risks becoming coupled to infrastructure details (database drivers, HTTP frameworks, third-party services), making it harder to test, reuse, and evolve independently.

## Decision

We follow the **Dependency Rule** from Clean Architecture:

> Source code dependencies must only point inward. Nothing in an inner circle can know anything about something in an outer circle.

Concretely for this project:

1. **`business/models`** must have **zero** imports from other internal packages. It defines the domain entities shared across layers.
2. **`business/usecases`** may only import **`business/models`**. When a use case needs to interact with infrastructure (e.g. persist data, send notifications), it must define an interface (driven port) within `business/usecases` that outer layers implement.
3. **Outer layers** (`api/*`, `pkg/*`, `repositories`, `worker`) may depend on `business/*` but never the reverse.

The driven ports pattern is already in use — see `business/usecases/*_driven_ports.go` and `business/usecases/*_driving_ports.go`.

### Current package dependency flow

```mermaid
graph TD
    subgraph entryPoints [Entry Points]
        main
        scripts
        cmd_keycloak["cmd/keycloak_test"]
    end

    subgraph apiLayer [API Layer]
        api_file["api/file"]
        api_income["api/income"]
        api_login["api/login"]
        api_reminder["api/reminder"]
        api_site["api/site"]
        api_user["api/user"]
    end

    subgraph businessLayer [Business Layer — must stay clean]
        biz_models["business/models"]
        biz_usecases["business/usecases"]
    end

    subgraph infraLayer [Infrastructure Layer]
        pkg_auth["pkg/auth"]
        pkg_config["pkg/config"]
        pkg_file["pkg/file"]
        pkg_mongo["pkg/mongo"]
        pkg_slack["pkg/slack"]
        pkg_utils["pkg/utils"]
    end

    repositories
    worker

    main --> api_file
    main --> api_income
    main --> api_login
    main --> api_reminder
    main --> api_site
    main --> api_user
    main --> biz_models
    main --> pkg_config
    main --> pkg_mongo
    main --> worker

    scripts --> api_income
    scripts --> biz_models
    scripts --> pkg_config
    scripts --> pkg_mongo

    cmd_keycloak --> pkg_auth

    api_file --> api_user
    api_file --> biz_models
    api_file --> pkg_mongo
    api_file --> pkg_utils

    api_income --> api_user
    api_income --> biz_models
    api_income --> biz_usecases
    api_income --> pkg_file
    api_income --> pkg_mongo
    api_income --> pkg_utils
    api_income --> repositories

    api_login --> api_site
    api_login --> api_user
    api_login --> biz_models
    api_login --> pkg_auth
    api_login --> pkg_mongo
    api_login --> pkg_utils

    api_reminder --> api_file
    api_reminder --> api_income
    api_reminder --> api_user
    api_reminder --> biz_models
    api_reminder --> pkg_mongo
    api_reminder --> pkg_utils
    api_reminder --> worker

    api_site --> biz_models
    api_site --> pkg_mongo
    api_site --> pkg_utils

    api_user --> api_site
    api_user --> biz_models
    api_user --> pkg_mongo
    api_user --> pkg_utils

    biz_usecases --> biz_models

    repositories --> biz_models
    repositories --> biz_usecases
    repositories --> pkg_mongo

    worker --> biz_models
    worker --> pkg_slack

    pkg_config --> biz_models
    pkg_file --> biz_models
    pkg_file --> pkg_utils
    pkg_mongo --> biz_models
    pkg_mongo --> pkg_config
```

### Current compliance

| Package | Internal imports | Clean? |
|---------|-----------------|--------|
| `business/models` | *(none)* | Yes |
| `business/usecases` | `business/models` only | Yes |

Both business packages comply with the Dependency Rule today.

## Consequences

### Positive

- **Testability** — business logic can be unit-tested with simple in-memory fakes; no database or HTTP server required.
- **Portability** — use cases are framework-independent and can be reused if the delivery mechanism changes (e.g. gRPC, CLI).
- **Clarity** — the dependency diagram makes coupling visible; violations are easy to spot in code review.

### Negative

- **Indirection** — use cases that need infrastructure must define driven-port interfaces, adding a layer of abstraction.
- **Discipline** — developers must resist the convenience of importing infrastructure packages directly from business code.
