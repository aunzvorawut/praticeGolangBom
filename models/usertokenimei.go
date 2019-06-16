package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type UserTokenImei struct {
	Id          int64     `orm:"pk;auto"`
	Version     int64     `orm:"default(0)"`
	DateCreated time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated time.Time `orm:"auto_now;type(datetime)"`
	SecUser     *SecUser  `orm:"null;rel(fk)"`
	DeviceObj   *Device   `orm:"null;rel(fk)"`
	Imei        string    `orm:"null"`
	AccessToken string    `orm:"null"`
	DateLogin   time.Time `orm:"null"`
	DateExpired time.Time `orm:"null"`
}

func init() {
	orm.RegisterModel(new(UserTokenImei))
}

// AddUserTokenImei insert a new UserTokenImei into database and returns
// last inserted Id on success.
func AddUserTokenImei(m *UserTokenImei) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserTokenImeiById retrieves UserTokenImei by Id. Returns error if
// Id doesn't exist
func GetUserTokenImeiById(id int64) (v *UserTokenImei, err error) {
	o := orm.NewOrm()
	v = &UserTokenImei{Id: id}
	if err = o.QueryTable(new(UserTokenImei)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUserTokenImei retrieves all UserTokenImei matches certain condition. Returns empty list if
// no records exist
func GetAllUserTokenImei(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserTokenImei))
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

	var l []UserTokenImei
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

// UpdateUserTokenImei updates UserTokenImei by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserTokenImeiById(m *UserTokenImei) (err error) {
	o := orm.NewOrm()
	v := UserTokenImei{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUserTokenImei deletes UserTokenImei by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUserTokenImei(id int64) (err error) {
	o := orm.NewOrm()
	v := UserTokenImei{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserTokenImei{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetUserTokenImeiByAceessTokenAndExpired(accessToken string) *UserTokenImei {

	today := time.Now()

	o := orm.NewOrm()
	d := new(UserTokenImei)
	err := o.QueryTable(new(UserTokenImei)).Filter("AccessToken", accessToken).Filter("DateExpired__gte", today).RelatedSel().Limit(1).One(d)
	if err != nil {
		return nil
	}

	return d
}

func GetUserTokenImeiByUserObjAndImeiAndDeviceObj(secUserObj *SecUser, deviceId string, deviceObj *Device) *UserTokenImei {

	o := orm.NewOrm()
	d := new(UserTokenImei)
	err := o.QueryTable(new(UserTokenImei)).Filter("SecUser", secUserObj).Filter("Imei", deviceId).Filter("DeviceObj", deviceObj).RelatedSel().Limit(1).One(d)
	if err != nil {
		return nil
	}

	return d
}

func GetUserTokenImeiByUserObjAndDeviceObj(secUserObj *SecUser, deviceObj *Device) *UserTokenImei {

	o := orm.NewOrm()
	d := new(UserTokenImei)
	err := o.QueryTable(new(UserTokenImei)).Filter("SecUser", secUserObj).Filter("DeviceObj", deviceObj).RelatedSel().Limit(1).One(d)
	if err != nil {
		return nil
	}

	return d
}

func GetUserTokenImeiByUserObjAndminDateLogin(secUserObj *SecUser) *UserTokenImei {

	o := orm.NewOrm()
	d := new(UserTokenImei)
	err := o.QueryTable(new(UserTokenImei)).Filter("secUserObj", secUserObj).OrderBy("DateLogin").RelatedSel().Limit(1).One(d)
	if err != nil {
		return nil
	}

	return d
}

func DeleteAllUserTokenByUserObj(secUserObj *SecUser) (err error) {

	o := orm.NewOrm()
	v := UserTokenImei{SecUser: secUserObj}

	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UserTokenImei{SecUser: secUserObj}); err == nil {
			fmt.Println("Number of records timeLimitScheduleStartEnd deleted in database:", num)
		}
	}
	return

}

func GetAllUserTokenImeiByUserObj(max, offset int, userObj *SecUser) (ml []*UserTokenImei) {

	o := orm.NewOrm()
	qs := o.QueryTable(new(UserTokenImei))

	if _, err := qs.Limit(max, offset).Filter("SecUser", userObj).RelatedSel().All(&ml); err == nil {
		return ml
	}

	return nil

}
