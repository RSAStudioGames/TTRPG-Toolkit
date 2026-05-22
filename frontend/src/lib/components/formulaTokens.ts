export type TokenGroup = 'dice' | 'modifiers' | 'logic' | 'functions' | 'variables';

export interface FormulaToken {
	group: TokenGroup;
	label: string;
	insert: string;
	/** Tooltip; omitted when label is sufficient */
	title?: string;
}

export const DICE_TOKENS: FormulaToken[] = [
	{ group: 'dice', label: 'd4', insert: '1d4' },
	{ group: 'dice', label: 'd6', insert: '1d6' },
	{ group: 'dice', label: '2d6', insert: '2d6', title: 'Two six-sided dice' },
	{ group: 'dice', label: 'd8', insert: '1d8' },
	{ group: 'dice', label: 'd10', insert: '1d10' },
	{ group: 'dice', label: 'd12', insert: '1d12' },
	{ group: 'dice', label: 'd20', insert: '1d20' }
];

export const MODIFIER_TOKENS: FormulaToken[] = [
	{ group: 'modifiers', label: 'Keep Highest', insert: 'kh', title: 'Keep N highest dice' },
	{ group: 'modifiers', label: 'Keep Lowest', insert: 'kl', title: 'Keep N lowest dice' },
	{ group: 'modifiers', label: 'Drop Lowest', insert: 'dl', title: 'Drop N lowest dice' },
	{ group: 'modifiers', label: 'Drop Highest', insert: 'dh', title: 'Drop N highest dice' },
	{ group: 'modifiers', label: 'Reroll', insert: 'r', title: 'Reroll results ≤ N' },
	{ group: 'modifiers', label: 'Explode', insert: 'x', title: 'Explode on result ≥ N' },
	{
		group: 'modifiers',
		label: 'Compound',
		insert: 'xo',
		title: 'Compound exploding dice on ≥ N'
	},
	{ group: 'modifiers', label: 'Count Successes', insert: 't', title: 'Count results ≥ N' },
	{ group: 'modifiers', label: 'Critical Count', insert: 'c', title: 'Count criticals ≥ N' },
	{ group: 'modifiers', label: 'Minimum', insert: 'min', title: 'Minimum result value' },
	{ group: 'modifiers', label: 'Maximum', insert: 'max', title: 'Maximum result value' }
];

export const LOGIC_TOKENS: FormulaToken[] = [
	{ group: 'logic', label: '+', insert: ' + ' },
	{ group: 'logic', label: '-', insert: ' - ' },
	{ group: 'logic', label: '*', insert: ' * ' },
	{ group: 'logic', label: '/', insert: ' / ' },
	{ group: 'logic', label: '%', insert: ' % ' },
	{ group: 'logic', label: '^', insert: ' ^ ' },
	{ group: 'logic', label: '==', insert: ' == ' },
	{ group: 'logic', label: '!=', insert: ' != ' },
	{ group: 'logic', label: '<', insert: ' < ' },
	{ group: 'logic', label: '>', insert: ' > ' },
	{ group: 'logic', label: '<=', insert: ' <= ' },
	{ group: 'logic', label: '>=', insert: ' >= ' },
	{ group: 'logic', label: 'and', insert: ' and ' },
	{ group: 'logic', label: 'or', insert: ' or ' },
	{ group: 'logic', label: 'not', insert: 'not ', title: 'Logical not' },
	{ group: 'logic', label: 'if', insert: 'if ', title: 'Conditional' },
	{ group: 'logic', label: 'then', insert: ' then ' },
	{ group: 'logic', label: 'else', insert: ' else ' }
];

export const FUNCTION_TOKENS: FormulaToken[] = [
	{ group: 'functions', label: 'max', insert: 'max(', title: 'Maximum of values' },
	{ group: 'functions', label: 'min', insert: 'min(', title: 'Minimum of values' },
	{ group: 'functions', label: 'floor', insert: 'floor(', title: 'Round down' },
	{ group: 'functions', label: 'ceil', insert: 'ceil(', title: 'Round up' },
	{ group: 'functions', label: 'round', insert: 'round(', title: 'Round to nearest' },
	{ group: 'functions', label: 'abs', insert: 'abs(', title: 'Absolute value' },
	{ group: 'functions', label: 'clamp', insert: 'clamp(', title: 'Clamp value between bounds' }
];

export const TOKEN_GROUP_LABELS: Record<TokenGroup, string> = {
	dice: 'Dice',
	modifiers: 'Modifiers',
	logic: 'Logic',
	functions: 'Functions',
	variables: 'Variables'
};

export const TOKEN_GROUP_ORDER: TokenGroup[] = [
	'dice',
	'modifiers',
	'logic',
	'functions',
	'variables'
];
