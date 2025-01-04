package marker

type Auditable interface {
	HasAuditModel() bool
}
