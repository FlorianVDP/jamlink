package animalInvariants

func ValidateAnimal(age int, name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}

	if err := ValidateAge(age); err != nil {
		return err
	}

	return nil
}
