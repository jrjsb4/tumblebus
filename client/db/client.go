package db

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type PaymentType int

const (
	Cash PaymentType = iota
	Check
	CreditCard
	Other
)

type PaymentFrequency int

const (
	Weekly PaymentFrequency = iota
	BiWeekly
	Monthly
	Quarterly
)

// The DB interface defines methods to manipulate the database of clients and schools.
// The reason it is implemented as an interface is to allow other NOSQL or SQL type databases
// to be used in the future. Currently only supporting MongoDB, jrb.
type DB interface {
	ListSchools() (schools []School, err error)
	FindSchoolByName(name string) (school *School, err error)
	ClientExist(FirstName, LastName string) (bool, string)
	GetSchoolById(id bson.ObjectId) (school *School, err error)
	FindClient(firstName, lastName string) (client *Client)
	GetClientId(firstName, lastName string) (string, err error)
	ListClients() (clients []Client, err error)
	FindClinentBySchool(school string) (clients []Client)
	FindClientByDob(dob string) ([]Client, error)
	AddSchool(school *School) (err error)
	AddClient(schoolName string, parent *Parent, children []Child, paymentInfo *PaymentMethod) (err error)
	CreateClient() (id string, err error)
	AddParent(ClientId, FirstName, LastName, Address, City, State, ZipCode, HomePhone, MobilePhone, EmailAddress string) (err error)
	UpdateSchool(school *School) (err error)
	UpdateClient(client *Client) (err error)
	UpdatePaymentMethod(client *Client, paymentInfo *PaymentMethod) (err error)
	AddPayment(id string, payment *Payment) (err error)
	DeleteSchool(school *School) (err error)
	DeleteClient(client *Client) (err error)
}

// Store master mgo Session
type MongoConnection struct {
	session *mgo.Session
}

// Hardcoded Database, Collection, and Hostname variables
var (
	databaseName         = "test"
	clientCollectionName = "clients"
	schoolCollectionName = "schools"
	hostname             = "mongodb://localhost"
	isDrop               = true
)

// Season contains infomation that relates to a school year season
type Season struct {
	Start           time.Time `bson:"start" json:"start"`
	End             time.Time `bson:"end" json:"end"`
	YearToDateTotal float64   `bson:"yeartodatetotal" json:"yeartodatetotal"`
}

// Schoool contains name, address and contact information for the school administrator
type School struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Address     string        `json:"address" bson"address"`
	City        string        `json:"city" bson:"city"`
	State       string        `json:"state" bson:"state"`
	ZipCode     string        `json:"zipcode" bson:"zipcode"`
	MainPhone   string        `json:"mainphone" bson:"mainphone"`
	ContactName string        `json:"contactname" bson:"contactname"`
	Url         string        `json:"url" bson:"url"`
	Seasons     []*Season     `json:"seasons" bson:"seasons"`
}

// PaymentMethod contains the information about how a client intends to pay for a Season
type PaymentMethod struct {
	//Id             bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Method         PaymentType      `bson:"method" json:"method"`
	Frequency      PaymentFrequency `bson:"frequency" json:"frequency"`
	UnitCost       float64          `bson:"unitcost" json:"unitcost"`
	StartDate      time.Time        `bson:"startdate" json:"startdate"`
	EndDate        time.Time        `bson:"enddate" json:"enddate"`
	CcNumber       string           `bson:"ccnumber" json:"ccnumber"`
	ExpirationDate time.Time        `bson:"expirationdate" json:"expirationdate"`
	SecurityCode   string           `bson:"securitycode" json:"securitycode"`
	CcName         string           `bson:"ccname" json:"ccname"`
}

// Payment contains information about an individual payment
type Payment struct {
	//Id     bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Method PaymentType `bson:"method" json:"method"`
	Date   time.Time   `bson:"date" json:"date"`
	Amount float64     `bson:"amount" json:"amount"`
}

// Parent contains name, address and contact information of the parent of the student
type Parent struct {
	FirstName    string `bson:"firstname" json:"firstname"`
	LastName     string `bson:"lastname" json:"lastname"`
	Address      string `bson:"address" json:"address"`
	City         string `bson:"city" json:"city"`
	State        string `bson:"state" json:"state"`
	ZipCode      string `bson:"zipcode" json:"zipcode"`
	HomePhone    string `bson:"homephone" json:"homephone"`
	MobilePhone  string `bson:"mobilephone" json:"mobilephone"`
	EmailAddress string `bson:"emailladdress" json:"emailaddress"`
}

// Child contains name and date of birth of the children of the parent
type Child struct {
	//Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	FirstName string    `bson:"firstname" json:"firstname"`
	LastName  string    `bson:"lastname" json:"lastname"`
	DOB       time.Time `bson:"dob" json:"dob"`
	Age       int       `bson:"age" json:"age"`
}

