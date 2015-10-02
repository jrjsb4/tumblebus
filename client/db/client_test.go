package db

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

func aTestNewConnection(t *testing.T) {
	t.Log("Connecting to mongodb...")
	c := NewConnection()
	defer c.CloseConnection()
}

func TestSchoolDb(t *testing.T) {
	t.Log("Connecting to mongodb...")
	c := NewConnection()

	defer c.CloseConnection()

	schoolList, err := c.ListSchools()
	if err != nil {
		t.Error(err.Error)
	}
	if len(schoolList) != 0 {
		t.Error("Error school collection should be empty")
	}

	t.Log("Adding School #1")
	school := School{
		Name:        "Holy Family",
		Address:     "South Street",
		City:        "Fitchburg",
		State:       "MA",
		ZipCode:     "01420",
		MainPhone:   "978-234-1234",
		ContactName: "Sister Julie",
		Url:         "http://www.holyfamily.org",
	}
	err = c.AddSchool(&school)
	if err != nil {
		schoolList, err = c.ListSchools()
		if err != nil {
			if len(schoolList) != 2 {
				t.Error("Failed to get the entire school list")
			}
		}
		t.Error("Failed to add School 1")
	}

	t.Log("Adding school #2")
	school = School{
		Name:        "Oakmont",
		Address:     "South Asburnham Road",
		City:        "Asburnham",
		State:       "MA",
		ZipCode:     "10452",
		MainPhone:   "978-234-1234",
		ContactName: "Sam Beauvais",
		Url:         "http://www.oakmont.org",
	}
	err = c.AddSchool(&school)
	if err != nil {
		t.Error("Failed to add School 2")
	}

	// Test the EnsureIndex that there are no two schools with the same name
	t.Log("Checking to make sure we can't add duplicate schools")
	school = School{
		Name:        "Holy Family",
		Address:     "South Street",
		City:        "Fitchburg",
		State:       "MA",
		ZipCode:     "01420",
		MainPhone:   "978-234-1234",
		ContactName: "Sister Julie",
		Url:         "http://www.holyfamily.org",
	}
	err = c.AddSchool(&school)
	if err != nil {
		t.Log(err.Error())
	}

	schoolList, err = c.ListSchools()
	if err != nil {
		t.Error("List schools incurred error")
	}

	if len(schoolList) != 2 {
		t.Error("Failed to get the entire school list")
	}

	t.Log("Len: ", len(schoolList))
	for _, school := range schoolList {
		t.Log("School Data: ", school)
	}

	s := &School{}
	s, err = c.FindSchoolByName("Holy Family")
	if err != nil {
		t.Error("Failed to find school")
	}
	t.Log("School: ", s)

	s.City = "Leominster"
	s.ZipCode = "01453"

	err = c.UpdateSchool(s)
	if err != nil {
		t.Error("Failed to update school")
	}

	t.Log("Display updated school information")
	s, err = c.FindSchoolByName("Holy Family")
	if err != nil {
		t.Error("Failed to find school")
	}
	t.Log("School: ", s)

	t.Log("Attempting to removing school")
	err = c.DeleteSchool(s)
	if err != nil {
		t.Error("Failed to remove school")
	}

	schoolList, err = c.ListSchools()
	if err != nil {
		t.Error("List schools incurred error")
	}

	if len(schoolList) != 1 {
		t.Error("Failed to get the entire school list")
	}

	t.Log("Len: ", len(schoolList))
	for _, school := range schoolList {
		t.Log("School Data: ", school)
	}
}

