drop table if exists pets;

create table pets (
  id      bigserial primary key,
  pet_name text,
  pet_type_breed text,
  pet_age text,
  pet_description text,
  owner_name    text,
  owner_address   text,
  owner_phone text,
  owner_email text,
  owner_preferred_contact text,
  vet_name text,
  vet_address text,
  vet_phone text,
  vet_email text,
  emergency_procedure text,
  pet_feeding_times text,
  pet_food_and_quantity text,
  pet_allowed_treats text,
  pet_sleep_location text,
  pet_bedtime text,
  pet_sleeping_habits text,
  pet_ok_w_kids_other_animals text,
  pet_ok_off_leash text,
  pet_sickness_history text,
  pet_desexed text,
  pet_vaccinated text,
  pet_toilet_trained text,
  pet_chews text,
  owner_happy_to_reimburse text,
  pet_anxieties text,
  pet_disallowed_activities text,
  pet_other_details text,
  created_ts  timestamp,
  created_by  text,
  updated_ts  timestamp,
  updated_by  text
);


drop table if exists stay;

create table stay (
  id      bigserial primary key,
  pet_id bigint,
  from_date date,
  to_date date,
  rate numeric(15,4),
  other_fee_1 numeric(15,4),
  other_fee_type_1 text,
  other_fee_2 numeric(15,4),
  other_fee_type_2 text,
  other_fee_3 numeric(15,4),
  other_fee_type_3 text,
  discount numeric(15,4),
  discount_reason numeric(15,4),
  total_charged numeric(15,4),
  total_paid  numeric(15,4),
  created_ts  timestamp,
  created_by  text,
  updated_ts  timestamp,
  updated_by  text
);




--
-- owner
-- -----------
-- owner_name
-- owner_address
-- owner_phone
-- owner_email
-- owner_preferred_contact
--
-- stay
-- ----------
--
-- from_date
-- to_date
-- rate
-- other_fee_1
-- other_fee_type_1
-- other_fee_2
-- other_fee_type_2
-- other_fee_2
-- other_fee_type_2
-- discount
-- total_charged
-- total_paid
--
-- pet
-- -----------
-- pet_name
-- pet_type_breed
-- pet_age
-- pet_description
-- vet_name
-- vet_address
-- vet_phone
-- vet_email
-- emergency_procedure
-- pet_feeding_times
-- pet_food_and_quantity
-- pet_allowed_treats
-- pet_sleep_location
-- pet_bedtime
-- pet_sleeping_habits
-- pet_ok_w_kids_other_animals
-- pet_ok_off_leash
-- pet_sickness_history
-- pet_desexed
-- pet_vaccinated
-- pet_toilet_trained
-- pet_chews
-- owner_happy_to_reimburse
-- pet_anxieties
-- pet_disallowed_activities
-- pet_other_details
--
--
--
-- =====================
-- Owner...
-- Name
-- Address
-- Contact number(s)
-- Email
-- From Date
--
-- Pet...
-- Name
-- Type / Breed
-- Age
-- Description
--
-- Emergency-info
-- Vet Name
-- Address / Location
-- Contact Number(s)
-- Emergency Procedure
--
-- Food...
-- Feeding times
-- What food and quantity?
-- Allowed treats?
--
-- Sleeping...
-- Where does your pet sleep?
-- What time does your pet go to bed?
-- Other sleeping habits
--
-- Misc...
-- How is your pet with kids and other animals?
-- OK to go off-leash in dog parks?
-- Any sickness or history of ill health?
-- Is your pet de-sexed and fully vaccinated?
-- Is your pet toilet trained?
-- Any anxieties? For example, afraid of thunder?
-- Is there anything they are not allowed to do? Go on beds etc.?
-- Anything else I should know to ensure they have a happy stay?



