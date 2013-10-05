package api

import (
	"encoding/json"
	"fmt"
	"github.com/briansimpson23/restserver/app"
	"net/http"
	"strings"
)

func (c Customer) Create(r *http.Request) (int, []byte) {

	//----------------------------------------------------------
	// get the payload from the HTTP request and parse it
	// into the needed variables
	//----------------------------------------------------------
	c.CompanyName = strings.TrimSpace(r.FormValue("companyName"))
	c.FirstName = strings.TrimSpace(r.FormValue("firstName"))
	c.LastName = strings.TrimSpace(r.FormValue("lastName"))
	c.Address1 = strings.TrimSpace(r.FormValue("address1"))
	c.Address2 = strings.TrimSpace(r.FormValue("address2"))
	c.City = strings.TrimSpace(r.FormValue("city"))
	c.State = strings.TrimSpace(r.FormValue("state"))
	c.ZipCode = strings.TrimSpace(r.FormValue("zipCode"))
	c.Country = strings.TrimSpace(r.FormValue("country"))
	c.Province = strings.TrimSpace(r.FormValue("province"))
	c.PostalCode = strings.TrimSpace(r.FormValue("postalCode"))
	c.Email = strings.TrimSpace(r.FormValue("email"))
	c.MainPhone = strings.TrimSpace(r.FormValue("mainPhone"))
	c.AlternatePhone = strings.TrimSpace(r.FormValue("alternatePhone"))

	//----------------------------------------------------------------------------
	// THIS IS ALL DEBUG OUTPUT
	//
	app.Debug("    +------------------------------------------------------")
	app.Debug("    | CUSTOMER CREATE DEBUG INFORMATION")
	app.Debug("    |")
	app.Debug("    |    INPUT")
	app.Debug(fmt.Sprintf("    |     CompanyName : %s", c.CompanyName))
	app.Debug(fmt.Sprintf("    |       FirstName : %s", c.FirstName))
	app.Debug(fmt.Sprintf("    |        LastName : %s", c.LastName))
	app.Debug(fmt.Sprintf("    |        Address1 : %s", c.Address1))
	app.Debug(fmt.Sprintf("    |        Address2 : %s", c.Address2))
	app.Debug(fmt.Sprintf("    |            City : %s", c.City))
	app.Debug(fmt.Sprintf("    |           State : %s", c.State))
	app.Debug(fmt.Sprintf("    |         ZipCode : %s", c.ZipCode))
	app.Debug(fmt.Sprintf("    |         Country : %s", c.Country))
	app.Debug(fmt.Sprintf("    |        Province : %s", c.Province))
	app.Debug(fmt.Sprintf("    |      PostalCode : %s", c.PostalCode))
	app.Debug(fmt.Sprintf("    |           Email : %s", c.Email))
	app.Debug(fmt.Sprintf("    |       MainPhone : %s", c.MainPhone))
	app.Debug(fmt.Sprintf("    |  AlternatePhone : %s", c.AlternatePhone))
	app.Debug("    |")
	app.Debug("    +------------------------------------------------------")
	//
	// END ALL THE DEBUG OUTPUT
	//----------------------------------------------------------------------------

	//-------------------------------------------------------------------------------------------------------------
	//-- validate that we've received everything all needed input
	//-------------------------------------------------------------------------------------------------------------
	if c.CompanyName == "" {
		app.Debug("missing companyName")
		return 400, BuildErrorMessage(400, "Missing Company Name", 2008)
	}
	if c.FirstName == "" {
		app.Debug("missing firstName")
		return 400, BuildErrorMessage(400, "Missing First Name", 2008)
	}
	if c.LastName == "" {
		app.Debug("missing LastName")
		return 400, BuildErrorMessage(400, "Missing Last Name", 2008)
	}
	if c.Email == "" {
		app.Debug("missing Email Address")
		return 400, BuildErrorMessage(400, "Missing Email Address", 2008)
	}

	//-------------------------------------------------------------------------------------------------------------
	//-- check to see if this email already exists.
	//-------------------------------------------------------------------------------------------------------------
	sql := "SELECT customerID FROM CUSTOMER WHERE email = ?"
	app.Debug(sql)
	row := app.Db.QueryRow(sql, c.Email)
	err := row.Scan(&c.CustomerID)
	if err != nil {
		// if the error is anything except "no rows found then something unplanned is wrong"
		if err.Error() != "sql: no rows in result set" {
			app.Error(fmt.Sprint("%s", err))
			return 500, BuildErrorMessage(500, "Unexpected error", 2008)
		}
	}

	// if this is a local registration then we've got a problem if the memberID is > 0
	if c.CustomerID > 0 {
		app.Error(fmt.Sprintf("email address already exists - [%s]", c.Email))
		return 409, BuildErrorMessage(409, "Duplicate Email Address", 2008)
	}

	//everything is OK so insert the new member record
	sql = "INSERT INTO CUSTOMER " +
		"SET " +
		"companyName=?, firstName = ?, lastName = ?, address1=?, address2=?, city=?, state=?, zipCode = ?, " +
		"country=?, province=?, postalCode=?, email=?, mainPhone=?, alternatePhone=?, dateAdded = now() "
	//s, _ := app.Db.Prepare(sql)
	app.Debug(sql)
	result, err := app.Db.Exec(sql, c.CompanyName, c.FirstName, c.LastName, c.Address1, c.Address2, c.City, c.State,
		c.ZipCode, c.Country, c.Province, c.PostalCode, c.Email, c.MainPhone, c.AlternatePhone)
	if err != nil {
		app.Error(err.Error())
		panic(err)
	}
	c.CustomerID, _ = result.LastInsertId()

	//----------------------------------------------------------
	// build the output
	//----------------------------------------------------------
	resp := Response{}
	resp.Code = 201
	resp.Status = "ok"
	resp.Data = c

	output, err := json.Marshal(resp)
	if err != nil {
		app.Error(err.Error())
		panic(err)
	}

	return 0, output

}
