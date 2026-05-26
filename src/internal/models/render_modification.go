package models

import "time"

// RenderModification representa la tabla salvia.render_modification.
// Define cómo alterar los campos de texto visibles al usuario (question.description,
// form_section.name/description, repeater_group.name/item_name) en función del formState.
// No tiene DeletedAt: las modificaciones se reemplazan, no se eliminan lógicamente.
type RenderModification struct {
	ID               string    `gorm:"type:varchar(36);primaryKey;default:gen_random_uuid()" json:"id"`
	TargetType       string    `gorm:"type:varchar(50);not null"                             json:"targetType"`
	TargetID         string    `gorm:"type:varchar(36);not null;index"                       json:"targetId"`
	TargetField      string    `gorm:"type:varchar(50);not null"                             json:"targetField"`
	ModificationType string    `gorm:"type:varchar(10);not null"                             json:"modificationType"`
	StatePath        string    `gorm:"type:varchar(255);not null"                            json:"statePath"`
	SearchString     *string   `gorm:"type:varchar(255)"                                     json:"searchString,omitempty"`
	CreatedAt        time.Time `                                                             json:"createdAt"`
	UpdatedAt        time.Time `                                                             json:"updatedAt"`
}

func (RenderModification) TableName() string { return "salvia.render_modification" }
