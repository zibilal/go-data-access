package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/zibilal/go-data-access/apptest/bootstrap"
	"log"
)

func main() {

	bstrp := bootstrap.GetBootstrapped()

	fmt.Println("Connection is successful --> ", bstrp.DbConnector)

	fmt.Println("Storing data")
	StoreData()

	fmt.Println("Data are stored")

	fmt.Println("Fetching the data")
	listUser := FetchData()
	fmt.Println("List user: ", listUser)

	fmt.Println()
	fmt.Println()

	fmt.Println("Fetching the user data")
	userType := FetchData2()
	fmt.Println("User type: ", userType)

	fmt.Println("Updating the user data")
	updatedUserType := UpdateData()
	fmt.Println("User type: ", updatedUserType)
}

type UserType struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	FullName string             `bson:"full_name"`
	Rate     float64            `bson:"rate"`
}

func StoreData() {
	data := []UserType{
		{
			FullName: "Name11", Rate: 0.21,
		},
		{
			FullName: "Name12", Rate: 0.22,
		},
		{
			FullName: "Name13", Rate: 5.57,
		},
	}

	btsrp := bootstrap.GetBootstrapped()
	deService, found := btsrp.ServiceMap["simpleservice"]
	if !found {
		log.Fatal("Simple service is not found")
	}

	err := deService.Serve(context.Background(), data[0], data[1], data[2])
	if err != nil {
		log.Fatalf("Simple service failed due to: %s", err.Error())
	}
}

func FetchData() []UserType {
	data := []UserType{}
	btsrp := bootstrap.GetBootstrapped()
	deService, found := btsrp.ServiceMap["fetchservice"]
	if !found {
		log.Fatalf("Fetch service is not found")
	}

	err := deService.Serve(context.Background(), &data)
	if err != nil {
		log.Fatalf("Fetch service failed due to: %s", err.Error())
	}

	return data
}

func FetchData2() UserType {
	data := UserType{}
	btsrp := bootstrap.GetBootstrapped()
	deService, found := btsrp.ServiceMap["fetchservice"]
	if !found {
		log.Fatal("Fetch service is not found")
	}

	query := struct {
		Id string `query:"_id"`
	}{
		Id: "5c9cf30fe2c7d9fac5d0afcc",
	}

	err := deService.Serve(context.Background(), &data, query)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func UpdateData() UserType {
	data := UserType{
		FullName: "Update Name",
		Rate:     2.015,
	}
	output := UserType{}
	btsrp := bootstrap.GetBootstrapped()
	deService, found := btsrp.ServiceMap["updateservice"]
	if !found {
		log.Fatal("Update service is not found")
	}

	query := struct {
		Id string `query:"_id"`
	}{
		Id: "5c9dc40532d94c2ecc3810ec",
	}

	err := deService.Serve(context.Background(), &output, data, query)
	if err != nil {
		log.Fatal(err)
	}

	return output
}
