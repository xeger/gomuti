package types

// ChooseCall is a tiebreaker when several allowed calls match a given set of
// parameters. If nil, the default behavior is to choose the most recently
// allowed call.
var ChooseCall func([]Call) Call
