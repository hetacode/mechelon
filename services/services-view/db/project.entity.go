package svvdb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectEntity struct
type ProjectEntity struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Services []ServiceEntity    `bson:"services"`
}

// ServiceEntity struct
type ServiceEntity struct {
	Name      string           `bson:"name"`
	UpdatedAt int64            `bson:"updated_at"` // same as Instance
	Instances []InstanceEntity `bson:"instances"`
}

// InstanceEntity struct
type InstanceEntity struct {
	Name      string `bsos:"name"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
	Status    string `bson:"status"`
}