// Client structure represents the Method that pertains to a single client.
type Client struct {
	Id            bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ParentInfo    Parent        `bson:"parent" json:"parent"`
	Children      []*Child      `bson:"children" json:"children"`
	PaymentMethod PaymentMethod `bson:"paymentmethod" json:"paymentmethod"`
	Payments      []*Payment    `bson:"payments" json:"payments"`
	School        string        `bson:"schoolid" json:"schoolid"`
}

// NewConnection creates a new connection to the mongoDB backend and returns the connection if successful.
func NewConnection() (c *MongoConnection) {
	c = new(MongoConnection)
	if c != nil {
		err := c.createConnection()
		if err != nil {
			panic(err)
		}
	}
	return
}

// createConnection attemps to connect to a mongoDB backend and create the collections.
func (c *MongoConnection) createConnection() (err error) {
	// create a new mongo database session

	c.session, err = mgo.Dial(hostname)
	if err == nil {
		c.session.SetMode(mgo.Monotonic, true)
		// Drop Database
		if isDrop {

			err = c.session.DB(databaseName).DropDatabase()
			if err != nil {
				err = errors.New("Failed to remove old database")
				return
			}
		}

		// Create the database
		dbs := c.session.DB(databaseName)

		// Create the School Collection
		schoolCollection := dbs.C(schoolCollectionName)
		if schoolCollection != nil {
			// Index
			index := mgo.Index{
				Key:      []string{"$text:name"},
				Unique:   true,
				DropDups: true,
			}

			err = schoolCollection.EnsureIndex(index)
			if err != nil {
				errStr := fmt.Sprintf("Collection (%s) could not be indexed properly", schoolCollectionName)
				err = errors.New(errStr)
				return
			}
		} else {
			errStr := fmt.Sprintf("Collection (%s) could not be created", schoolCollectionName)
			err = errors.New(errStr)
			return
		}
		// Create the Client Collection
		clientCollection := dbs.C(clientCollectionName)
		if clientCollection == nil {
			errStr := fmt.Sprintf("Collection (%s) could not be created", clientCollectionName)
			err = errors.New(errStr)
			return
		}
	}
	return
}

// CloseConnection to the mongoDB.
func (c *MongoConnection) CloseConnection() {
	if c.session != nil {
		c.session.Close()
	}
}

// getSeesionAndCollection returns the School and Client connections.
func (c *MongoConnection) getSessionAndCollection() (session *mgo.Session, client, school *mgo.Collection, err error) {
	if c.session != nil {
		session = c.session.Copy()
		client = session.DB(databaseName).C(clientCollectionName)
		school = session.DB(databaseName).C(schoolCollectionName)
	} else {
		err = errors.New("No session found")
	}
	return
}

// getSchoolId returns the School associated with the Id.
func (c *MongoConnection) getSchoolId(name string) (id bson.ObjectId, err error) {
	session, _, schoolCollection, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()
	//oid := bson.ObjectIdHex(id)
	//err = schoolCollection.Find(oid).One(&school)
	school := &School{}
	err = schoolCollection.Find(bson.M{"name": name}).One(&school)

	return school.Id, err
}

// getSchoolById returns the School associated with the Id.
func (c *MongoConnection) GetSchoolById(id bson.ObjectId) (school *School, err error) {
	session, _, schoolCollection, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()
	//oid := bson.ObjectIdHex(id)
	//err = schoolCollection.Find(oid).One(&school)
	err = schoolCollection.Find(bson.M{"_id": id}).One(&school)

	return
}

// FindSchoolByName stored returns a School assicated with the name of the school.
func (c *MongoConnection) FindSchoolByName(name string) (school *School, err error) {
	session, _, schoolCollection, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()

	err = schoolCollection.Find(bson.M{"name": name}).One(&school)

	return
}

// ListSchools returns a list of all the Schools in the collection
func (c *MongoConnection) ListSchools() (schools []School, err error) {
	session, _, schoolCollection, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()

	err = schoolCollection.Find(nil).All(&schools)
	return
}

// AddSchool to the School collection.
func (c *MongoConnection) AddSchool(school *School) (err error) {
	session, _, schoolCollection, err := c.getSessionAndCollection()
	if err != nil {
		return
	}

	defer session.Close()

	err = schoolCollection.Insert(
		bson.M{
			"name":        school.Name,
			"address":     school.Address,
			"city":        school.City,
			"state":       school.State,
			"zipcode":     school.ZipCode,
			"contactname": school.ContactName,
			"mainphone":   school.MainPhone,
			"url":         school.Url,
		},
	)

	if err != nil {
		if mgo.IsDup(err) {
			err = errors.New("Duplicate name exists for the school name")
		}
	}
	return
}

