package salvia_config

//Definimos todo lo de permisos y menús

//Permisos generales según el rol.
var PermissionsByRole map[string]map[string]bool = map[string]map[string]bool{
	//Victim COntact
	"get_victim_contact": {
		"op": true,
		"ro": true,
		"et": true,
		"sv": true,
		"no": true,
	},
	"get_victim_contacts": {
		"op": true,
		"ro": true,
		"et": true,
		"sv": true,
		"no": true,
	},
	"get_victim_contacts_by_all": {
		"op": true,
		"ro": true,
		"et": true,
		"sv": true,
		"no": true,
	},
	"set_victim_contact": {
		"op": true,
		"ro": true,
		"et": true,
	},
	"put_victim_contact": {
		"op": true,
		"ro": true,
	},
	//Victim Case
	"get_victim_case": {
		"op": true,
		"do": true,
		"ro": true,
		"et": true,
		"us": true,
		"sv": true,
		"no": true,
		"fo": true,
	},
	"get_victim_cases": {
		"op": true,
		"ro": true,
		"do": true,
		"et": true,
		"us": true,
		"sv": true,
		"no": true,
		"fo": true,
	},
	"get_victim_case_by_document": {
		"op": true,
		"ro": true,
		"et": true,
		"sv": true,
		"no": true,
	},
	"get_victim_cases_by_all": {
		"op": true,
		"ro": true,
		"et": true,
		"sv": true,
		"no": true,
	},
	"set_victim_case": {
		"op": true,
		"ro": true,
		"et": true,
		"sv": true,
		"no": true,
	},
	"update_victim_case": {
		"op": true,
		"ro": true,
	},
	"report_victim_cases": {
		"op": true,
		"ro": true,
		"sv": true,
		"no": true,
	},
	"approve_victim_case": {
		"op": true,
		"ro": true,
	},
	//EntityBranch
	"get_entity_branches_by_towncode": {
		"op": true,
		"ro": true,
		"et": true,
		"ad": true,
		"sv": true,
		"no": true,
	},
	"get_entity_branches_by_towncode_and_entity": {
		"op": true,
		"ro": true,
		"do": true,
		"et": true,
		"ad": true,
		"sv": true,
		"no": true,
	},
	"get_entity_branches_with_moments": {
		"op": true,
		"ro": true,
		"et": true,
		"ad": true,
		"sv": true,
		"no": true,
	},
	"get_entity_branches": {
		"do": true,
		"ro": true,
	},
	"set_entity_branch": {
		"do": true,
	},
	"put_entity_branch": {
		"do": true,
	},
	"update_entity_branch": {
		"do": true,
	},
	//Moment
	"update_moment": {
		"et": true,
	},
	//Alert
	"get_alerts_by_victim_case": {
		"op": true,
		"ro": true,
		"et": true,
		"sv": true,
		"no": true,
	},
	"get_alerts_by_town_code": {
		"op": true,
		"ro": true,
		"et": true,
		"sv": true,
		"no": true,
	},
	"get_alerts_in_danger": {
		"op": true,
		"ro": true,
		"sv": true,
		"no": true,
	},
	"get_alerts_by_all": {
		"sv": true,
		"op": true,
		"ro": true,
		"et": true,
	},
	//CaseLog
	"set_case_log": {
		"op": true,
		"ro": true,
		"et": true,
		"us": true,
		"sv": true,
		"no": true,
		"do": true,
	},
	//Files
	"load_plain_files": {
		"ad": true,
	},
	//AssignOperators
	"assign_operators": {
		"sv": true,
	},
	//FollowUp
	"set_follow_up": {
		"ro": true,
	},
	//FollowUpEntry
	"update_follow_up_entry": {
		"ro": true,
	},
	"set_follow_up_entry": {
		"ro": true,
	},
	//FollowUpEntryActing
	"set_follow_up_entry_acting": {
		"do": true,
	},
	//Barrier
	"get_barrier": {
		"ro": true,
	},
	//Entity
	"get_entity": {
		"do": true,
	},
	//Feminicide
	"get_feminicide": {
		"ro": true,
		"fo": true,
	},
	"get_feminicides": {
		"ro": true,
		"fo": true,
	},
	"set_feminicide": {
		"ro": true,
		"fo": true,
	},
	//FeminicideRisk
	"get_feminicide_risk": {
		"ro": true,
		"fo": true,
	},
	"get_feminicides_risk": {
		"ro": true,
		"fo": true,
	},
	"set_feminicide_risk": {
		"ro": true,
		"fo": true,
	},

	// ── Módulo de Seguimiento (HU-027) ──────────────────────────────────────
	"get_seguimiento_detalle_caso": {"ad": true, "sv": true, "op": true},
	"get_seguimiento_detalle":      {"ad": true, "sv": true, "op": true},
	"get_seguimiento_formulario":   {"op": true},
	"get_mis_seguimientos_dia":     {"op": true},
	"get_seguimientos_area":        {"ad": true, "sv": true},
	"generate_calendario_seguimiento": {"ad": true, "sv": true, "op": true},
}

