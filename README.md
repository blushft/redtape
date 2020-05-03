
# Redtape [![Go Report Card](https://goreportcard.com/badge/github.com/blushft/redtape)](https://goreportcard.com/report/github.com/blushft/redtape) ![Go](https://github.com/blushft/redtape/workflows/Go/badge.svg) [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/blushft/redtape) [![license](http://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/blushft/redtape/master/LICENSE)

## A flexible policy engine for Go

Redtape is based on the excellent [ory/ladon](https://github.com/ory/ladon) package which is no longer maintained.

## Instalallation

```bash
go get github.com/blushft/redtape
```

## Usage

### Roles

Roles are the basic permission unit for Redtape. A role is identified by a simple string like `view_comments` and can contain other roles. 

```golang
myrole := redtape.NewRole("edit_comments")

myrole.AddRole(redtape.NewRole("view_comments"))
```

A role can also hold a name and description for display purposes.

```golang
myrole.Name = "Comment Editor"
myrole.Description = "This role is able to edit comments."
```

The role manager interface allows a persistence and storage mechanism for interacting with roles.

```golang
manager := redtape.NewRoleManager()
manager.Create(myrole)

role, err := manager.Get("edit_comments")
```

### Requests

A request specifies a set of values that can be processed to determine what permission to apply.

```golang
req := redtape.NewRequest("/comments", "GET", "edit_comments", "post")

fmt.Printf(
    "Request:\nResource: %s, Action: %s, Role: %s, Scope: %s",
    req.Resource, // The resource describes what is being accessed
    req.Action,  // The action describes what is being done to the resource
    req.Role, // The role indicates what role the caller holds
    req.Scope, // The scope describes a context for the resource or action
)
```

Requests also contains a context that can carry metadata into policy objects like conditions.

```golang
req := NewRequest("/comments", "GET", "edit_comments", "post",
 map[string]interface{}{
    "headers": httpRequest.Header,
})
```

If you'd like to append an existing context with this metadata, use the `NewRequestWithContext` method.



### Policies

Policies describe what permissions to apply to requests. A policy can contain a range of Resources, Actions, Roles, and Scopes used to match that policy to the request. Further, you can use conditions to express logic needed to determine a `PolicyEffect`.

Policies are built with the `PolicyOptions` struct using a functional options pattern.

```golang
policy, err := redtape.NewPolicy(
    redtape.PolicyName("allow_edit_comments"),
    redtape.PolicyDescription("Allows members of the edit_comments role to edit comments"),
    redtape.SetResources("/comments"),
    redtape.SetActions("GET", "POST", "PUT"),
    redtape.WithRole("edit_comments"),
    redtape.PolicyAllow(),
)
```

To enable efficient storage, you can also unmarshal policy options from json.

```golang
j := `
{
    "name": "allow_edit_comments",
    "description": "Allows members of the edit_comments role to edit comments.",
    "roles": [
        "edit_comments"
    ],
    "resources": [
        "/comments"
    ],
    "actions": [
        "GET",
        "POST",
        "PUT"
    ],
    "effect": "allow"
}
`
var opts redtape.PolicyOptions
err := json.Unmarshal([]byte(j), &opts)

policy := redtape.NewPolicy(redtape.SetPolicyOptions(opts))
```

### Conditions

Conditions can be applied to policies to add additional logic to the application of permissions.

TODO: Document usage API for conditions.


### PolicyManager

The policy manager interface provides basic methods to allow you to load policies from memory, a storage backend, or files. The default manager is memory backed without persistence.

```golang
manager := redtape.NewManager()

err := manager.Create(myPolicy)
```

### Enforcer

An enforcer brings together a `PolicyManager` and `Matcher` to enforce permssions on requests.

```golang
enforcer, err := redtape.NewDefaultEnforcer(manager)

err := enforcer.Enforce(myRequest)
```

The default enforcer uses the default matcher which allows resources, actions, and scopes to be matched with wildcards. 

Policies are evaluated in order to ensure matches against actions, then resources, then roles, then scopes, and finally conditions. If any matched policy evaluates to `PolicyEffect` deny, the request is actively denied. If no policy matches and the package level `DefaultPolicyEffect` is deny (the default), the request is implicitly denied.

Permission is determined by the error value returned by `Enforce()`. A `nil` error is considered permission allowed.

```golang
if err := enforcer.Enforce(req); err != nil {
    log.Println("Request Denied")
    return
}

// Do the request here
```

### Todo
- [x] RoleManager interface
- [ ] SQL backend for managers
- [ ] KV Store backend for managers
- [ ] URL backend for managers
- [ ] Improve `Condition` API
- [ ] Expand `Scope` utilities
- [ ] Improve `context.Context` interopertation
- [ ] Create middlewares for popular frameworks
- [ ] Increased test coverage
- [ ] Examples