func TestAddClient(t *testing.T) {
	isDrop = false
	t.Log("Connecting to mongodb...")
	c := NewConnection()

	defer c.CloseConnection()

	parent := Parent{
		FirstName:    "Jason",
		LastName:     "Beauvais",
		Address:      "60 Drepanos Drive",
		City:         "Fitchburg",
		State:        "MA",
		ZipCode:      "01420",
		HomePhone:    "978-665-9618",
		MobilePhone:  "978-505-2403",
		EmailAddress: "jrjsb4@gmail.com",
	}

	children := make([]Child, 2)

	children[0].FirstName = "Jacob"
	children[0].LastName = "Beauvais"
	children[0].DOB = time.Date(1996, time.September, 13, 0, 0, 0, 0, time.Local)
	children[0].Age = 19

	children[1].FirstName = "Samuel"
	children[1].LastName = "Beauvais"
	children[1].DOB = time.Date(1999, time.April, 6, 0, 0, 0, 0, time.Local)
	children[1].Age = 16

	paymentInfo := PaymentMethod{
		Method:       CreditCard,
		Frequency:    BiWeekly,
		UnitCost:     15.00,
		StartDate:    time.Date(2015, time.May, 19, 0, 0, 0, 0, time.Local),
		EndDate:      time.Date(2016, time.May, 19, 0, 0, 0, 0, time.Local),
		CcNumber:     "1223 2344 1234 2344",
		SecurityCode: "123",
		CcName:       "Joe Blow",
	}
	t.Log("Looking for Oakmont")
	s, err := c.FindSchoolByName("Oakmont")
	if err != nil {
		t.Error("Failed to find school")
	}
	t.Log("School: ", s)

	err = c.AddClient("Oakmont", &parent, children, &paymentInfo)
	if err != nil {
		t.Error("Failed to insert client info")
	}

	parent = Parent{
		FirstName:    "Mary",
		LastName:     "Keller",
		Address:      "16 Harris Drive",
		City:         "Nrthborogh",
		State:        "MA",
		ZipCode:      "01732",
		HomePhone:    "978-665-9618",
		MobilePhone:  "978-505-2403",
		EmailAddress: "mkeller@gmail.com",
	}

	children[0].FirstName = "Simon"
	children[0].LastName = "Keller"
	children[0].DOB = time.Date(1999, time.April, 13, 0, 0, 0, 0, time.Local)
	children[0].Age = 19

	children[1].FirstName = "Matt"
	children[1].LastName = "Keller"
	children[1].DOB = time.Date(2001, time.November, 30, 0, 0, 0, 0, time.Local)
	children[1].Age = 16

	paymentInfo = PaymentMethod{
		Method:       Cash,
		Frequency:    Weekly,
		UnitCost:     10.00,
		StartDate:    time.Date(2015, time.February, 19, 0, 0, 0, 0, time.Local),
		EndDate:      time.Date(2016, time.February, 19, 0, 0, 0, 0, time.Local),
		CcNumber:     "1223 2344 1234 2344",
		SecurityCode: "123",
		CcName:       "Mary Keller",
	}

	err = c.AddClient("Oakmont", &parent, children, &paymentInfo)
	if err != nil {
		t.Error("Failed to insert client info")
	}

	client := &Client{}
	client, err = c.FindClient("Jaosn", "Beauvais")
	if err != nil {
		t.Error("Unable to find client")
	}

	t.Log("Client: ", client)

	clients := []Client{}
	clients, err = c.ListClients()
	if err != nil {
		t.Error("Unable to find client")
	}

	for _, p := range clients {
		t.Log("Client: ", p)

		for _, ch := range p.Children {
			t.Log("Child: ", ch)
		}
		for _, pay := range p.Payments {
			t.Log("Payment: ", pay)
		}
	}

	var id string
	id, err = c.GetClientId("Jason", "Beauvais")
	if err != nil {
		t.Error("Failed to get id")
	}

	t.Log("Adding three payments to: ", id)
	payment := Payment{
		Method: 1,
		Date:   time.Now(),
		Amount: 10.34,
	}

	err = c.AddPayment(id, &payment)
	if err != nil {
		t.Error("Failed to add payment #1: ", err.Error())
	}

	payment = Payment{
		Method: 3,
		Date:   time.Now(),
		Amount: 33.21,
	}

	err = c.AddPayment(id, &payment)
	if err != nil {
		t.Error("Failed to add payment #2")
	}

	payment = Payment{
		Method: 2,
		Date:   time.Now(),
		Amount: 100.98,
	}

	err = c.AddPayment(id, &payment)
	if err != nil {
		t.Error("Failed to add payment #3")
	}

	client, err = c.FindClient("Jaosn", "Beauvais")
	if err != nil {
		t.Error("Unable to find client")
	}

	t.Log("Client: ", client)
	for _, ch := range client.Children {
		t.Log("Child: ", ch)
	}
	for _, pay := range client.Payments {
		t.Log("Payment: ", pay)
	}
	kidschool := &School{}
	kidschool, err = c.GetSchoolById(bson.ObjectIdHex(client.School))
	if err != nil {
		t.Error("Unable to find school")
	}
	t.Log("School: ", kidschool)

	t.Log("Looking for cliens born 4/1999")
	dob := time.Date(1999, time.April, 0, 0, 0, 0, 0, time.Local)
	clients, err = c.FindClientByDob(dob)
	if err != nil {
		t.Error(err.Error())
	}
	if len(clients) != 2 {
		t.Error("Search found a more than two client with data of 4/1999")
	}
	for _, p := range clients {
		t.Log("AGE Client: ", p)

		for _, ch := range p.Children {
			t.Log("Child: ", ch)
		}
		for _, pay := range p.Payments {
			t.Log("Payment: ", pay)
		}
	}

	t.Log("Looking for Clients born 11/2001")
	dob = time.Date(2001, time.November, 0, 0, 0, 0, 0, time.Local)
	clients, err = c.FindClientByDob(dob)
	if err != nil {
		t.Error(err.Error())
	}
	if len(clients) != 1 {
		t.Error("Search found a more than one client with data of 11/2001")
	}
	for _, p := range clients {
		t.Log("AGE Client: ", p)

		for _, ch := range p.Children {
			t.Log("Child: ", ch)
		}
		for _, pay := range p.Payments {
			t.Log("Payment: ", pay)
		}
	}

	t.Log("Looking for Clients born 1/1971")
	dob = time.Date(1971, time.January, 0, 0, 0, 0, 0, time.Local)
	clients, err = c.FindClientByDob(dob)
	if err != nil {
		t.Error(err.Error())
	}
	if len(clients) > 0 {
		t.Error("Search found a client with data of 1/1971")
	}
	for _, p := range clients {
		t.Log("AGE Client: ", p)

		for _, ch := range p.Children {
			t.Log("Child: ", ch)
		}
		for _, pay := range p.Payments {
			t.Log("Payment: ", pay)
		}
	}
}
