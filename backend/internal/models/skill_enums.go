package models

const (
	SkillTypeBinary    = "binary"
	SkillTypeMultiTier = "multi_tier"
	SkillTypeNumeric   = "numeric"
	SkillTypeStepDie   = "step_die"
	SkillTypeRankName  = "rank_name"
)

var allowedSkillTypes = map[string]struct{}{
	SkillTypeBinary:    {},
	SkillTypeMultiTier: {},
	SkillTypeNumeric:   {},
	SkillTypeStepDie:   {},
	SkillTypeRankName:  {},
}

func IsAllowedSkillType(v string) bool {
	_, ok := allowedSkillTypes[v]
	return ok
}
