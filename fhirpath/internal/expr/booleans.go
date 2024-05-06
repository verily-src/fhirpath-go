package expr

import "github.com/verily-src/fhirpath-go/fhirpath/system"

func evaluateAnd(left []system.Boolean, right []system.Boolean) system.Collection {
	if len(left) > 0 && len(right) > 0 {
		result := system.Boolean(left[0] && right[0])
		return system.Collection{result}
	}
	// returns false if either boolean is false, regardless of whether or not the other is empty.
	if (len(left) == 1 && !left[0]) || (len(right) == 1 && !right[0]) {
		return system.Collection{system.Boolean(false)}
	}
	return system.Collection{}
}

func evaluateOr(left []system.Boolean, right []system.Boolean) system.Collection {
	if len(left) > 0 && len(right) > 0 {
		result := system.Boolean(left[0] || right[0])
		return system.Collection{result}
	}
	// returns false if either boolean is true, regardless of whether or not the other is empty.
	if (len(left) == 1 && left[0]) || (len(right) == 1 && right[0]) {
		return system.Collection{system.Boolean(true)}
	}
	return system.Collection{}
}

func evaluateXor(left []system.Boolean, right []system.Boolean) system.Collection {
	if len(left) > 0 && len(right) > 0 {
		result := system.Boolean(left[0] != right[0])
		return system.Collection{result}
	}
	return system.Collection{}
}

func evaluateImplies(left []system.Boolean, right []system.Boolean) system.Collection {
	if len(left) > 0 && len(right) > 0 {
		result := system.Boolean(!left[0] || right[0])
		return system.Collection{result}
	}
	// returns true if left is false, or if right is true, regardless of whether or not the other is empty.
	if (len(left) > 0 && !left[0]) || (len(right) > 0 && right[0]) {
		return system.Collection{system.Boolean(true)}
	}
	return system.Collection{}
}
