package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ConfigScreen struct {
	Id             int64     `orm:"pk;auto"`
	Version        int64     `orm:"default(0)"`
	DateCreated    time.Time `orm:"auto_now_add;type(datetime)"`
	LastUpdated    time.Time `orm:"auto_now;type(datetime)"`
	Start          time.Time `orm:"null"`
	End            time.Time `orm:"null"`
	Enabled        bool      `orm:"null;default(false)"`
	IsShowOnce     bool      `orm:"null;default(false)"`
	VersionDisplay int       `orm:"null;size(255)"`
	Description    string
}

func init() {
	orm.RegisterModel(new(ConfigScreen))
}

// AddConfigScreen insert a new ConfigScreen into database and returns
// last inserted Id on success.
func AddConfigScreen(m *ConfigScreen) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetConfigScreenById retrieves ConfigScreen by Id. Returns error if
// Id doesn't exist
func GetConfigScreenById(id int64) (v *ConfigScreen, err error) {
	o := orm.NewOrm()
	v = &ConfigScreen{Id: id}
	if err = o.QueryTable(new(ConfigScreen)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllConfigScreen retrieves all ConfigScreen matches certain condition. Returns empty list if
// no records exist
func GetAllConfigScreen(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ConfigScreen))
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

	var l []ConfigScreen
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

// UpdateConfigScreen updates ConfigScreen by Id and returns error if
// the record to be updated doesn't exist
func UpdateConfigScreenById(m *ConfigScreen) (err error) {
	o := orm.NewOrm()
	v := ConfigScreen{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteConfigScreen deletes ConfigScreen by Id and returns error if
// the record to be deleted doesn't exist
func DeleteConfigScreen(id int64) (err error) {
	o := orm.NewOrm()
	v := ConfigScreen{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&ConfigScreen{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetconfigScreenObjBylogicSplashScreen() *ConfigScreen {

	today := time.Now()
	o := orm.NewOrm()
	d := new(ConfigScreen)
	err := o.QueryTable(new(ConfigScreen)).Filter("Enabled", true).OrderBy("-LastUpdated").Filter("Start__lt", today).Filter("End__gt", today).RelatedSel().Limit(1).One(d)
	if err != nil {
		return nil
	}

	return d
}