package redtape

import "log"

// AuditLevel specifies which operations are logged through an Auditor.
type AuditLevel int

const (
	// AuditNone disables audit output.
	AuditNone AuditLevel = iota
	// AuditDeny audits only deny effects.
	AuditDeny
	// AuditAllow audits allow effects and lower.
	AuditAllow
	// AuditRequest audits requests and lower.
	AuditRequest
	// AuditAll audits everything.
	AuditAll
)

// Auditor interface allows logging requests and results of policy operations.
type Auditor interface {
	LogRequest(req *Request)
	LogPolicyEffect(req *Request, effect PolicyEffect)
}

// NewConsoleAuditor returns an Auditor that prints the audit log to stdout.
func NewConsoleAuditor(lvl AuditLevel) Auditor {
	return &consoleAuditor{
		lvl: lvl,
	}
}

type consoleAuditor struct {
	lvl AuditLevel
}

const (
	logfmt   = "%s action=%s resource=%s role=%s scope=%s\n"
	logReq   = "[AUDIT_REQ]:"
	logAllow = "[AUDIT_ALLOW]:"
	logDeny  = "[AUDIT_DENY]:"
)

// LogRequest prints the request to console if AuditLevel is at or above AuditRequest.
func (a *consoleAuditor) LogRequest(req *Request) {
	if a.lvl >= AuditRequest {
		log.Printf(logfmt, logReq, req.Action, req.Resource, req.Role, req.Scope)
	}
}

// LogRequest prints the effect of a request to console if AuditLevel is at or above the PolicyEffect.
func (a *consoleAuditor) LogPolicyEffect(req *Request, effect PolicyEffect) {
	switch {
	case effect == PolicyEffectDeny && a.lvl >= AuditDeny:
		log.Printf(logfmt, logDeny, req.Action, req.Resource, req.Role, req.Scope)
	case effect == PolicyEffectAllow && a.lvl >= AuditAllow:
		log.Printf(logfmt, logAllow, req.Action, req.Resource, req.Role, req.Scope)
	}
}
