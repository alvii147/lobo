package isotope

import "fmt"

// Email generates a randomized email addresss.
func Email() string {
	return fmt.Sprintf(
		"%s@%s.%s",
		String(12, true, false, true),
		String(10, true, false, false),
		String(3, true, false, false),
	)
}

// Password generates a randomized password.
func Password() string {
	password := String(20, true, true, true)

	return password
}
