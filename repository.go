package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

// connect to the Db
func init() {
	var err error
	Db, err = sql.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func retrievePet(id int) (pet Pet, err error) {
	pet = Pet{}
	err = Db.QueryRow("select id, name, owner from pet where id = $1", id).Scan(&pet.Id, &pet.Name, &pet.Owner)
	return
}

func retrievePets(offset int, limit int) (pets []Pet, err error) {

	rows, err := Db.Query("select id, name, owner from pet offset $1 limit $2", offset, limit)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var pet Pet
		err := rows.Scan(&pet.Id, &pet.Name, &pet.Owner)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(pet.Id, pet.Name, pet.Owner)
		pets = append(pets, pet)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func (pet *Pet) create() (err error) {
	statement := "insert into pet (name, owner) values ($1, $2) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(pet.Name, pet.Owner).Scan(&pet.Id)
	return
}

func (pet *Pet) updatePet() (err error) {
	_, err = Db.Exec("update pet set name = $2, owner = $3 where id = $1", pet.Id, pet.Name, pet.Owner)
	return
}

func (pet *Pet) deletePet() (err error) {
	_, err = Db.Exec("delete from pet where id = $1", pet.Id)
	return
}

//================================

func retrievePets1(offset int, limit int) (pets []Pet1, err error) {

	rows, err := Db.Query("select id, " +
		" pet_name, pet_type_breed, pet_age, pet_description, owner_name, " +
	    " owner_address, owner_phone, owner_email, owner_preferred_contact, vet_name, " +
		" vet_address, vet_phone, vet_email, emergency_procedure, pet_feeding_times, " +
		" pet_food_and_quantity, pet_allowed_treats, pet_sleep_location, pet_bedtime, pet_sleeping_habits, " +
		" pet_ok_w_kids_other_animals, pet_ok_off_leash, pet_sickness_history, pet_desexed, pet_vaccinated, " +
		" pet_toilet_trained, pet_chews, owner_happy_to_reimburse, pet_anxieties, pet_disallowed_activities, " +
		" pet_other_details " +
		"from pets offset $1 limit $2", offset, limit)
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var pet Pet1
		err := rows.Scan(&pet.Id, &pet.PetName, &pet.PetTypeBreed, &pet.PetAge, &pet.PetDescription,
			&pet.OwnerName, &pet.OwnerAddress, &pet.OwnerPhone, &pet.OwnerEmail, &pet.OwnerPreferredContact,
			&pet.VetName, &pet.VetAddress, &pet.VetPhone, &pet.VetEmail, &pet.EmergencyProcedure,
			&pet.PetFeedingTimes, &pet.PetFoodAndQuantity, &pet.PetAllowedTreats, &pet.PetSleepLocation, &pet.PetBedtime,
			&pet.PetSleepingHabits, &pet.PetOkWithKidsAndOtherAnimals, &pet.PetOkOffLeash, &pet.PetSicknessHistory, &pet.PetDesexed,
			&pet.PetVaccinated, &pet.PetToiletTrained, &pet.PetChews, &pet.OwnerHappyToReimburse, &pet.PetAnxieties,
			&pet.PetDisallowedActivities, &pet.PetOtherDetails)

		if err != nil {
			log.Println(err)
			//return
		}

		log.Println(pet.Id, pet.PetName, pet.PetTypeBreed, pet.PetAge, pet.PetDescription,
			pet.OwnerName, pet.OwnerAddress, pet.OwnerPhone, pet.OwnerEmail, pet.OwnerPreferredContact,
			pet.VetName, pet.VetAddress, pet.VetPhone, pet.VetEmail, pet.EmergencyProcedure,
			pet.PetFeedingTimes, pet.PetFoodAndQuantity, pet.PetAllowedTreats, pet.PetSleepLocation, pet.PetBedtime,
			pet.PetSleepingHabits, pet.PetOkWithKidsAndOtherAnimals, pet.PetOkOffLeash, pet.PetSicknessHistory, pet.PetDesexed,
			pet.PetVaccinated, pet.PetToiletTrained, pet.PetChews, pet.OwnerHappyToReimburse, pet.PetAnxieties,
			pet.PetDisallowedActivities, pet.PetOtherDetails)
		pets = append(pets, pet)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

func (pet *Pet1) create() (err error) {
	statement := "insert into pets (pet_name, pet_type_breed, pet_age, pet_description, owner_name, " +
		" owner_address, owner_phone, owner_email, owner_preferred_contact, vet_name, " +
		" vet_address, vet_phone, vet_email, emergency_procedure, pet_feeding_times, " +
		" pet_food_and_quantity, pet_allowed_treats, pet_sleep_location, pet_bedtime, pet_sleeping_habits, " +
		" pet_ok_w_kids_other_animals, pet_ok_off_leash, pet_sickness_history, pet_desexed, pet_vaccinated, " +
		" pet_toilet_trained, pet_chews, owner_happy_to_reimburse, pet_anxieties, pet_disallowed_activities, " +
		" pet_other_details " +
		//" , created_ts, created_by, updated_ts, updated_by " +
		" ) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10," +
	    " $11,$12,$13,$14,$15,$16,$17,$18,$19,$20," +
	    " $21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31" +
        //" \"\",\"james\",\"\",\"james\" + " +
		" ) returning id"

	stmt, err := Db.Prepare(statement)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(pet.PetName, pet.PetTypeBreed, pet.PetAge, pet.PetDescription,
		pet.OwnerName, pet.OwnerAddress, pet.OwnerPhone, pet.OwnerEmail, pet.OwnerPreferredContact,
		pet.VetName, pet.VetAddress, pet.VetPhone, pet.VetEmail, pet.EmergencyProcedure,
		pet.PetFeedingTimes, pet.PetFoodAndQuantity, pet.PetAllowedTreats, pet.PetSleepLocation, pet.PetBedtime,
		pet.PetSleepingHabits, pet.PetOkWithKidsAndOtherAnimals, pet.PetOkOffLeash, pet.PetSicknessHistory, pet.PetDesexed,
		pet.PetVaccinated, pet.PetToiletTrained, pet.PetChews, pet.OwnerHappyToReimburse, pet.PetAnxieties,
		pet.PetDisallowedActivities, pet.PetOtherDetails).Scan(&pet.Id)
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}