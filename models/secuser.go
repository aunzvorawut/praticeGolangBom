package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type SecUser struct {
	Id                     int64     `orm:"pk;auto"`
	Version                int64     `orm:"default(0)"`
	DateCreated            time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated            time.Time `orm:"auto_now;type(datetime)"`
	Enabled                bool      `orm:"default(true)"`
	AccountExpired         bool
	AccountLocked          bool
	PasswordExpired        bool
	Username               string
	Password               string
	Facebookid             string
	ImageProfile           string
	NickNameSocial         string
	NickNameSocialFacebook string
	IsFacebook             bool `orm:"default(false)"`
	Msisdn                 string
	Status                 string
	RegisterDate           time.Time
	LastVisited            time.Time
	IdCard                 string
	PhoneNo                string
	Fixbbid                string
	PasswordNoEncrypt      string
	CountTerminateCode     int
	RoleString             string
	CountSong              int `orm:"default(0)"`
	MsisdnMobile           string
	FfbidBox               string
	TypeTid                string //todo คืออะไร
	ValueVid               string
	Privatev1              string
	Passwordv1             string
	PasswordNoEncryptv1    string
	NetworkType            string
	CurrentPackage         string
}

func init() {
	orm.RegisterModel(new(SecUser))
}

// AddSecUser insert a new SecUser into database and returns
// last inserted Id on success.
func AddSecUser(m *SecUser) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetSecUserById retrieves SecUser by Id. Returns error if
// Id doesn't exist
func GetSecUserById(id int64) (v *SecUser, err error) {
	o := orm.NewOrm()
	v = &SecUser{Id: id}
	if err = o.QueryTable(new(SecUser)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllSecUser retrieves all SecUser matches certain condition. Returns empty list if
// no records exist
func GetAllSecUser(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(SecUser))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []SecUser
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateSecUser updates SecUser by Id and returns error if
// the record to be updated doesn't exist
func UpdateSecUserById(m *SecUser) (err error) {
	o := orm.NewOrm()
	v := SecUser{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteSecUser deletes SecUser by Id and returns error if
// the record to be deleted doesn't exist
func DeleteSecUser(id int64) (err error) {
	o := orm.NewOrm()
	v := SecUser{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&SecUser{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetUserByAccessToken(accessToken string) *UserTokenImei {
	o := orm.NewOrm()
	d := new(UserTokenImei)
	err := o.QueryTable(new(UserTokenImei)).Filter("AccessToken", accessToken).RelatedSel().Limit(1).One(d)
	if err != nil {
		return nil
	}

	return d
}

func GetRealUserByAccessToken(accessToken string) *SecUser {
	o := orm.NewOrm()
	d := new(UserTokenImei)
	err := o.QueryTable(new(UserTokenImei)).Filter("AccessToken", accessToken).RelatedSel().Limit(1).One(d)
	if err != nil {
		return nil
	}

	return d.SecUser

}

func GetSecUserByUsernameAndEnabled(username string, enabled bool) *SecUser {
	o := orm.NewOrm()
	d := new(SecUser)
	err := o.QueryTable(new(SecUser)).Filter("Username", username).Filter("Enabled", enabled).RelatedSel().Limit(1).One(d)
	if err != nil {
		return nil
	}

	return d
}

func GetSecUserByUsername(username string) *SecUser {
	o := orm.NewOrm()
	d := new(SecUser)
	err := o.QueryTable(new(SecUser)).Filter("Username", username).RelatedSel().Limit(1).One(d)
	if err != nil {
		return nil
	}

	return d
}
