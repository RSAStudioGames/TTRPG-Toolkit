package models

const (
	AttributeTypeNumeric     = "numeric"
	AttributeTypeStepDie     = "step_die"
	AttributeTypeDescriptive = "descriptive"
	AttributeTypeRankTier    = "rank_tier"
	AttributeTypeCustom      = "custom"

	ModifierDisplaySigned   = "signed"
	ModifierDisplayAbsolute = "absolute"

	NumericFormatInteger = "integer"
	NumericFormatFloat   = "float"

	CachingRuleLive      = "live"
	CachingRuleOnTrigger = "on_trigger"
)

var allowedAttributeTypes = map[string]struct{}{
	AttributeTypeNumeric:     {},
	AttributeTypeStepDie:     {},
	AttributeTypeDescriptive: {},
	AttributeTypeRankTier:    {},
	AttributeTypeCustom:      {},
}

var allowedModifierDisplays = map[string]struct{}{
	ModifierDisplaySigned:   {},
	ModifierDisplayAbsolute: {},
}

var allowedNumericFormats = map[string]struct{}{
	NumericFormatInteger: {},
	NumericFormatFloat:   {},
}

var allowedCachingRules = map[string]struct{}{
	CachingRuleLive:      {},
	CachingRuleOnTrigger: {},
}

func IsAllowedAttributeType(v string) bool {
	_, ok := allowedAttributeTypes[v]
	return ok
}

func IsAllowedModifierDisplay(v string) bool {
	_, ok := allowedModifierDisplays[v]
	return ok
}

func IsAllowedNumericFormat(v string) bool {
	_, ok := allowedNumericFormats[v]
	return ok
}

func IsAllowedCachingRule(v string) bool {
	_, ok := allowedCachingRules[v]
	return ok
}