//Menús del sistema. Se filtra por el tipo de menú (ie/ default) y luego por la sección principal repesentada por el locale.
var Menu map[string]map[string]map[string][]map[string]string = map[string]map[string]map[string][]map[string]string{
	"sp": {
		"op": {
			Locale["sp"]["menu_title_victim_contact"]: {
				{
					"label":       Locale["sp"]["menu_get_victim_contacts"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
					"customClass": "fas fa-eye icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_get_victim_contacts_by_all"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
					"customClass": "fas fa-eye icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_set_victim_contact"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"] + "/" + Locale["sp"]["new"],
					"customClass": "fas fa-archive unable",
				},
			},
		},
		"ro": {
			Locale["sp"]["menu_title_victim_contact"]: {
				{
					"label":       Locale["sp"]["menu_get_victim_contacts"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
					"customClass": "fas fa-eye icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_get_victim_contacts_by_all"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
					"customClass": "fas fa-eye icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_set_victim_contact"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"] + "/" + Locale["sp"]["new"],
					"customClass": "fas fa-archive unable",
				},
			},
		},
		"et": {
			Locale["sp"]["menu_title_victim_contact"]: {
				{
					"label":       Locale["sp"]["menu_get_victim_contacts"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
					"customClass": "fas fa-eye icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_get_victim_contacts_by_all"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
					"customClass": "fas fa-eye icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_set_victim_contact"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"] + "/" + Locale["sp"]["new"],
					"customClass": "fas fa-archive unable",
				},
			},
		},
		"us": {
			Locale["sp"]["menu_title_victim_cases"]: {
				{
					"label":       Locale["sp"]["menu_get_victim_cases"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
					"customClass": "fas fa-eye icoacciones",
				},
			},
		},
	},
}

//Menús secundarios del sistema para definir las acciones de edición, eliminación, etc. Se filtra por el idioma, luego por el rol, nombre del menú y finalmente sus respectivos items.
var MenuTools map[string]map[string]map[string][]map[string]string = map[string]map[string]map[string][]map[string]string{
	"sp": {
		"op": {
			"menu_tool_get_victim_contacts": {
				{
					"label":       Locale["sp"]["menu_tool_get_victim_contact"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_tool_new_victim_case"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      Locale["sp"]["new"],
					"customClass": "fas fa-plus-circle",
				},
			},
			"menu_tool_get_victim_cases": {
				{
					"label":       Locale["sp"]["menu_tool_get_victim_case"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
			"menu_tool_get_alerts": {
				{
					"label":       Locale["sp"]["menu_tool_get_alert"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["Alert"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
			"menu_tool_get_case_logs": {
				{
					"label":       Locale["sp"]["menu_tool_get_case_logs"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["CaseLog"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
		},
		"ro": {
			"menu_tool_get_victim_contacts": {
				{
					"label":       Locale["sp"]["menu_tool_get_victim_contact"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_tool_new_victim_case"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      Locale["sp"]["new"],
					"customClass": "fas fa-plus-circle",
				},
			},
			"menu_tool_get_victim_cases": {
				{
					"label":       Locale["sp"]["menu_tool_get_victim_case"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
			"menu_tool_get_alerts": {
				{
					"label":       Locale["sp"]["menu_tool_get_alert"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["Alert"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
			"menu_tool_get_case_logs": {
				{
					"label":       Locale["sp"]["menu_tool_get_case_logs"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["CaseLog"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
		},
		"do": {
			"menu_tool_get_entity_branches": {
				{
					"label":       Locale["sp"]["menu_tool_update_entity_branch"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["EntityBranch"],
					"action":      Locale["sp"]["update"],
					"customClass": "fas fa-pencil-square-o icoacciones",
				},
			},
			"menu_tool_get_victim_cases": {
				{
					"label":       Locale["sp"]["menu_tool_get_victim_case"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
		},
		"et": {
			"menu_tool_get_victim_contacts": {
				{
					"label":       Locale["sp"]["menu_tool_get_victim_contact"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
				{
					"label":       Locale["sp"]["menu_tool_new_victim_case"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      Locale["sp"]["new"],
					"customClass": "fas fa-plus-circle",
				},
			},
			"menu_tool_get_victim_cases": {
				{
					"label":       Locale["sp"]["menu_tool_get_victim_case"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
			"menu_tool_get_alerts": {
				{
					"label":       Locale["sp"]["menu_tool_get_alert"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
		},
		"us": {
			"menu_tool_get_victim_cases": {
				{
					"label":       Locale["sp"]["menu_tool_get_victim_case"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
			"menu_tool_get_alerts": {
				{
					"label":       Locale["sp"]["menu_tool_get_alert"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
		},
		"sv": {
			"menu_tool_get_victim_contacts": {
				{
					"label":       Locale["sp"]["menu_tool_get_victim_contact"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimContact"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
			"menu_tool_get_victim_cases": {
				{
					"label":       Locale["sp"]["menu_tool_get_victim_case"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
			"menu_tool_get_alerts": {
				{
					"label":       Locale["sp"]["menu_tool_get_alert"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["Alert"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
			"menu_tool_get_case_logs": {
				{
					"label":       Locale["sp"]["menu_tool_get_case_logs"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["CaseLog"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
		},
		"no": {
			"menu_tool_get_victim_cases": {
				{
					"label":       Locale["sp"]["menu_tool_get_victim_case"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
		},
		"fo": {
			"menu_tool_get_victim_cases": {
				{
					"label":       Locale["sp"]["menu_tool_get_victim_case"],
					"path":        "/" + Locale["sp"]["salvia"] + "/" + Locale["sp"]["VictimCase"],
					"action":      "",
					"customClass": "fas fa-eye icoacciones",
				},
			},
		},
	},
}
