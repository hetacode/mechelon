package svvdb

// ServiceEntity struct
type ServiceEntity struct {
	ProjectName string           `bson:"project_name"`
	ServiceName string           `bson:"service_name"`
	UpdatedAt   int64            `bson:"updated_at"` // same as Instance
	Instances   []InstanceEntity `bson:"instances"`
}

// InstanceEntity struct
type InstanceEntity struct {
	Name      string `bsos:"name"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
	Status    string `bson:"status"`
}