// UpdateSchool updates an existing School collection with new informaion.
func (c *MongoConnection) UpdateSchool(school *School) (err error) {
	session, _, schoolCollection, err := c.getSessionAndCollection()
	if err != nil {
		return
	}

	defer session.Close()

	id, err := c.getSchoolId(school.Name)
	if err != nil {
		return
	}

	err = schoolCollection.Update(
		bson.M{"_id": id},
		bson.M{
			"name":        school.Name,
			"address":     school.Address,
			"city":        school.City,
			"state":       school.State,
			"zipcode":     school.ZipCode,
			"contactname": school.ContactName,
			"mainphone":   school.MainPhone,
			"url":         school.Url,
		},
	)

	return
}

// DeleteSchool removes a School from the collection
func (c *MongoConnection) DeleteSchool(school *School) (err error) {
	session, _, schoolCollection, err := c.getSessionAndCollection()
	if err != nil {
		return
	}

	defer session.Close()

	id, err := c.getSchoolId(school.Name)
	if err != nil {
		return
	}

	err = schoolCollection.Remove(bson.M{"_id": id})
	return
}

func (c *MongoConnection) ClientExist(FirstName, LastName string) (bool, string) {
	session, clientCollection, _, err := c.getSessionAndCollection()
	if err != nil {
		return false, string("")
	}
	defer session.Close()

	client := Client{}

	bsonSelect := bson.M{"parent": bson.M{"$elemMatch": bson.M{"firstname": FirstName, "lastname": LastName}}}

	err = clientCollection.Find(nil).Select(bsonSelect).One(&client)
	if err != nil {
		return false, string("")
	}

	return true, client.Id.Hex()
}

func (c *MongoConnection) CreateClient() (id string, err error) {
	session, clientCollection, _, sessionErr := c.getSessionAndCollection()
	if sessionErr != nil {
		err = sessionErr
		return
	}
	defer session.Close()

	clientId := bson.NewObjectId()

	// Empty Parent
	parent := Parent{}

	// Empty child
	children := []Child{}

	// Empty payment
	payments := []Payment{}

	// PaymentInfo
	paymentInfo := PaymentMethod{}

	// school
	err = clientCollection.Insert(
		bson.M{
			"_id":           clientId,
			"parent":        parent,
			"children":      children,
			"paymentmethod": paymentInfo,
			"payments":      payments,
			"schoolid":      bson.ObjectIdHex("0"),
		},
	)

	return clientId.Hex(), err
}

func (c *MongoConnection) AddParent(ClientId, FirstName, LastName, Address, City, State, ZipCode, HomePhone, MobilePhone, EmailAddress string) (err error) {
	session, clientCollection, _, sessionErr := c.getSessionAndCollection()
	if sessionErr != nil {
		err = sessionErr
		return
	}
	defer session.Close()

	bsonParent := bson.M{
		"firstname":    FirstName,
		"lastname":     LastName,
		"address":      Address,
		"city":         City,
		"state":        State,
		"zipcode":      ZipCode,
		"homephone":    HomePhone,
		"mobilephone":  MobilePhone,
		"emailaddress": EmailAddress,
	}

	err = clientCollection.Update(bson.M{"_id": bson.ObjectIdHex(ClientId)}, bsonParent)

	return
}

// AddCient to the Client collection.
func (c *MongoConnection) AddClient(schoolName string, parent *Parent, children []Child, paymentInfo *PaymentMethod) (err error) {
	school := School{}
	//school, err = c.FindSchoolByName(schoolName)
	//if err != nil {
	session, clientCollection, schoolCollection, sessionErr := c.getSessionAndCollection()
	if sessionErr != nil {
		err = sessionErr
		return
	}
	defer session.Close()

	err = schoolCollection.Find(bson.M{"name": schoolName}).One(&school)
	if err != nil {
		fmt.Println("Unable to find school")
		return
	}

	bsonParent := bson.M{
		"firstname":    parent.FirstName,
		"lastname":     parent.LastName,
		"address":      parent.Address,
		"city":         parent.City,
		"state":        parent.State,
		"zipcode":      parent.ZipCode,
		"homephone":    parent.HomePhone,
		"mobilephone":  parent.MobilePhone,
		"emailaddress": parent.EmailAddress,
	}

	bsonChildren := make([]bson.M, len(children))

	for index, child := range children {
		bsonChild := bson.M{
			"firstname": child.FirstName,
			"lastname":  child.LastName,
			"dob":       child.DOB,
			"age":       child.Age,
		}
		bsonChildren[index] = bsonChild
	}

	bsonPaymentInfo := bson.M{
		"method":         paymentInfo.Method,
		"frequency":      paymentInfo.Frequency,
		"unitcost":       paymentInfo.UnitCost,
		"startdate":      paymentInfo.StartDate,
		"enddate":        paymentInfo.EndDate,
		"ccnumber":       paymentInfo.CcNumber,
		"expirationdate": paymentInfo.ExpirationDate,
		"securitycode":   paymentInfo.SecurityCode,
		"ccname":         paymentInfo.CcName,
	}

	// Enter a empty payment
	bsonPayment := []Payment{}

	err = clientCollection.Insert(
		bson.M{
			"parent":        bsonParent,
			"children":      bsonChildren,
			"paymentmethod": bsonPaymentInfo,
			"payments":      bsonPayment,
			"schoolid":      school.Id.Hex(),
		},
	)

	return
}

