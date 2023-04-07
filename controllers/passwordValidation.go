package cont

func IsvalidatePass(password string) (string, bool) {

	if len(password) < 8 {

		return "Password is too short", false
	}
	hasUpperCase := false
	hasLowerCase := false
	hasNumbers := false
	hasSpecial := false

	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpperCase = true
		} else if char >= 'a' && char <= 'z' {
			hasLowerCase = true
		} else if char >= '0' && char <= '9' {
			hasNumbers = true
		} else if char >= '!' && char <= '/' {
			hasSpecial = true
		} else if char >= ':' && char <= '@' {
			hasSpecial = true
		}
	}

	if !hasUpperCase {
		return "Password do not contain upperCase Character", false

	}

	if !hasLowerCase {
		return "Password do not contain lowerCase Character", false

	}

	if !hasNumbers {
		return "Password do not contain any numbers", false

	}

	if !hasSpecial {
		return "Password do not contain any special character", false

	}
	return "", true
}