drop database if exists petstay;
create database petstay;

revoke all privileges on database petstay from admin;
drop user if exists admin;
create user admin with password 'admin';
grant all privileges on database petstay to admin;

revoke all privileges on database petstay from fido;
drop user if exists fido;
create user fido with password 'woof';
grant all privileges on database petstay to fido;