// ListClients provides an entire list of all clients in the collection
func (c *MongoConnection) ListClients() (clients []Client, err error) {
	session, clientCollection, _, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()
	err = clientCollection.Find(nil).All(&clients)

	return
}

// FindClient returns a list of clients associated by first and last name.
func (c *MongoConnection) GetClientId(firstName, lastName string) (id string, err error) {
	session, clientCollection, _, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()

	tmp := Client{}

	bsonSelect := bson.M{"parent": bson.M{"$elemMatch": bson.M{"firstname": firstName, "lastname": lastName}}}

	err = clientCollection.Find(nil).Select(bsonSelect).One(&tmp)
	if err != nil {
		return "", err
	}

	id = tmp.Id.Hex()

	return id, err
}

// FindClient returns a list of clients associated by first and last name.
func (c *MongoConnection) FindClient(firstName, lastName string) (client *Client, err error) {
	session, clientCollection, _, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()

	tmp := Client{}

	bsonSelect := bson.M{"parent": bson.M{"$elemMatch": bson.M{"firstname": firstName, "lastname": lastName}}}

	err = clientCollection.Find(nil).Select(bsonSelect).One(&tmp)

	err = clientCollection.Find(bson.M{"_id": tmp.Id}).One(&client)

	return
}

// FindClientBySchool returns a list of clients associated with a particular school
func (c *MongoConnection) FindClinentBySchool(name string) (clients []Client, err error) {
	session, clientCollection, schoolCollection, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()

	school := &School{}
	bsonQuery := bson.M{"name": name}
	err = schoolCollection.Find(bsonQuery).One(&school)
	if err != nil {
		return
	}
	bsonQuery = bson.M{"schoolid": school.Id}
	err = clientCollection.Find(bsonQuery).All(&clients)

	return
}

// FindClientByDob returns the list of names with a set Year and Month. T
// This function is intended to be used for marketing purposes to search all
// children with a specific birth date so we can send out birthday notices to clients
func (c *MongoConnection) FindClientByDob(dob time.Time) (clients []Client, err error) {
	session, clientCollection, _, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()

	// Add a month to the begining month
	endMonth := dob.AddDate(0, 1, 0)

	bsonDateSearch := bson.M{"children.dob": bson.M{"$gt": dob, "$lt": endMonth}}

	err = clientCollection.Find(bsonDateSearch).All(&clients)

	return
}

// UpdateClient infotmation currently stored in the collection
func (c *MongoConnection) UpdateClient(client *Client) (err error) {
	/*
		session, clientCollection, _, err := c.getSessionAndCollection()
		if err != nil {
			return
		}

		defer session.Close()AddPayment(id bson.ObjectId, payment *Payment) (err error)
	*/
	//err = clientCollection.Find(nil).Select(bson.M{
	//	"parent": bson.M{"$elemMatch": bson.M{"firstname": client.Parent.FirstName, "lastname": client.Parent.LastName}},
	//}).One(&client)

	return
}

// UpdatePaymentMethod is used to update the clients payment information associated with a client.
func (c *MongoConnection) UpdatePaymentMethod(id string, paymentInfo *PaymentMethod) (err error) {
	session, clientCollection, _, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()

	bsonPaymentInfo := bson.M{"paymentmethod": bson.M{
		"method":         paymentInfo.Method,
		"frequency":      paymentInfo.Frequency,
		"startdate":      paymentInfo.StartDate,
		"enddate":        paymentInfo.EndDate,
		"ccnumber":       paymentInfo.CcNumber,
		"expirationdate": paymentInfo.ExpirationDate,
		"securitycode":   paymentInfo.SecurityCode,
		"ccname":         paymentInfo.CcName},
	}

	err = clientCollection.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bsonPaymentInfo)
	return
}

// AddPayment to the payments list associated with a particular client
func (c *MongoConnection) AddPayment(id string, payment *Payment) (err error) {
	session, clientCollection, _, err := c.getSessionAndCollection()
	if err != nil {
		return
	}
	defer session.Close()

	bsonPayment := bson.M{"$push": bson.M{"payments": bson.M{"method": payment.Method, "date": payment.Date, "amount": payment.Amount}}}

	err = clientCollection.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bsonPayment)
	return
}

// Delete a client from the collection
func (c *MongoConnection) DeleteClient(client *Client) (err error) {
	return
}
